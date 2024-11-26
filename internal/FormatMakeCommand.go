package internal

import (
	"fmt"
	"strings"
)

func FormatMakeCommand(tags map[string]string, registry string) string {
	var tagStrings []string
	for key, value := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s=%s", key, value))
	}
	tagLines := strings.Join(tagStrings, "\n")

	makeCommand := fmt.Sprintf("make REGISTRY=%s \\\n%s \\\nupdate", registry, strings.Join(tagStrings, " \\\n"))

	return fmt.Sprintf("%s\n\n%s", tagLines, makeCommand)
}
