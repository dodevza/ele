package fileio

import "testing"

func Test_AddTextToEmptyTextAppender_UpdateActualText(t *testing.T) {
	file := newMapFile("Original")
	appender := NewTextAppender(file)

	appender.Text("NewLine")

	if *file.content == "Original" {
		t.Fatalf("Did not update text")
	}
}

func Test_AddMultipeLines_UpdateActualText(t *testing.T) {
	file := newMapFile("")
	appender := NewTextAppender(file)

	appender.Text("Line1")
	appender.Text("Line2")

	if *file.content != "Line1\nLine2\n" {
		t.Fatalf("Did add second line, %s", *file.content)
	}
}
