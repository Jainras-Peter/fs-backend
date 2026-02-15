package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookingDocument represents a document in the "Booking" collection
type BookingDocument struct {
	MBLNumber   string   `bson:"mbl_number"`
	ShipmentIDs []string `bson:"shipment_ids"`
	Status      string   `bson:"status"`
}

// BookingRepository defines read operations on the "Booking" collection
type BookingRepository interface {
	FindByMBLNumber(ctx context.Context, mblNumber string) (*BookingDocument, error)
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
