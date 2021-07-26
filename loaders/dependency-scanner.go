package loaders

import (
	"ele/utils/fileio"
	"ele/utils/text"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// DependencyScanner Scans file for requirements
type DependencyScanner struct {
	commentprefixes []string
	requireMatcher  *regexp.Regexp
}

// Scan ...
func (scn *DependencyScanner) Scan(scanner fileio.FileScanner) []string {
	requirements := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if text.IsWhitespaceOrEmpty(line) {
			continue
		}

		for _, match := range scn.requireMatcher.FindAllStringSubmatch(line, -1) {
			lineRequirements := text.SplitOnCommonCSVCharacters(match[2])

			requirements = append(requirements, lineRequirements...)
		}

	}

	return requirements
}

func createDependencyScanner(fileExt string) *DependencyScanner {

	commentprefixes := getCommentPrefixes(fileExt)
	commentFlags := make([]string, 0)
	for _, pfx := range commentprefixes {
		commentFlags = append(commentFlags, fmt.Sprintf("\\%s", pfx))
	}
	commentFlagsString := strings.Join(commentFlags, "")

	regexString := fmt.Sprintf("(?i)[%s][ \t]*(require[:	 ]+)([\\w\"\"].*)", commentFlagsString)

	requirementMatcher, err := regexp.Compile(regexString)
	if err != nil {
		log.Fatalf("Could not compile requirements matcher")
	}

	return &DependencyScanner{commentprefixes: commentprefixes, requireMatcher: requirementMatcher}
}

func getCommentPrefixes(filename string) []string {
	if strings.HasSuffix(filename, ".sql") {
		return []string{"-"}
	}

	return []string{"#"}
}
