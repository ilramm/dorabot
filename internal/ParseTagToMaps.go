package internal

import "strings"

func ParseTagsToMap(tags string) map[string]string {
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
