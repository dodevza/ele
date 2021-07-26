package versionparsers

import (
	"log"
	"regexp"
)

// RegexParser ...
type RegexParser struct {
	regex *regexp.Regexp
}

// Parse ...
func (parser *RegexParser) Parse(text string) (string, bool) {
	isVersion := parser.regex.MatchString(text)
	return text, isVersion
}

// NewRegexParser ...
func NewRegexParser(regex string) *RegexParser {
	r, err := regexp.Compile(regex)

	if err != nil {
		log.Fatalf("Could not create regex (%s) parser: %s", regex, err)
	}

	return &RegexParser{regex: r}
}
