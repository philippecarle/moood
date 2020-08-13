package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/middlewares"
	"github.com/philippecarle/moood/api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EntriesHandler embeds entries-related handler func
type EntriesHandler struct {
	bus        io.Writer
	collection collections.EntriesCollection
}

// NewEntriesHandler creates an entries handler
func NewEntriesHandler(b io.Writer, r collections.EntriesCollection) EntriesHandler {
	return EntriesHandler{bus: b, collection: r}
}

// PostEntry is a gin handler func
func (e *EntriesHandler) PostEntry(c *gin.Context) {
	user := c.MustGet(middlewares.IdentityKey).(models.User)

	var entry models.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	entry.UserID = user.ID

	if err := e.collection.Insert(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	j, _ := json.Marshal(entry)
	_, err := e.bus.Write(j)

	if err != nil {
		err := e.collection.Delete(entry)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry.UserID = primitive.NilObjectID

	c.JSON(http.StatusCreated, entry)
}
