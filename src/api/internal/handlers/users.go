package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philippecarle/mood/api/internal/auth"
	"github.com/philippecarle/mood/api/internal/collection"
	"github.com/philippecarle/mood/api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UsersHandler embeds users-related handler func
type UsersHandler struct {
	repository collection.UsersCollection
}

// NewUserHandler creates an users handler
func NewUserHandler(u collection.UsersCollection) UsersHandler {
	return UsersHandler{repository: u}
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

	if err := u.repository.Insert(&storedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := u.repository.FindOneByUserName(storedUser.Username)
	fmt.Println(createdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// Login is the endpoint where a user can get a token
func (u *UsersHandler) Login(c *gin.Context) {
	var creds models.Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.repository.FindOneByUserName(creds.Username)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.EncodedPassword), []byte(creds.Password)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{})
	}

	token := auth.NewToken(user)

	c.JSON(http.StatusOK, gin.H{"token": token})
}
