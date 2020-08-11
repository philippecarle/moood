package collection

import (
	"context"
	"time"

	"github.com/philippecarle/mood/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// EntriesCollection expose the necessary method to read and write entries
type EntriesCollection struct {
	collection *mongo.Collection
}

// NewEntriesCollection creates an EntriesCollection
func NewEntriesCollection(c *mongo.Database) EntriesCollection {
	return EntriesCollection{collection: c.Collection("entries")}
}

// Insert a new Entry in the database
func (er *EntriesCollection) Insert(e *models.Entry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := er.collection.InsertOne(ctx, e)
	if err != nil {
		return err
	}

	e.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

// Delete an Entry in the database
func (er *EntriesCollection) Delete(e models.Entry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := er.collection.DeleteOne(ctx, bson.M{"id": e.ID})
	if err != nil {
		return err
	}
	return nil
}
