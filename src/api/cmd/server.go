package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/philippecarle/mood/api/internal/bus"
	"github.com/philippecarle/mood/api/internal/collection"
	"github.com/philippecarle/mood/api/internal/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup will run the gin server
func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	cred := options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017").SetAuth(cred))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("mood")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db.Client().Connect(ctx)
	defer cancel()
	er := collection.NewEntriesCollection(db.Collection("entries"))

	conn := bus.Connection{}
	conn.Init()
	h := handlers.NewEntriesHandler(&conn, er)

	r.POST("/register", handlers.NotImplemented)
	r.POST("/login", handlers.NotImplemented)
	r.POST("/entries", h.PostEntry)
	r.GET("/entries", handlers.NotImplemented)

	return r
}
