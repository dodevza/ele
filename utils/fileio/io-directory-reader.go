package fileio

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// IODirectoryReader ...
type IODirectoryReader struct {
	path string
	name string
}

// Path ...
func (reader *IODirectoryReader) Path() string {
	return reader.path
}

// Name ...
func (reader *IODirectoryReader) Name() string {
	return reader.name
}

// SubDirectories ...
func (reader *IODirectoryReader) SubDirectories() []string {
	files, err := ioutil.ReadDir(reader.path)
	if err != nil {
		log.Fatal(err)
	}

	directories := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			directories = append(directories, f.Name())
		}
	}

	return directories
}

// Up ...
func (reader *IODirectoryReader) Up() (DirectoryReader, bool) {

	parent := filepath.Dir(reader.path)
	if parent != reader.path {
		return NewIODirectoryReader(parent), true
	}

	return nil, false

}

// Files ...
func (reader *IODirectoryReader) Files() []string {
	files, err := ioutil.ReadDir(reader.path)
	if err != nil {
		log.Fatal(err)
	}

	filesNames := make([]string, 0)
	for _, f := range files {
		if f.IsDir() == false {
			filesNames = append(filesNames, f.Name())
		}
	}

	return filesNames
}

// FileScanner ...
func (reader *IODirectoryReader) FileScanner(filename string) (FileScanner, error) {
	return NewIOFileScanner(reader.path + "/" + filename)
}

// FileAppender ...
func (reader *IODirectoryReader) FileAppender(filename string) FileAppender {
	return NewIOFileAppender(reader.path + "/" + filename)
}

// RemoveFile ...
func (reader *IODirectoryReader) RemoveFile(filename string) {
	DeleteFileIfExists(reader.path + "/" + filename)
}

// MoveFile ...
func (reader *IODirectoryReader) MoveFile(from string, to string) error {
	return os.Rename(reader.path+"/"+from, reader.path+"/"+to)
}

// RemoveDirectory ...
func (reader *IODirectoryReader) RemoveDirectory(folder string) {
	DeleteDirectoryIfExists(reader.path + "/" + folder)
}

// CreateDirectory ...
func (reader *IODirectoryReader) CreateDirectory(name string) DirectoryReader {
	CreateDirectoryIfNotExists(reader.path + "/" + name)
	return reader.GetReader(name)
}

// GetReader ...
func (reader *IODirectoryReader) GetReader(subdirectory string) DirectoryReader {
	return NewIODirectoryReader(reader.path + "/" + subdirectory)
}

// NewIODirectoryReader ...
func NewIODirectoryReader(path string) *IODirectoryReader {
	sections := strings.Split(path, "/")
	name := ""
	if len(sections) > 0 {
		name = sections[len(sections)-1]
	}

	return &IODirectoryReader{path: path, name: name}
}
