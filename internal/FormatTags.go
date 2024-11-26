package internal

import (
	"fmt"
	"strings"
)

func FormatTags(tags map[string]string) string {
	var tagStrings []string
	for key, value := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(tagStrings, "\n")
}
