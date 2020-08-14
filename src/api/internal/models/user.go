package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Credentials is a struct that models the structure of a user login request body
type Credentials struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

// User is a struct that models the structure of a user in the DB
type User struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	FirstName       string             `json:"firstName" bson:"firstName"`
	LastName        string             `json:"lastName" bson:"lastName"`
	EncodedPassword string             `json:"-" bson:"password"`
	CreatedAt       primitive.DateTime `json:"createdAt" bson:"createdAt"`
}

// UserRegistration is a struct that models the structure of a user registration request body
type UserRegistration struct {
	Email      string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	ClearPassword string `json:"password"`
}
