package utils

import (
	"strings"
)

func FormatEndpoint(endpoint string) string {
	endpoint = strings.ReplaceAll(endpoint, " ", "")
	endpoint = strings.ReplaceAll(endpoint, "/", "-")
	endpoint = strings.ReplaceAll(endpoint, "?", "")
	return endpoint
}
