package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philippecarle/mood/api/internal/collection"
	"github.com/philippecarle/mood/api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EntriesHandler embeds entries-related handler func
type EntriesHandler struct {
	bus        io.Writer
	repository collection.EntriesCollection
}

// NewEntriesHandler creates an entries handler
func NewEntriesHandler(b io.Writer, r collection.EntriesCollection) EntriesHandler {
	return EntriesHandler{bus: b, repository: r}
}

// PostEntry is a gin handler func
func (e *EntriesHandler) PostEntry(c *gin.Context) {
	var entry models.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := e.repository.Insert(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	j, _ := json.Marshal(entry)
	_, err := e.bus.Write(j)

	if err != nil {
		err := e.repository.Delete(entry)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}
