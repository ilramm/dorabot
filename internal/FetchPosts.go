package internal

import (
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var serverURL string
var botToken string

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get server URL and bot token from environment variables
	serverURL = os.Getenv("SERVER_URL")
	botToken = os.Getenv("BOT_TOKEN")

	if serverURL == "" || botToken == "" {
		log.Fatalf("Missing SERVER_URL or BOT_TOKEN in .env file")
	}
}

func FetchPosts(channelID string) ([]byte, error) {
	url := serverURL + "api/v4/channels/" + channelID + "/posts"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+botToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
