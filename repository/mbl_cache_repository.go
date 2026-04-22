package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MBLCacheDocument represents a cached MBL extraction result in the "MBL_Cache" collection
type MBLCacheDocument struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty"`
	FileHash      string                 `bson:"file_hash"`
	Engine        string                 `bson:"engine"`
	MBLNumber     string                 `bson:"mbl_number"`
	ExtractedData map[string]interface{} `bson:"extracted_data"`
	CreatedAt     time.Time              `bson:"created_at"`
}

// MBLCacheRepository defines operations on the "MBL_Cache" collection
type MBLCacheRepository interface {
	FindByFileHashAndEngine(ctx context.Context, fileHash, engine string) (*MBLCacheDocument, error)
	Insert(ctx context.Context, doc *MBLCacheDocument) error
}

type mblCacheRepository struct {
	collection *mongo.Collection
}

// NewMBLCacheRepository creates a new MBLCacheRepository backed by the "MBL_Cache" collection
func NewMBLCacheRepository(db *mongo.Database) MBLCacheRepository {
	return &mblCacheRepository{
		collection: db.Collection("MBL_Cache"),
	}
}

func (r *mblCacheRepository) FindByFileHashAndEngine(ctx context.Context, fileHash, engine string) (*MBLCacheDocument, error) {
	var doc MBLCacheDocument
	err := r.collection.FindOne(ctx, bson.M{
		"file_hash": fileHash,
		"engine":    engine,
	}).Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *mblCacheRepository) Insert(ctx context.Context, doc *MBLCacheDocument) error {
	doc.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}
