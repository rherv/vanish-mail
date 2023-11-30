package util

import "strings"

func IsHtml(contentType string) bool {
	return strings.Contains(contentType, "html")
}

func IsText(contentType string) bool {
	return strings.Contains(contentType, "text")
}
