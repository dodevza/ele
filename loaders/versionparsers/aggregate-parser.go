package versionparsers

import (
	"ele/config"
	"ele/constants"
)

// AggregateParser ...
type AggregateParser struct {
	parsers []VersionParser
}

// Parse ...
func (parser *AggregateParser) Parse(text string) (string, bool) {
	outputstring := text
	isversion := false
	for _, p := range parser.parsers {
		outputstring, isversion = p.Parse(text)
		if isversion {
			return outputstring, isversion
		}
	}

	return text, false
}

// NewAggregateParser ...
func NewAggregateParser(parsers ...VersionParser) *AggregateParser {

	list := make([]VersionParser, 0)
	for _, p := range parsers {
		list = append(list, p)
	}
	return &AggregateParser{parsers: list}
}

// BuildFromConfig ...
func BuildFromConfig(config *config.AppConfig) *AggregateParser {
	if config == nil || config.VersionParsers == nil {
		return NewAggregateParser()
	}

	list := make([]VersionParser, 0)
	for _, conf := range config.VersionParsers {
		switch conf.Type {
		case constants.DATE:
			list = append(list, NewDateParser(conf.Parameter))
			break
		case constants.REGEX:
			list = append(list, NewRegexParser(conf.Parameter))
			break
		case constants.SEMANTIC:
			list = append(list, NewSemanticParser(conf.Parameter))
			break
		default:
			break
		}

	}
	return NewAggregateParser(list...)
}
