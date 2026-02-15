package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ShipperDocument represents a document in the "shippers" collection
type ShipperDocument struct {
	ShipperID      string `bson:"shipper_id"`
	ShipperName    string `bson:"shipper_name"`
	ShipperAddress string `bson:"shipper_address"`
	ShipperContact string `bson:"shipper_contact"`
}

// ShipperRepository defines read operations on the "shippers" collection
type ShipperRepository interface {
	FindByShipperIDs(ctx context.Context, shipperIDs []string) ([]ShipperDocument, error)
}

type shipperRepository struct {
	collection *mongo.Collection
}

// NewShipperRepository creates a new ShipperRepository backed by the "shippers" collection
func NewShipperRepository(db *mongo.Database) ShipperRepository {
	return &shipperRepository{
		collection: db.Collection("shippers"),
	}
}

func (r *shipperRepository) FindByShipperIDs(ctx context.Context, shipperIDs []string) ([]ShipperDocument, error) {
	filter := bson.M{"shipper_id": bson.M{"$in": shipperIDs}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []ShipperDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}
