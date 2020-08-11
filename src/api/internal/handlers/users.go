package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philippecarle/moood/api/internal/collection"
	"github.com/philippecarle/moood/api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UsersHandler embeds users-related handler func
type UsersHandler struct {
	collection collection.UsersCollection
}

// NewUserHandler creates an users handler
func NewUserHandler(u collection.UsersCollection) UsersHandler {
	return UsersHandler{collection: u}
}

// Register is the endpoint where a user can register
func (u *UsersHandler) Register(c *gin.Context) {
	var user models.UserRegistration
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedUser := models.User{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.ClearPassword), 8)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	storedUser.EncodedPassword = string(hashedPassword)
	storedUser.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := u.collection.Insert(&storedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := u.collection.FindOneByUserName(storedUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}
