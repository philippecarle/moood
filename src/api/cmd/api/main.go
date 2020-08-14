package main

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/philippecarle/moood/api/internal/bus"
	"github.com/philippecarle/moood/api/internal/config"
	"github.com/philippecarle/moood/api/internal/database"
	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/handlers"
	"github.com/philippecarle/moood/api/internal/middlewares"
)

func main() {
	conf := config.GetConfig()

	r := gin.Default()
	r.Use(cors.Default())

	client := database.NewClient(conf.Mongo.User, conf.Mongo.Password)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := client.Database("mood")
	db.Client().Connect(ctx)

	conn := bus.GetConnection(conf.RabbitMQ.User, conf.RabbitMQ.Password, conf.RabbitMQ.URL, conf.RabbitMQ.Port)

	uc := collections.NewUsersCollection(db)
	ec := collections.NewEntriesCollection(db)

	e := handlers.NewEntriesHandler(&conn, ec)
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

	r.Run()
}
