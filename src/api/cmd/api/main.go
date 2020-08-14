package main

import (
	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/philippecarle/moood/api/internal/bus"
	"github.com/philippecarle/moood/api/internal/database"
	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/handlers"
	"github.com/philippecarle/moood/api/internal/middlewares"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	client := database.NewClient()

	db := client.Database("mood")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db.Client().Connect(ctx)
	defer cancel()

	conn := bus.Connection{}
	conn.Init()

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
