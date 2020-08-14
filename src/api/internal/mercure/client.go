package mercure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/philippecarle/moood/api/internal/models"
)

// Client is a struct in charge of the http interaction with the Mercure Hub
type Client struct {
	baseURL string
}

// NewClient returns a new Mercure Client
func NewClient(url string, port int) Client {
	return Client{baseURL: fmt.Sprintf("%s:%d%s", url, port, "/.well-known/mercure")}
}

// PublishEntryUpdate publishes an entry update to the mercure hub
func (c *Client) PublishEntryUpdate(entry models.FullEntry, tokenString string) error {
	data, _ := json.Marshal(entry)
	form := url.Values{}
	form.Add("topic", "/users/"+entry.UserID.Hex())
	form.Add("data", string(data))

	req, err := http.NewRequest(http.MethodPost, c.baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
