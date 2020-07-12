package entries

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philippecarle/mood/api/internal/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller embeds entries-related handler func
type Controller struct {
	Channel  *amqp.Channel
	Database *mongo.Database
}

// PostEntry is a gin handler func
func (e *Controller) PostEntry(c *gin.Context) {
	var entry models.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entriesCollection := e.Database.Collection("entries")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := entriesCollection.InsertOne(ctx, entry)
	if err != nil {
		panic(err)
	}

	entry.ID = res.InsertedID.(primitive.ObjectID)

	err = e.Channel.Publish(
		"",
		"spacy",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(entry.Content),
		})

	if err != nil {
		// todo: delete the entry in Mongo
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}
