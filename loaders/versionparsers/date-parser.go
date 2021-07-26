package versionparsers

import (
	"github.com/bykof/gostradamus"
)

// DateParser ...
type DateParser struct {
	format string
}

// Parse ...
func (parser *DateParser) Parse(text string) (string, bool) {

	dateTime, err := gostradamus.Parse(text, parser.format)
	if err != nil {
		return text, false
	}
	// Always make it alpha-numerical sortable
	return dateTime.Format("YYYY-MM-DD"), true
}

// NewDateParser ...
func NewDateParser(format string) *DateParser {
	if format == "" {
		format = "YYYY-MM-DD"
	}
	return &DateParser{format: format}
}
