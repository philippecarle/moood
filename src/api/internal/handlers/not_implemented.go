package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotImplemented is a dummy gin handler func
func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Not implemented yet"})
}
