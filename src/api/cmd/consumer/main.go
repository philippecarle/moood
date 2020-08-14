package main

import (
	"context"
	"time"

	"github.com/philippecarle/moood/api/internal/bus"
	"github.com/philippecarle/moood/api/internal/config"
	"github.com/philippecarle/moood/api/internal/database"
	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/mercure"
)

func main() {
	conf := config.GetConfig()

	c := database.NewClient(conf.Mongo.User, conf.Mongo.Password)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := c.Database("mood")
	db.Client().Connect(ctx)

	conn := bus.GetConnection(conf.RabbitMQ.User, conf.RabbitMQ.Password, conf.RabbitMQ.URL, conf.RabbitMQ.Port)

	msgs := conn.Consume()
	stopChan := make(chan bool)

	ec := collections.NewEntriesCollection(db)
	mc := mercure.NewClient(conf.Mercure.URL, conf.Mercure.Port)

	go bus.ConsumeProcessedEntries(msgs, ec, mc)

	<-stopChan
}
