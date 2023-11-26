package util

import (
	"fmt"
	"time"
)

func GenerateTimestamp(creation time.Time) string {
	duration := time.Now().Sub(creation)

	minutes := duration.Minutes()
	hours := duration.Hours()

	if minutes < 1 {
		return "just now"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes ago", int(minutes))
	} else if minutes < 120 {
		return fmt.Sprintf("1 hour ago")
	} else {
		return fmt.Sprintf("%d hours ago", int(hours))
	}
}
