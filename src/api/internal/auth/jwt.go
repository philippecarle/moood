package auth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/philippecarle/mood/api/internal/models"
)

// UserClaims represents the JWT user information
type UserClaims struct {
	Mercure mercureSubscriptions `json:"mercure"`
	jwt.StandardClaims
}

type mercureSubscriptions struct {
	Subscriptions []string `json:"subscribe"`
}

// NewToken generates a JWT from a user struc
func NewToken(u models.User) string {
	claims := UserClaims{
		mercureSubscriptions{Subscriptions: []string{"/users/" + u.ID.String()}},
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "moood",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(os.Getenv("JWT_PRIVATE_KEY")))

	if err != nil {
		panic(err)
	}

	return ss
}
