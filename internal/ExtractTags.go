package internal

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

func ExtractTags(posts []byte) string {
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

	return FormatTags(extractedTags)
}
