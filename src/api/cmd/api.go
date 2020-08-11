package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/philippecarle/moood/api/internal/bus"
	"github.com/philippecarle/moood/api/internal/collection"
	"github.com/philippecarle/moood/api/internal/handlers"
	"github.com/philippecarle/moood/api/internal/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Server will run the gin server
func Server() *gin.Engine {
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

	conn := bus.Connection{}
	conn.Init()
	uc := collection.NewUsersCollection(db)
	e := handlers.NewEntriesHandler(&conn, collection.NewEntriesCollection(db))
	u := handlers.NewUserHandler(uc)
	f := middlewares.JWTMiddleWareFactory{UsersCollection: uc}

	authMiddleware := f.NewJWTMiddleWare()

	r.POST("/register", u.Register)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refresh", authMiddleware.RefreshHandler)

	auth := r.Group("")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.POST("/entries", e.PostEntry)
		auth.GET("/entries", handlers.NotImplemented)
	}

	return r
}
