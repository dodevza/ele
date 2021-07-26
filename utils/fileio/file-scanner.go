package fileio

import "strings"

// FileScanner ...
type FileScanner interface {
	Scan() bool
	Text() string
	Close()
}

// ReadAll ...
func ReadAll(scn FileScanner) string {
	var str strings.Builder

	for scn.Scan() {
		line := scn.Text()
		str.WriteString(line)
		str.WriteString("\n")
	}
	return str.String()
}
