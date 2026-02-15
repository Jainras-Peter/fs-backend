package repository

import (
	"context"
	"fs-backend/models/mbl_schema"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MBLRepository defines operations on the "MBL" collection
type MBLRepository interface {
	InsertMBL(ctx context.Context, doc *mbl_schema.MBLDocument) error
	FindByMBLNumber(ctx context.Context, mblNumber string) (*mbl_schema.MBLDocument, error)
}

type mblRepository struct {
	collection *mongo.Collection
}

// NewMBLRepository creates a new MBLRepository backed by the "MBL" collection
func NewMBLRepository(db *mongo.Database) MBLRepository {
	return &mblRepository{
		collection: db.Collection("MBL"),
	}
}

func (r *mblRepository) InsertMBL(ctx context.Context, doc *mbl_schema.MBLDocument) error {
	doc.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *mblRepository) FindByMBLNumber(ctx context.Context, mblNumber string) (*mbl_schema.MBLDocument, error) {
	var doc mbl_schema.MBLDocument
	err := r.collection.FindOne(ctx, bson.M{"mbl.bill_of_lading_no": mblNumber}).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
