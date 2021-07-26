package git

import "testing"

// An asterisk "*" matches anything except a slash. The character "?" matches any one character except "/". The range notation, e.g. [a-zA-Z], can be used to match one of the characters in a range. See fnmatch(3) and the FNM_PATHNAME flag for a more detailed description.

func Test_Questionmark_ReplacedQuestionmark(t *testing.T) {
	line := replaceQuesionsWithRegex("?")

	if line != "[^\\/]" {
		t.Fatalf("Did not replace '?': %s", line)
	}
}

func Test_MultipleQuestionmark_ReplacedAllQuesionmarks(t *testing.T) {
	line := replaceQuesionsWithRegex("?ub?titute")

	if line != "[^\\/]ub[^\\/]titute" {
		t.Fatalf("Did not replace all '?': %s", line)
	}
}

func Test_QuestionmarkPattern_MatchSubdirectory(t *testing.T) {
	ign := NewIgnore("/sub/directory")
	ign.ReadLine("a/???/b")

	if matched, _ := ign.MatchDirectory("/sub/directory/a/abc/b"); matched == false {
		t.Fatalf("'a/???/b' did not match directory: '/sub/directory/a/abc/b'")
	}
}

func Test_QuestionmarkPattern_DontMatchSubdirectory(t *testing.T) {
	ign := NewIgnore("/sub/directory")
	ign.ReadLine("a/???/b")

	if matched, _ := ign.MatchDirectory("/sub/directory/a/ab/b"); matched == true {
		t.Fatalf("'a/???/b' matched directory: '/sub/directory/a/ab/b'")
	}
}
