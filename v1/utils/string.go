package v1_utils

import "strings"

func SanitizeName(name string) string {
	return strings.ReplaceAll(name, "*", "")
}
