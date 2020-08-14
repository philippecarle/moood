package main

import (
	"context"
	"time"

	"github.com/philippecarle/moood/api/internal/bus"
	"github.com/philippecarle/moood/api/internal/database"
	"github.com/philippecarle/moood/api/internal/database/collections"
)

func main() {
	client := database.NewClient()
	db := client.Database("mood")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db.Client().Connect(ctx)
	defer cancel()

	ec := collections.NewEntriesCollection(db)

	conn := bus.Connection{}
	conn.Init()

	msgs := conn.Consume()
	stopChan := make(chan bool)

	go bus.ConsumeProcessedEntries(msgs, ec)

	<-stopChan
}
