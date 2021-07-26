package fileio

// TextAppender ...
type TextAppender struct {
	file *mapFile
}

// Text ...
func (appender *TextAppender) Text(line string) {
	content := *appender.file.content
	content = content + line + "\n"
	appender.file.content = &content
}

// Close ...
func (appender *TextAppender) Close() {

}

// NewTextAppender ...
func NewTextAppender(file *mapFile) *TextAppender {
	return &TextAppender{file: file}
}
