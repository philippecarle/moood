package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient configures and returns a MongoDB client
func NewClient(u string, pw string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017").SetAuth(options.Credential{
		Username: u,
		Password: pw,
	}))
	if err != nil {
		panic(err)
	}

	return client
}
