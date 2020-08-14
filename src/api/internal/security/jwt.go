package security

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateConsumerToken returns a token which grants publishing permissions for one hour
func GenerateConsumerToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
		"mercure": map[string][]string{
			"publish": {"*"},
		},
	})
	tokenString, err := token.SignedString(getPrivateKey())

	if err != nil {
		log.Fatalf("Could not generate JWT, reason: %s", err.Error())
	}

	return tokenString
}

// ValidateToken makes sure the token is valid, not expired and signed with the right method
func ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method used to sign the JWT: %v", token.Header["alg"])
		}

		return getPrivateKey(), nil
	})

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return fmt.Errorf("Could not validate the JWT")
	}

	return nil
}

func getPrivateKey() []byte {
	return []byte(os.Getenv("JWT_PRIVATE_KEY"))
}
