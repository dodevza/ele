package text

import (
	"regexp"
	"strings"
	"unicode"
)

// Trim all whitespaces characters
func Trim(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// IsWhitespaceOrEmpty Checks if a string is empty or only whitespaces
func IsWhitespaceOrEmpty(str string) bool {
	text := Trim(str)
	return len(text) == 0
}

// StartsWithAny Trims and check if the string starts with the paramters
func StartsWithAny(str string, params ...string) bool {
	for _, p := range params {
		if strings.HasPrefix(Trim(str), p) {
			return true
		}
	}
	return false
}

// SplitOnCommonCSVCharacters Split on ;,:| and spaces except if in single or double quotes
func SplitOnCommonCSVCharacters(str string) []string {
	regex, _ := regexp.Compile("[^\\t;,| \"']+|\"([^\"]*)\"|'([^']*)'")

	found := regex.FindAllString(str, -1)

	trimmed := make([]string, 0)

	for _, f := range found {
		trimmed = append(trimmed, strings.Trim(strings.Trim(f, "\""), "'"))
	}

	return trimmed
}

// RegexArray ...
type RegexArray []*regexp.Regexp

// GetMatchExpressions ...
func GetMatchExpressions(searchPaths ...string) RegexArray {
	result := make(RegexArray, 0)
	for _, s := range searchPaths {
		result = append(result, GetMatchExpression(s))
	}
	return result
}

// MatchAny ...
func (list RegexArray) MatchAny(value string) bool {
	if len(list) == 0 {
		return false
	}
	for _, r := range list {
		if r.MatchString(value) {
			return true
		}
	}

	return false

}

// GetMatchExpression ...
func GetMatchExpression(searchPath string) *regexp.Regexp {
	if searchPath == "" {
		searchPath = "*"
	}

	searchString := "(?i)^" + strings.Replace(searchPath, "*", ".*", -1) + "$"
	return regexp.MustCompile(searchString)
}
