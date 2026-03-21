package repository

import (
	"context"
	"fs-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type HBLDocRepository interface {
	InsertMany(ctx context.Context, docs []models.HBLDoc) error
}

type hblDocRepository struct {
	collection *mongo.Collection
}

func NewHBLDocRepository(db *mongo.Database) HBLDocRepository {
	return &hblDocRepository{
		collection: db.Collection("HBL_Doc"),
	}
}

func (r *hblDocRepository) InsertMany(ctx context.Context, docs []models.HBLDoc) error {
	if len(docs) == 0 {
		return nil
	}

	now := time.Now()
	values := make([]interface{}, 0, len(docs))
	for i := range docs {
		docs[i].CreatedAt = now
		docs[i].UpdatedAt = now
		values = append(values, docs[i])
	}

	_, err := r.collection.InsertMany(ctx, values)
	return err
}
