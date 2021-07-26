package fileio

import (
	"strings"
)

// TextScanner Mock to read files from memmory (Unit Testing)
type TextScanner struct {
	text string
}

// Scan ...
func (scn *TextScanner) Scan() bool {

	return len(scn.text) > 0
}

// Text ...
func (scn *TextScanner) Text() string {
	index := strings.Index(scn.text, "\n")
	if index < 0 { //Should simulate scan behiour
		result := scn.text
		scn.text = ""
		return result
	} else if index == 0 {
		scn.text = scn.text[1:]
		return ""
	}

	result := scn.text[0:index]
	scn.text = scn.text[index+1:]

	return result
}

// Close ...
func (scn *TextScanner) Close() {

}

// NewTextScanner ...
func NewTextScanner(text string) *TextScanner {
	return &TextScanner{text: text}
}
