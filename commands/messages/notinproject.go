package messages

import (
	"ele/utils/doc"
	"os"
)

// NotInProject ...
func NotInProject() {
	doc.Paragraph(doc.Error("Not an ele project (or any of the parent directories): .ele"))

	doc.Paragraph(doc.Info("run ele init in the root of the project, or directory you want to run your migrations from"))
	os.Exit(1)
}
