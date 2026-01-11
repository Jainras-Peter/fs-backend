package repository

import (
    "context"
    "fs-backend/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
    Create(ctx context.Context, post *models.Post) error
    FindAll(ctx context.Context) ([]models.Post, error)
}

type postRepository struct {
    collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) PostRepository {
    return &postRepository{
        collection: db.Collection("posts"),
    }
}

func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
    _, err := r.collection.InsertOne(ctx, post)
    return err
}

func (r *postRepository) FindAll(ctx context.Context) ([]models.Post, error) {
    var posts []models.Post
    cursor, err := r.collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    if err = cursor.All(ctx, &posts); err != nil {
        return nil, err
    }
    return posts, nil
}
