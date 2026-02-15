package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ShipmentDocument represents a document in the "shipments" collection
type ShipmentDocument struct {
	ShipmentID       string  `bson:"shipment_id"`
	ShipperID        string  `bson:"shipper_id"`
	GoodsDescription string  `bson:"goods_description"`
	PackagesCount    int     `bson:"packages_count"`
	GrossWeight      float64 `bson:"gross_weight"`
	NetWeight        float64 `bson:"net_weight"`
	Volume           float64 `bson:"volume"`
	MarksAndNumbers  string  `bson:"marks_and_numbers"`
	Measurement      string  `bson:"measurement"`
}

// ShipmentRepository defines read operations on the "shipments" collection
type ShipmentRepository interface {
	FindByShipmentIDs(ctx context.Context, shipmentIDs []string) ([]ShipmentDocument, error)
	FindByShipperIDs(ctx context.Context, shipperIDs []string) ([]ShipmentDocument, error)
}

type shipmentRepository struct {
	collection *mongo.Collection
}

// NewShipmentRepository creates a new ShipmentRepository backed by the "shipments" collection
func NewShipmentRepository(db *mongo.Database) ShipmentRepository {
	return &shipmentRepository{
		collection: db.Collection("shipments"),
	}
}

func (r *shipmentRepository) FindByShipmentIDs(ctx context.Context, shipmentIDs []string) ([]ShipmentDocument, error) {
	filter := bson.M{"shipment_id": bson.M{"$in": shipmentIDs}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []ShipmentDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *shipmentRepository) FindByShipperIDs(ctx context.Context, shipperIDs []string) ([]ShipmentDocument, error) {
	filter := bson.M{"shipper_id": bson.M{"$in": shipperIDs}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []ShipmentDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}
