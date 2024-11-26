package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleLatestCommand(w http.ResponseWriter, r *http.Request) {
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

	rawPosts, err := FetchPosts(channelID)
	if err != nil {
		log.Printf("Failed to fetch posts: %v", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	latestTags := ExtractTags(rawPosts)

	registry := "gcr.io/pr-binom"
	formattedOutput := FormatMakeCommand(ParseTagsToMap(latestTags), registry)

	response := map[string]string{
		"response_type": "ephemeral",
		"text":          formattedOutput,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
