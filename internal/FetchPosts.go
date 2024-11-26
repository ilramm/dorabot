package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	serverURL = "https://mm.binom.dev/"      // Replace with your Mattermost server URL
	botToken  = "1j4rwqf8ijg89mchd9yous1cqw" // Replace with your bot token
)

func FetchPosts(channelID string) ([]byte, error) {
	url := fmt.Sprintf("%sapi/v4/channels/%s/posts", serverURL, channelID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", botToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
