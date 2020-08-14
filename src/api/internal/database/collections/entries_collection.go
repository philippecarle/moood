package collections

import (
	"context"
	"log"
	"time"

	"github.com/philippecarle/moood/api/internal/models"
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

// FindEntry finds an Entry in the database
func (er *EntriesCollection) FindEntry(id primitive.ObjectID) models.FullEntry {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var entry models.FullEntry
	res := er.collection.FindOne(ctx, bson.M{"_id": id})
	if res.Err() == mongo.ErrNoDocuments {
		log.Fatal(res.Err())
	}
	res.Decode(&entry)

	return entry
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

// UpdateEntry update an entry after it's been processed
func (er *EntriesCollection) UpdateEntry(e *models.FullEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := er.collection.ReplaceOne(
		ctx,
		bson.M{"_id": e.ID},
		e,
	)

	return err
}
