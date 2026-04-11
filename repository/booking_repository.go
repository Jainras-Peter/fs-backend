package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookingDocument represents a document in the "Booking" collection
type BookingDocument struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MBLNumber          string             `bson:"mbl_number" json:"mbl_number"`
	ShipmentIDs        []string           `bson:"shipment_ids" json:"shipment_ids"`
	Mode               string             `bson:"mode" json:"mode"` // FCL or LCL
	CarrierName        string             `bson:"carrier_name" json:"carrier_name"`
	EstimatedDeparture string             `bson:"estimated_departure" json:"estimated_departure"`
	EstimatedArrival   string             `bson:"estimated_arrival" json:"estimated_arrival"`
	Status             string             `bson:"status" json:"status"`
	CreatedAt          time.Time          `bson:"created_at" json:"created_at"`
}

// BookingRepository defines read operations on the "Booking" collection
type BookingRepository interface {
	FindByMBLNumber(ctx context.Context, mblNumber string) (*BookingDocument, error)
	CreateBooking(ctx context.Context, doc *BookingDocument) error
	AddShipmentToBooking(ctx context.Context, mblNumber, shipmentID string) error
	FindByShipmentID(ctx context.Context, shipmentID string) (*BookingDocument, error)
	GetAllBookings(ctx context.Context) ([]BookingDocument, error)
	UpdateBookingStatus(ctx context.Context, id primitive.ObjectID, status string) error
	RemoveShipmentFromBooking(ctx context.Context, shipmentID string) error
}

type bookingRepository struct {
	collection *mongo.Collection
}

// NewBookingRepository creates a new BookingRepository backed by the "Booking" collection
func NewBookingRepository(db *mongo.Database) BookingRepository {
	return &bookingRepository{
		collection: db.Collection("Booking"),
	}
}

func (r *bookingRepository) FindByMBLNumber(ctx context.Context, mblNumber string) (*BookingDocument, error) {
	var doc BookingDocument
	err := r.collection.FindOne(ctx, bson.M{"mbl_number": mblNumber}).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *bookingRepository) CreateBooking(ctx context.Context, doc *BookingDocument) error {
	doc.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *bookingRepository) AddShipmentToBooking(ctx context.Context, mblNumber, shipmentID string) error {
	filter := bson.M{"mbl_number": mblNumber}
	update := bson.M{"$addToSet": bson.M{"shipment_ids": shipmentID}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *bookingRepository) FindByShipmentID(ctx context.Context, shipmentID string) (*BookingDocument, error) {
	var doc BookingDocument
	err := r.collection.FindOne(ctx, bson.M{"shipment_ids": shipmentID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil seamlessly when not found
		}
		return nil, err
	}
	return &doc, nil
}

func (r *bookingRepository) GetAllBookings(ctx context.Context) ([]BookingDocument, error) {
	var bookings []BookingDocument
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}

	if bookings == nil {
		bookings = []BookingDocument{}
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateBookingStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *bookingRepository) RemoveShipmentFromBooking(ctx context.Context, shipmentID string) error {
	filter := bson.M{"shipment_ids": shipmentID}
	update := bson.M{"$pull": bson.M{"shipment_ids": shipmentID}}
	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}
