package utils

import "strings"

func BuildLikePattern(keywork string) string {
	keywork = strings.ToLower(strings.TrimSpace(keywork))
	var builder strings.Builder
	builder.WriteByte('%')

	for _, r := range keywork {
		builder.WriteRune(r)
		builder.WriteByte('%')
	}
	return builder.String()
}
