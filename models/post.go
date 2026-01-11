package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
    ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title   string             `bson:"title" json:"title"`
    Content string             `bson:"content" json:"content"`
}
