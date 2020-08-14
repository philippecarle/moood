package bus

import (
	"encoding/json"
	"log"

	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/mercure"
	"github.com/philippecarle/moood/api/internal/models"
	"github.com/philippecarle/moood/api/internal/security"
	"github.com/streadway/amqp"
)

// ConsumeProcessedEntries reads entries from the processed channel and update the database
func ConsumeProcessedEntries(msgs <-chan amqp.Delivery, ec collections.EntriesCollection) {
	tokenString := security.GenerateConsumerToken()

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)

		if err := security.ValidateToken(tokenString); err != nil {
			tokenString = security.GenerateConsumerToken()
		}

		entry := &models.Entry{}
		if err := json.Unmarshal(d.Body, entry); err != nil {
			log.Printf("Could not process message: %s", err)
			continue
		}

		fullEntry := ec.FindEntry(entry.ID)
		if err := json.Unmarshal(d.Body, &fullEntry); err != nil {
			log.Printf("Could not process message: %s", err)
			continue
		}

		if err := ec.UpdateEntry(&fullEntry); err != nil {
			log.Printf("Could not update entry in database: %s", err)
			continue
		}

		if err := mercure.PublishEntryUpdate(fullEntry, tokenString); err != nil {
			log.Fatalf("Could not create request to Mercure Hub, reason: %s", err.Error())
		}
	}
}
