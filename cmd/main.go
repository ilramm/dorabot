package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	serverURL = "https://mm.binom.dev/"      // Replace with your Mattermost server URL
	botToken  = "1j4rwqf8ijg89mchd9yous1cqw" // Replace with your bot token
)

func main() {
	http.HandleFunc("/command/latest", handleLatestCommand)

	port := ":8080"
	fmt.Printf("Bot is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleLatestCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	command := r.FormValue("command")
	if command != "/latest" {
		http.Error(w, "Invalid command", http.StatusBadRequest)
		return
	}

	channelID := r.FormValue("channel_id")

	rawPosts, err := fetchPosts(channelID)
	if err != nil {
		log.Printf("Failed to fetch posts: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	latestTags := extractTags(rawPosts)

	registry := "gcr.io/pr-binom"
	formattedOutput := formatMakeCommand(parseTagsToMap(latestTags), registry)

	response := map[string]string{
		"response_type": "ephemeral",
		"text":          formattedOutput,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func fetchPosts(channelID string) ([]byte, error) {
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

func extractTags(posts []byte) string {
	var postList map[string]interface{}
	if err := json.Unmarshal(posts, &postList); err != nil {
		log.Printf("Failed to unmarshal posts: %v", err)
		return "Error parsing posts"
	}

	defaultTags := map[string]string{
		"FRONTEND_TAG":   "413-c60ced72cb9727a049ed04a31ca15f0f1131524d",
		"BACKEND_TAG":    "v2.26.21",
		"CLICKHOUSE_TAG": "v2.5.13",
		"TD_TAG":         "2904-cf81d9218e441cad4554866ecc9cd358503f921b",
		"PROXY_TAG":      "v1.0.34",
		"POSTGRES_TAG":   "v1.0.6",
	}

	labels := map[string]string{
		"FRONTEND_TAG":   "Binom Frontend",
		"BACKEND_TAG":    "Binom API",
		"CLICKHOUSE_TAG": "Binom Clickhouse",
		"TD_TAG":         "Binom Traffic Distribution",
		"PROXY_TAG":      "Binom Proxy",
		"POSTGRES_TAG":   "Binom Postgres",
	}

	extractedTags := make(map[string]string)
	latestTimestamps := make(map[string]int64)

	for key, value := range defaultTags {
		extractedTags[key] = value
		latestTimestamps[key] = 0
	}

	versionRegex := regexp.MustCompile(`(?i)Version:\s*(\S+)`)

	postsMap := postList["posts"].(map[string]interface{})

	for _, post := range postsMap {
		postMap := post.(map[string]interface{})
		props, exists := postMap["props"].(map[string]interface{})
		if !exists {
			continue
		}

		attachments, hasAttachments := props["attachments"].([]interface{})
		if !hasAttachments {
			continue
		}

		createAt, _ := postMap["create_at"].(float64)

		for _, attachment := range attachments {
			attachmentMap := attachment.(map[string]interface{})
			authorName, hasAuthor := attachmentMap["author_name"].(string)
			text, hasText := attachmentMap["text"].(string)

			if !hasAuthor || !hasText {
				continue
			}

			versionMatch := versionRegex.FindStringSubmatch(text)
			if len(versionMatch) < 2 {
				continue
			}

			version := versionMatch[1]
			for tag, label := range labels {
				if strings.Contains(authorName, label) {
					if int64(createAt) > latestTimestamps[tag] {
						extractedTags[tag] = version
						latestTimestamps[tag] = int64(createAt)
					}
				}
			}
		}
	}

	return formatTags(extractedTags)
}

func formatTags(tags map[string]string) string {
	var tagStrings []string
	for key, value := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(tagStrings, "\n")
}

func formatMakeCommand(tags map[string]string, registry string) string {
	var tagStrings []string
	for key, value := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
	}
	tagLines := strings.Join(tagStrings, "\n")

	makeCommand := fmt.Sprintf("make REGISTRY=%s \\\n%s \\\nupdate", registry, strings.Join(tagStrings, " \\\n"))

	return fmt.Sprintf("%s\n\n%s", tagLines, makeCommand)
}

func parseTagsToMap(tags string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(tags, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
