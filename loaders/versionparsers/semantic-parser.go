package versionparsers

import (
	"log"
	"regexp"
)

// SemanticParser ...
type SemanticParser struct {
	prefix string
	regex  *regexp.Regexp
}

// Parse ...
func (parser *SemanticParser) Parse(text string) (string, bool) {
	isVersion := parser.regex.MatchString(text)
	return text, isVersion
}

// NewSemanticParser ...
func NewSemanticParser(prefix string) *SemanticParser {
	if prefix == "" {
		prefix = "V"
	}

	r, err := regexp.Compile("^" + prefix + "[\\-_]{0,1}[0-9].*")

	if err != nil {
		log.Fatalf("Could not create sematic parser: %s", err)
	}

	return &SemanticParser{prefix: prefix, regex: r}
}
