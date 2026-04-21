package repository

import (
	"context"
	"fs-backend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InfoToDocRepository interface {
	Create(ctx context.Context, doc *models.InfoToDoc) error
}

type infoToDocRepository struct {
	collection *mongo.Collection
}

func NewInfoToDocRepository(db *mongo.Database) InfoToDocRepository {
	return &infoToDocRepository{
		collection: db.Collection("info-to-doc"),
	}
}

func (r *infoToDocRepository) Create(ctx context.Context, doc *models.InfoToDoc) error {
	result, err := r.collection.InsertOne(ctx, doc)
	if err == nil {
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			doc.ID = oid
		}
	}
	return err
}
