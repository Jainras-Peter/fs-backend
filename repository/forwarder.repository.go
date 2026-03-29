package repository

import (
	"context"
	"fmt"
	"fs-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ForwarderRepository interface {
	Create(ctx context.Context, forwarder *models.Forwarder) error
	FindByUsername(ctx context.Context, username string) (*models.Forwarder, error)
	GetNextForwarderID(ctx context.Context) (string, error)
}

type forwarderRepository struct {
	collection *mongo.Collection
}

func NewForwarderRepository(db *mongo.Database) ForwarderRepository {
	return &forwarderRepository{
		collection: db.Collection("forwarders"),
	}
}

func (r *forwarderRepository) Create(ctx context.Context, forwarder *models.Forwarder) error {
	_, err := r.collection.InsertOne(ctx, forwarder)
	return err
}

func (r *forwarderRepository) FindByUsername(ctx context.Context, username string) (*models.Forwarder, error) {
	var forwarder models.Forwarder
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&forwarder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if no user found, not an error
		}
		return nil, err
	}
	return &forwarder, nil
}

func (r *forwarderRepository) GetNextForwarderID(ctx context.Context) (string, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return "", err
	}
	// e.g., FWD001, FWD002
	return fmt.Sprintf("FWD%03d", count+1), nil
}
