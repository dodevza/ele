package fileio

// FileAppender ...
type FileAppender interface {
	Text(line string)
	Close()
}
