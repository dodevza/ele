package git

import (
	"ele/utils/doc"
	"ele/utils/fileio"
	"ele/utils/text"
	"regexp"
	"strings"
)

// Ignore Specification implementation Will be used for ignoring .gitignore and .eleignore
type Ignore struct {
	patterns          map[string]string
	sequence          ignoreSequence
	directoryMatchers ignoreSequence
	fileMatchers      ignoreSequence
	builded           bool
	path              string
	ignore            bool
}

type ignoreItem struct {
	pattern      string
	negate       bool
	testfullpath bool
	regex        *regexp.Regexp
}

type ignoreSequence []*ignoreItem

func (sequence ignoreSequence) match(value string) (bool, bool) {
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
func (ign *Ignore) Include(lines ...string) {
	for _, line := range lines {
		ign.ReadLine(line)
	}
}

// ReadLine ...
func (ign *Ignore) ReadLine(line string) {
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

	seq := ignoreItem{pattern: line, negate: negate}
	ign.sequence = append(ign.sequence, &seq)
}

func (ign *Ignore) buildPatterns() {
	if ign.builded {
		return
	}
	dirPatterns := make(ignoreSequence, 0)
	filePatterns := make(ignoreSequence, 0)
	for _, value := range ign.sequence {

		ptrn, fullpath := ign.convertToRegex(value.pattern)
		if ptrn != nil {
			seq := &ignoreItem{pattern: value.pattern, negate: value.negate, regex: ptrn, testfullpath: fullpath}
			dirPatterns = append(dirPatterns, seq)
			// If there is a separator at the end of the pattern then the pattern will only match directories, otherwise the pattern can match both files and directories.
			// For example, a pattern doc/frotz/ matches doc/frotz directory, but not a/doc/frotz directory; however frotz/ matches frotz and a/frotz that is a directory (all paths are relative from the .gitignore file).

			if strings.HasSuffix(value.pattern, "/") == false {
				filePatterns = append(filePatterns, seq)
			}
		}
	}

	ign.directoryMatchers = dirPatterns
	ign.fileMatchers = filePatterns
	ign.builded = true
}

func (ign *Ignore) convertToRegex(pattern string) (*regexp.Regexp, bool) {

	if strings.Contains(pattern, "/") == false {
		regexString := replaceAsterixWithRegex(escapePeriods(pattern))
		regexString = replaceQuesionsWithRegex(regexString)
		regexString = "^" + regexString
		regex, err := regexp.Compile(regexString)

		if err != nil {
			doc.Line(doc.Infof("Could not parse ignore pattern: %s", pattern))
			return nil, false
		}
		return regex, false
	}

	prefix := ign.path
	if strings.HasSuffix(pattern, "/") {
		prefix = ".*"
	}

	// A leading "**" followed by a slash means match in all directories. For example, "**/foo" matches file or directory "foo" anywhere, the same as pattern "foo". "**/foo/bar" matches file or directory "bar" anywhere that is directly under directory "foo".
	if strings.HasPrefix(pattern, "**") == false {

		if strings.HasPrefix(pattern, "/") == false && strings.HasSuffix(pattern, "/") == false {
			if strings.HasPrefix(prefix, "/") {
				prefix = fileio.AppendPathChar(prefix)
			} else {
				prefix = fileio.AppendPathChar("/" + prefix)
			}

		}
	}

	if strings.Contains(pattern, "/") == false {
		prefix += ".*__dash___"
	}

	regexString := "^" + prefix + replaceAsterixWithRegex(escapePeriods(pattern))
	regexString = strings.ReplaceAll(regexString, "/", "\\/*")
	regexString = replaceQuesionsWithRegex(regexString)
	regexString = strings.ReplaceAll(regexString, "__dash___", "\\/")
	regex, err := regexp.Compile(regexString)

	if err != nil {
		doc.Line(doc.Infof("Could not parse ignore patter: %s", pattern))
		return nil, false
	}
	return regex, true
}

// MatchDirectory ...
func (ign *Ignore) MatchDirectory(path string) (bool, bool) {
	ign.buildPatterns()

	return ign.directoryMatchers.match(fileio.AppendPathChar(path))
}

// MatchFile ...
func (ign *Ignore) MatchFile(path string) (bool, bool) {
	ign.buildPatterns()

	if len(ign.fileMatchers) == 0 {
		return true, false
	}
	return ign.fileMatchers.match(path)
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

// NewIgnore ...
func NewIgnore(path string) *Ignore {
	sequence := make([]*ignoreItem, 0)
	return &Ignore{path: path, sequence: sequence, builded: false, ignore: true}
}
