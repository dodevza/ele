package fileio

// DirectoryReader ...
type DirectoryReader interface {
	SubDirectories() []string
	Files() []string
	Up() (DirectoryReader, bool)
	Name() string
	Path() string

	FileScanner(filename string) (FileScanner, error)
	FileAppender(filename string) FileAppender
	CreateDirectory(folder string) DirectoryReader
	RemoveFile(filename string)
	RemoveDirectory(folder string)
	MoveFile(from string, to string) error
	GetReader(subdirectory string) DirectoryReader
}

// NewInMemoryReader ...
func NewInMemoryReader() DirectoryReader {
	mapreader := NewMapReader()
	return &mapreader
}
