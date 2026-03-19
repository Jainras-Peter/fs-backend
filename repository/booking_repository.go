package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookingDocument represents a document in the "Booking" collection
type BookingDocument struct {
	MBLNumber   string    `bson:"mbl_number"`
	ShipmentIDs []string  `bson:"shipment_ids"`
	Status      string    `bson:"status"`
	CreatedAt   time.Time `bson:"created_at"`
}

// BookingRepository defines read operations on the "Booking" collection
type BookingRepository interface {
	FindByMBLNumber(ctx context.Context, mblNumber string) (*BookingDocument, error)
	CreateBooking(ctx context.Context, doc *BookingDocument) error
	AddShipmentToBooking(ctx context.Context, mblNumber, shipmentID string) error
	FindByShipmentID(ctx context.Context, shipmentID string) (*BookingDocument, error)
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
