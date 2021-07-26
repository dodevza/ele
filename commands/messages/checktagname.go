package messages

import (
	"ele/utils/doc"
	"os"
)

// CheckTagName ...
func CheckTagName(tagName string) {
	if len(tagName) == 0 {
		doc.Paragraph(doc.Error("No tag name provided"))
		os.Exit(1)
	}
}
