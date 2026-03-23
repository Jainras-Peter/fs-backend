package repository

import (
	"context"
	"fs-backend/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HBLDocRepository interface {
	InsertMany(ctx context.Context, docs []models.HBLDoc) error
	CountTotal(ctx context.Context) (int64, error)
	GetRecent(ctx context.Context, limit int64) ([]models.HBLDoc, error)
	DeleteByID(ctx context.Context, id string) error
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

func (r *hblDocRepository) CountTotal(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

func (r *hblDocRepository) GetRecent(ctx context.Context, limit int64) ([]models.HBLDoc, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []models.HBLDoc
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *hblDocRepository) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
