package collection

import (
	"context"
	"time"

	"github.com/philippecarle/moood/api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UsersCollection expose the necessary method to read and write users
type UsersCollection struct {
	collection *mongo.Collection
}

// NewUsersCollection creates an UsersCollection
func NewUsersCollection(c *mongo.Database) UsersCollection {
	return UsersCollection{collection: c.Collection("users")}
}

// Insert a new User in the database
func (er *UsersCollection) Insert(e *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := er.collection.InsertOne(ctx, e)
	if err != nil {
		return err
	}

	e.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

// FindOneByUserName finds a User in the database by their username
func (er *UsersCollection) FindOneByUserName(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User

	err := er.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	return user, err
}
