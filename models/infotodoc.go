package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoToDoc struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Filename  string             `json:"filename" bson:"filename"`
	Type      string             `json:"type" bson:"type"` // e.g. "Bill of Lading"
	Data      interface{}        `json:"data" bson:"data"`
	URL       string             `json:"url" bson:"url"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type InfoToDocCreateRequest struct {
	Template string      `json:"Template"`
	Filename string      `json:"Filename"` // passed as reference name from frontend
	Data     interface{} `json:"data"`
}
