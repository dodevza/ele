package search

import (
	"ele/utils/doc"
	"ele/utils/text"
	"regexp"
	"strings"
)

// Path ...
type Path struct {
	patterns map[string]string
	sequence pathSequence
	builded  bool
}

type pathItem struct {
	pattern      string
	negate       bool
	regex        *regexp.Regexp
	testfullpath bool
}

type pathSequence []*pathItem

func (sequence pathSequence) match(value string) (bool, bool) {
	matched := false
	ignore := false
	fullpath := value
	paths := strings.Split(value, "/")
	lastPath := value
	if len(paths) > 0 {
		lastPath = paths[len(paths)-1]
	}
	for _, seq := range sequence {

		testValue := fullpath
		if !seq.testfullpath {
			testValue = lastPath
		}
		match := seq.regex.MatchString(testValue)
		if !match {
			continue
		}
		matched = true
		if seq.negate {
			ignore = false
		} else {
			ignore = true
		}
	}
	return matched, ignore
}

// Include ...
func (ign *Path) Include(lines ...string) {
	for _, line := range lines {
		ign.ReadLine(line)
	}
}

// ReadLine ...
func (ign *Path) ReadLine(line string) {
	// A blank line matches no files, so it can serve as a separator for readability.
	if text.IsWhitespaceOrEmpty(line) {
		return
	}

	// A line starting with # serves as a comment. Put a backslash ("\") in front of the first hash for patterns that begin with a hash.
	if strings.HasPrefix(line, "#") {
		return
	}

	if strings.HasPrefix(line, "\\#") {
		line = line[1:]
	}

	//Trailing spaces are ignored unless they are quoted with backslash ("\").
	line = removeTrailingSpaces(line)

	// An optional prefix "!" which negates the pattern;
	negate := false
	if strings.HasPrefix(line, "!") {
		negate = true
		line = line[1:]
	}

	if strings.HasPrefix(line, "\\!") {
		line = line[1:]
	}

	seq := pathItem{pattern: line, negate: negate}
	ign.sequence = append(ign.sequence, &seq)
}

func (ign *Path) buildPatterns() {
	if ign.builded {
		return
	}

	for _, value := range ign.sequence {

		ptrn := ign.convertToRegex(value.pattern)
		if ptrn != nil {
			value.regex = ptrn
			value.testfullpath = true
		}
	}
	ign.builded = true
}

func (ign *Path) convertToRegex(pattern string) *regexp.Regexp {

	regexString := "^.*" + replaceAsterixWithRegex(escapePeriods(pattern))
	regexString = strings.ReplaceAll(regexString, "/", "\\/*")
	regexString = replaceQuesionsWithRegex(regexString)
	regex, err := regexp.Compile(regexString)

	if err != nil {
		doc.Line(doc.Infof("Could not parse ignore patter: %s", pattern))
		return nil
	}
	return regex
}

// Match ...
func (ign *Path) Match(path string) (bool, bool) {
	ign.buildPatterns()

	return ign.sequence.match(path)
}

func removeTrailingSpaces(line string) string {
	if len(line) == 0 {
		return line
	}

	pos := len(line) - 1
	char := line[len(line)-1]
	for char == ' ' && pos > 0 {

		before := line[pos-1]
		if before == '\\' {
			break
		}

		char = before
		pos--
	}
	if pos == 0 {
		return ""
	}

	return line[:pos+1]
}

func replaceQuesionsWithRegex(line string) string {
	return strings.ReplaceAll(line, "?", "[^\\/]")
}
func escapePeriods(line string) string {
	return strings.ReplaceAll(line, ".", "\\.")
}

func replaceAsterixWithRegex(line string) string {
	if line == "*" {
		return ".*"
	}

	if len(line) <= 1 {
		return line
	}

	var sb strings.Builder

	pos := 0
	streak := 0
	for pos < len(line) {
		char := line[pos]

		if char != '*' {
			if streak == 1 || streak == 2 {
				sb.WriteString(".*")
			}
			sb.WriteByte(char)
			streak = 0
		} else {
			streak++
		}

		pos++

		if streak > 0 {
			if pos == len(line) && streak <= 2 {
				sb.WriteString(".*")
			} else if streak == 2 { // Directory **
				sb.WriteString(".*")
			} else if streak > 2 { // Escaped *
				sb.WriteString("*")
			}
		}
	}
	return sb.String()
}

// NewCollection ...
func NewCollection() *Path {
	sequence := make([]*pathItem, 0)
	return &Path{sequence: sequence, builded: false}
}
