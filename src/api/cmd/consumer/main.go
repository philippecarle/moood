package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/philippecarle/moood/api/internal/bus"
	"github.com/philippecarle/moood/api/internal/database"
	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/models"
	"github.com/streadway/amqp"
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

	go processMessages(msgs, ec)

	<-stopChan
}

func processMessages(msgs <-chan amqp.Delivery, ec collections.EntriesCollection) {
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		entry := &models.Entry{}
		if err := json.Unmarshal(d.Body, entry); err != nil {
			log.Printf("Could not process message: %s", err)
			return
		}

		fullEntry := ec.FindEntry(entry.ID)
		if err := json.Unmarshal(d.Body, &fullEntry); err != nil {
			log.Printf("Could not process message: %s", err)
			return
		}

		err := ec.UpdateEntry(&fullEntry)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
