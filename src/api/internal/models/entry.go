package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Entry defines the structure of an Entry
type Entry struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Content string             `json:"content" binding:"required" bson:"content"`
}
