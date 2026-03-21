package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ShipperDocument represents a document in the "shippers" collection
type ShipperDocument struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShipperID      string             `bson:"shipper_id" json:"shipper_id"`
	ShipperName    string             `bson:"shipper_name" json:"shipper_name"`
	ShipperAddress string             `bson:"shipper_address" json:"shipper_address"`
	ShipperContact string             `bson:"shipper_contact" json:"shipper_contact"`
}

// ShipperRepository defines read operations on the "shippers" collection
type ShipperRepository interface {
	FindByShipperIDs(ctx context.Context, shipperIDs []string) ([]ShipperDocument, error)
	CreateShipper(ctx context.Context, doc ShipperDocument) (*mongo.InsertOneResult, error)
	FindAllShippers(ctx context.Context) ([]ShipperDocument, error)
	UpdateShipper(ctx context.Context, id primitive.ObjectID, doc map[string]interface{}) (*mongo.UpdateResult, error)
	DeleteShipper(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error)
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

func (r *shipperRepository) CreateShipper(ctx context.Context, doc ShipperDocument) (*mongo.InsertOneResult, error) {
	return r.collection.InsertOne(ctx, doc)
}

func (r *shipperRepository) FindAllShippers(ctx context.Context) ([]ShipperDocument, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []ShipperDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	// Return empty array instead of nil if no documents
	if docs == nil {
		docs = []ShipperDocument{}
	}
	return docs, nil
}

func (r *shipperRepository) UpdateShipper(ctx context.Context, id primitive.ObjectID, doc map[string]interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": doc}
	return r.collection.UpdateOne(ctx, filter, update)
}

func (r *shipperRepository) DeleteShipper(ctx context.Context, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return r.collection.DeleteOne(ctx, filter)
}
