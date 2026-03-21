package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HBLDoc struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Filename  string             `bson:"filename" json:"filename"`
	URL       string             `bson:"url" json:"url"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
