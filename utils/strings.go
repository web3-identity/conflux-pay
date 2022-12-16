package utils

import "regexp"

func ReplaceEmoji(source string, replacement string) string {
	reg := regexp.MustCompile(`[\x{00A1}-\x{00AF}]|[\x{1F000}-\x{1FFFF}]|[\x{2000}-\x{3000}]`)
	converted := reg.ReplaceAll([]byte(source), []byte(replacement))
	return string(converted)
}
