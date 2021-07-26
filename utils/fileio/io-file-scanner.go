package fileio

import (
	"bufio"
	"os"
)

// IOFileScanner ...
type IOFileScanner struct {
	file    *os.File
	scanner *bufio.Scanner
}

// Scan ...
func (scn *IOFileScanner) Scan() bool {

	return scn.scanner.Scan()
}

// Text ...
func (scn *IOFileScanner) Text() string {
	return scn.scanner.Text()
}

// Close ...
func (scn *IOFileScanner) Close() {
	scn.file.Close()
}

// NewIOFileScanner ...
func NewIOFileScanner(path string) (*IOFileScanner, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	return &IOFileScanner{file: file, scanner: scanner}, nil
}
