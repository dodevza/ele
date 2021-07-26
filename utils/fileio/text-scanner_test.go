package fileio

import "testing"

func Test_SingleLine(t *testing.T) {
	scan := NewTextScanner("A single line")

	readCount := 0
	text := ""
	for scan.Scan() {
		readCount++
		text = scan.Text()
	}

	if readCount != 1 {
		t.Errorf("Read more than once")
	}

	if text != "A single line" {
		t.Errorf("Didn't read full line: %s", text)
	}
}

func Test_MultipleLinse(t *testing.T) {
	scan := NewTextScanner(`First Line
Second Line
Third Line`)

	readCount := 0
	lines := make([]string, 0)
	for scan.Scan() {
		readCount++
		lines = append(lines, scan.Text())
	}

	if readCount != 3 {
		t.Errorf("Read more than once")
	}

	if lines[0] != "First Line" {
		t.Errorf("Didn't read full line: %s", lines[0])
	}

	if lines[1] != "Second Line" {
		t.Errorf("Didn't read full line: %s", lines[1])
	}

	if lines[2] != "Third Line" {
		t.Errorf("Didn't read full line: %s", lines[2])
	}
}
