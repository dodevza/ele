package fileio

import (
	"bufio"
	"log"
	"os"
)

// IOFileAppender ...
type IOFileAppender struct {
	file   *os.File
	writer *bufio.Writer
}

// Text ...
func (appender *IOFileAppender) Text(line string) {
	appender.writer.WriteString(line)
	appender.writer.WriteString("\n")
	err := appender.writer.Flush()
	if err != nil {
		log.Fatalf("Appending to file: %s", err)
	}
}

// Close ...
func (appender *IOFileAppender) Close() {
	appender.file.Close()
}

// NewIOFileAppender ...
func NewIOFileAppender(path string) *IOFileAppender {

	var file *os.File
	var fileError error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, fileError = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
	} else {
		file, fileError = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0755)
	}
	if fileError != nil {
		log.Fatal(fileError)
	}

	writer := bufio.NewWriter(file)
	return &IOFileAppender{file: file, writer: writer}
}
