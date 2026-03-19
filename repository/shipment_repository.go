package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	GetNextShipmentID(ctx context.Context) (string, error)
	GetAllShipments(ctx context.Context) ([]ShipmentDocument, error)
	InsertShipment(ctx context.Context, doc *ShipmentDocument) error
	UpdateShipment(ctx context.Context, shipmentID string, doc *ShipmentDocument) error
	DeleteShipment(ctx context.Context, shipmentID string) error
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

func (r *shipmentRepository) GetNextShipmentID(ctx context.Context) (string, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "shipment_id", Value: -1}})
	var lastShipment ShipmentDocument
	err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&lastShipment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "SHIP001", nil
		}
		return "", err
	}

	var lastID int
	fmt.Sscanf(lastShipment.ShipmentID, "SHIP%d", &lastID)
	newID := fmt.Sprintf("SHIP%03d", lastID+1)
	return newID, nil
}

func (r *shipmentRepository) GetAllShipments(ctx context.Context) ([]ShipmentDocument, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
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

func (r *shipmentRepository) InsertShipment(ctx context.Context, doc *ShipmentDocument) error {
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *shipmentRepository) UpdateShipment(ctx context.Context, shipmentID string, doc *ShipmentDocument) error {
	filter := bson.M{"shipment_id": shipmentID}
	
	updateDoc := *doc
	updateDoc.ShipmentID = shipmentID
	
	update := bson.M{"$set": updateDoc}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *shipmentRepository) DeleteShipment(ctx context.Context, shipmentID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"shipment_id": shipmentID})
	return err
}
