package repository

import (
	"context"
	"fs-backend/models/hbl_schema"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// HBLRepository defines operations on the "HBL" collection
type HBLRepository interface {
	InsertHBL(ctx context.Context, doc *hbl_schema.HBLDocument) error
	UpdateHBL(ctx context.Context, hblNumber string, data hbl_schema.HBLData) error
	FindByHBLNumber(ctx context.Context, hblNumber string) (*hbl_schema.HBLDocument, error)
}

type hblRepository struct {
	collection *mongo.Collection
}

// NewHBLRepository creates a new HBLRepository backed by the "HBL" collection
func NewHBLRepository(db *mongo.Database) HBLRepository {
	return &hblRepository{
		collection: db.Collection("HBL"),
	}
}

func (r *hblRepository) InsertHBL(ctx context.Context, doc *hbl_schema.HBLDocument) error {
	doc.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *hblRepository) UpdateHBL(ctx context.Context, hblNumber string, data hbl_schema.HBLData) error {
	filter := bson.M{"hbl_number": hblNumber}
	update := bson.M{"$set": bson.M{"hbl": data}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *hblRepository) FindByHBLNumber(ctx context.Context, hblNumber string) (*hbl_schema.HBLDocument, error) {
	filter := bson.M{"hbl_number": hblNumber}
	var doc hbl_schema.HBLDocument
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	return &doc, err
}
