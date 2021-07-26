package doc

import (
	"strings"
)

// FitLeft ...
func FitLeft(value string, size int) string {
	diff := len(value) - size
	if diff >= 0 {
		return value[:size]
	}

	var sb strings.Builder
	sb.WriteString(value)
	for i := 0; i < -diff; i++ {
		sb.WriteString(" ")
	}
	return sb.String()
}

// FitRight ...
func FitRight(value string, size int) string {
	diff := len(value) - size
	if diff >= 0 {
		return value[:size]
	}

	var sb strings.Builder

	for i := 0; i < -diff; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString(value)
	return sb.String()
}

// GenerateString ...
func GenerateString(value string, times int) string {
	var sb strings.Builder

	for i := 0; i < times; i++ {
		sb.WriteString(value)
	}
	return sb.String()
}
