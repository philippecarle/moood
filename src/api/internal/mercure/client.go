package mercure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/philippecarle/moood/api/internal/models"
)

// PublishEntryUpdate publishes an entry update to the mercure hub
func PublishEntryUpdate(entry models.FullEntry, tokenString string) error {
	data, _ := json.Marshal(entry)
	form := url.Values{}
	form.Add("topic", "/users/"+entry.UserID.Hex())
	form.Add("data", string(data))
	url := fmt.Sprintf("%s:%s%s", os.Getenv("MERCURE_HUB_URL"), os.Getenv("MERCURE_HUB_POST"), "/.well-known/mercure")
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
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
