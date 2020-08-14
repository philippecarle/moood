package bus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/philippecarle/moood/api/internal/database/collections"
	"github.com/philippecarle/moood/api/internal/models"
	"github.com/streadway/amqp"
)

// ConsumeProcessedEntries reads entries from the processed channel and update the database
func ConsumeProcessedEntries(msgs <-chan amqp.Delivery, ec collections.EntriesCollection) {
	tokenString := generateToken()

	for d := range msgs {
		if err := validateToken(tokenString); err != nil {
			log.Printf("JWT Invalid: %s", err.Error())
			tokenString = generateToken()
		}

		log.Printf("Received a message: %s", d.Body)
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
		data, _ := json.Marshal(fullEntry)
		form := url.Values{}
		form.Add("topic", "/users/"+fullEntry.UserID.Hex())
		form.Add("data", string(data))
		url := fmt.Sprintf("%s:%s%s", os.Getenv("MERCURE_HUB_URL"), os.Getenv("MERCURE_HUB_POST"), "/.well-known/mercure")
		req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
		log.Println(form.Encode())
		if err != nil {
			log.Fatalf("Could not create request to Mercure Hub, reason: %s", err.Error())
		}
		req.Header.Add("Authorization", "Bearer "+tokenString)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response: ", err)
		}

		r, _ := ioutil.ReadAll(resp.Body)
		log.Println(string([]byte(r)))
	}
}

func generateToken() string {
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

func validateToken(tokenString string) error {
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
