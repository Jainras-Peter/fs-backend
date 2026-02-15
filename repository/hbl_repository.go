package repository

import (
	"context"
	"fs-backend/models/hbl_schema"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// HBLRepository defines operations on the "HBL" collection
type HBLRepository interface {
	InsertHBL(ctx context.Context, doc *hbl_schema.HBLDocument) error
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
