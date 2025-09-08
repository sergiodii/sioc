package text

import "strings"

func Sanitize(text string) string {
	return strings.ReplaceAll(text, "*", "")
}
