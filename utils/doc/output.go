package doc

import (
	"fmt"
	"io"
	"os"
)

var output io.Writer

// SetOutput ...
func SetOutput(w io.Writer) {
	output = w
}

func writer() io.Writer {
	if output == nil {
		output = os.Stdout
	}
	return output
}

// Heading ...
func Heading(params ...interface{}) {

	w := writer()
	fmt.Fprint(w, "\n")
	Concat(params...)
	fmt.Fprint(w, "\n\n")
}

// Progress ...
func Progress(params ...interface{}) {
	Concat(params...)
}

// Complete ...
func Complete(status interface{}) {
	w := writer()
	fmt.Fprint(w, " - ")
	fmt.Fprint(w, status)
	fmt.Fprint(w, "\n")
}

// Paragraph ...
func Paragraph(params ...interface{}) {

	w := writer()
	fmt.Fprint(w, "\n")
	Concat(params...)
	fmt.Fprint(w, "\n\n")
}

// NewLine ...
func NewLine() {
	w := writer()
	fmt.Fprint(w, "\n")
}

// Line ...
func Line(params ...interface{}) {

	w := writer()
	Concat(params...)
	fmt.Fprint(w, "\n")
}

// Concat ...
func Concat(params ...interface{}) {

	w := writer()
	for _, p := range params {
		fmt.Fprintf(w, "%s ", p)
	}
}

// NewSlice ...
func NewSlice() []interface{} {
	return make([]interface{}, 0)
}

// Cell Print Cell formatting
func Cell(value string) {
	fmt.Fprintf(writer(), "%s", Data(value))
}

// Divider ...
func Divider(length int) {
	writer := writer()
	fmt.Fprint(writer, GenerateString("-", length))
	fmt.Fprintln(writer)
}
