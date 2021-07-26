package fileio

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

// MapDirectoryReader ...
type MapDirectoryReader struct {
	directories map[string]*MapDirectoryReader
	files       map[string]*mapFile
	parent      *MapDirectoryReader
	path        string
	name        string
}

type mapFile struct {
	content *string
}

func newMapFile(content string) *mapFile {
	return &mapFile{content: &content}
}

// Path ...
func (reader *MapDirectoryReader) Path() string {
	return reader.path
}

// Name ...
func (reader *MapDirectoryReader) Name() string {
	return reader.name
}

// SubDirectories ...
func (reader *MapDirectoryReader) SubDirectories() []string {

	keys := make([]string, 0)

	for key := range reader.directories {

		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// Up ...
func (reader *MapDirectoryReader) Up() (DirectoryReader, bool) {

	if reader.parent != nil {
		return reader.parent, true
	}

	return nil, false
}

// Files ...
func (reader *MapDirectoryReader) Files() []string {
	keys := make([]string, 0)

	for key := range reader.files {

		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// FileScanner ...
func (reader *MapDirectoryReader) FileScanner(filename string) (FileScanner, error) {
	file, ok := reader.files[filename]
	if !ok {
		return nil, fmt.Errorf("No File: %s", filename)
	}
	return NewTextScanner(*file.content), nil
}

// FileAppender ...
func (reader *MapDirectoryReader) FileAppender(filename string) FileAppender {
	file, ok := reader.files[filename]
	if !ok {
		file = newMapFile("")
		reader.files[filename] = file
	}
	return NewTextAppender(file)
}

// RemoveFile ...
func (reader *MapDirectoryReader) RemoveFile(filename string) {
	delete(reader.files, filename)
}

// MoveFile ...
func (reader *MapDirectoryReader) MoveFile(from string, to string) error {
	paths := strings.Split(from, "/")
	curDir := reader
	pathCount := len(paths)
	for index, p := range paths {
		if index == pathCount-1 {
			file, ok := curDir.files[p]
			if !ok {
				return fmt.Errorf("No File Found")
			}
			reader.moveFile(file, to)
			delete(curDir.files, p)
		} else {
			nextDir, ok := curDir.directories[p]
			if !ok {
				return fmt.Errorf("No File Found")
			}
			curDir = nextDir
		}
	}
	return nil
}

func (reader *MapDirectoryReader) moveFile(file *mapFile, to string) {
	paths := strings.Split(to, "/")
	curDir := reader
	pathCount := len(paths)
	for index, p := range paths {
		if index == pathCount-1 {
			_, ok := curDir.files[p]
			if ok {
				log.Fatalf("File already exist")
			}

			curDir.files[p] = file

		} else {
			nextDir := curDir.AddDirectory(p)

			curDir = nextDir
		}
	}
}

// CreateDirectory ...
func (reader *MapDirectoryReader) CreateDirectory(name string) DirectoryReader {
	dr := reader.AddDirectory(name)
	return dr
}

// RemoveDirectory ...
func (reader *MapDirectoryReader) RemoveDirectory(folder string) {
	delete(reader.directories, folder)
}

// GetReader ...
func (reader *MapDirectoryReader) GetReader(subdirectory string) DirectoryReader {
	sub, ok := reader.directories[subdirectory]

	if ok == false {
		log.Fatalf("No Directory")
	}

	return sub
}

// AddDirectory ...
func (reader *MapDirectoryReader) AddDirectory(name string) *MapDirectoryReader {
	dir, ok := reader.directories[name]

	if ok == false {
		dirValue := NewMapReader()
		dir = &dirValue
		dir.parent = reader
		if reader.path != "/" {
			dir.path = reader.path + "/" + name
		} else {
			dir.path = "/" + name
		}

		dir.name = name
		reader.directories[name] = dir
	}

	return dir
}

// AddFile ...
func (reader *MapDirectoryReader) AddFile(name string, content string) *MapDirectoryReader {

	_, ok := reader.files[name]

	if ok == false {
		contentFile := newMapFile(content)
		reader.files[name] = contentFile
	}

	return reader
}

// Back Go back one directory
func (reader *MapDirectoryReader) Back() *MapDirectoryReader {
	return reader.parent
}

// NewMapReader ...
func NewMapReader() MapDirectoryReader {
	directories := make(map[string]*MapDirectoryReader, 0)
	files := make(map[string]*mapFile, 0)
	reader := MapDirectoryReader{directories: directories, files: files, path: "/", name: ""}

	return reader
}
