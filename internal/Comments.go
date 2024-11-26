package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Comments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	command := r.FormValue("command")
	if command != "/ahead" {
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

	latestTags, commentsMap := ExtractTagsComments(rawPosts)

	if len(latestTags) == 0 {
		http.Error(w, "No tags found", http.StatusBadRequest)
		return
	}

	response := formatTagsWithComments(latestTags, commentsMap)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response_type": "ephemeral",
		"text":          response,
	})
}

// formatTagsWithComments formats the tags and associated comments into a readable output.
func formatTagsWithComments(tags map[string]string, comments map[string][]string) string {
	var output strings.Builder
	for tag, version := range tags {
		output.WriteString(fmt.Sprintf("%s=%s\n", tag, version))
		if commentList, exists := comments[tag]; exists {
			for _, comment := range commentList {
				output.WriteString(fmt.Sprintf("  Comment: %s\n", comment))
			}
		}
	}
	return output.String()
}
