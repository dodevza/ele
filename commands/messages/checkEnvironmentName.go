package messages

import (
	"ele/utils/doc"
	"os"
)

// CheckEnvironmentName ...
func CheckEnvironmentName(envName string) {
	if len(envName) == 0 {
		doc.Paragraph(doc.Error("No environment name provided"))
		os.Exit(1)
	}
}
