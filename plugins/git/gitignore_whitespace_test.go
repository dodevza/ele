package git

import "testing"

func Test_TrailingSpaces_Removed(t *testing.T) {
	line := removeTrailingSpaces("This is a line   ")

	if line != "This is a line" {
		t.Fatalf("Did not remove trailing spaces correctly: '%s'", line)
	}
}

func Test_WhitespaceLines_Removed(t *testing.T) {
	line := removeTrailingSpaces("   ")

	if len(line) != 0 {
		t.Fatalf("Did not remove all whitespaces: '%s'", line)
	}
}

func Test_LineWithEscapeCharacter_StaysTheSame(t *testing.T) {
	line := removeTrailingSpaces("This is a line  \\ ")

	if line != "This is a line  \\ " {
		t.Fatalf("Removed whitespaces: '%s'", line)
	}
}

func Test_LineWithEscapeCharacterNotEscapingAllSpaces_StaysTheSame(t *testing.T) {
	line := removeTrailingSpaces("This is a line  \\    ")

	if line != "This is a line  \\ " {
		t.Fatalf("Didn't remove correct amount of whitespace: '%s'", line)
	}
}
