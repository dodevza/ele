package git

import "testing"

// An asterisk "*" matches anything except a slash. The character "?" matches any one character except "/". The range notation, e.g. [a-zA-Z], can be used to match one of the characters in a range. See fnmatch(3) and the FNM_PATHNAME flag for a more detailed description.

// A leading "**" followed by a slash means match in all directories. For example, "**/foo" matches file or directory "foo" anywhere, the same as pattern "foo". "**/foo/bar" matches file or directory "bar" anywhere that is directly under directory "foo".

// A trailing "/**" matches everything inside. For example, "abc/**" matches all files inside directory "abc", relative to the location of the .gitignore file, with infinite depth.

// A slash followed by two consecutive asterisks then a slash matches zero or more directories. For example, "a/**/b" matches "a/b", "a/x/b", "a/x/y/b" and so on.

// Other consecutive asterisks are considered regular asterisks and will match according to the previous rules.

// -------------------------
// Astirix Replacement tests

func Test_AstrixPatternSingleAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("*")
	if line != ".*" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternBlank_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("")
	if line != "" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternTextWithoutAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("myvaluewithou/asterix")
	if line != "myvaluewithou/asterix" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternDoubleAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("**")
	if line != ".*" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternTextWithAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("hello*")
	if line != "hello.*" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternTextWithDoublwAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("hello**")
	if line != "hello.*" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternTextSingleAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("hel*lo")
	if line != "hel.*lo" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternTextWithMoreThanDoubleAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("hello***")
	if line != "hello.**" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternTextWithMoreThanDoubleComplexAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("***hel***lo***")
	if line != ".**hel.**lo.**" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

func Test_AstrixPatternTextWithMoreThanTrippleComplexAsterix_RegexString(t *testing.T) {
	line := replaceAsterixWithRegex("****hel****lo****")
	if line != ".***hel.***lo.***" {
		t.Fatalf("Did not return regex: %s", line)
	}
}

// ------------------------
// Directory matching with astirix

// A leading "**" followed by a slash means match in all directories. For example, "**/foo" matches file or directory "foo" anywhere, the same as pattern "foo". "**/foo/bar" matches file or directory "bar" anywhere that is directly under directory "foo".
func Test_LeadingAsterix_MathDirectory(t *testing.T) {
	ign := NewIgnore("")
	ign.ReadLine("**/files")

	if matched, _ := ign.MatchDirectory("/sub/dir/filesinhere"); matched == false {
		t.Fatalf("'**/files' did not match directory: '/sub/dir/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/sub/filesinhere"); matched == false {
		t.Fatalf("'**/files' did not match directory: '/sub/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/filesinhere"); matched == false {
		t.Fatalf("'**/files' did not match directory: '/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/corrections/inhere"); matched == true {
		t.Fatalf("'**/files' matched directory: '/corrections/inhere'")
	}
}

// A trailing "/**" matches everything inside. For example, "abc/**" matches all files inside directory "abc", relative to the location of the .gitignore file, with infinite depth.
func Test_TrailingAsterix_MathDirectory(t *testing.T) {
	ign := NewIgnore("")
	ign.ReadLine("files/**")

	if matched, _ := ign.MatchDirectory("/files/inhere"); matched == false {
		t.Fatalf("'files/**' did not match directory: '/files/inhere'")
	}

	if matched, _ := ign.MatchDirectory("/files/somefiles/inhere"); matched == false {
		t.Fatalf("'files/**' did not match directory: '/files/somefiles/inhere'")
	}

	if matched, _ := ign.MatchDirectory("/files"); matched == false {
		t.Fatalf("'files/**' did not match directory: '/files'")
	}

	if matched, _ := ign.MatchDirectory("/corrections/files"); matched == true {
		t.Fatalf("'files/**' matched directory: '/corrections/files'")
	}
	if matched, _ := ign.MatchDirectory("/corrections"); matched == true {
		t.Fatalf("'files/**' matched directory: '/corrections'")
	}
}

// A slash followed by two consecutive asterisks then a slash matches zero or more directories. For example, "a/**/b" matches "a/b", "a/x/b", "a/x/y/b" and so on.
func Test_BetweenSlashesAsterix_MathDirectory(t *testing.T) {
	ign := NewIgnore("")
	ign.ReadLine("a/**/b")

	if matched, _ := ign.MatchDirectory("/a/b"); matched == false {
		t.Fatalf("'a/**/b' did not match directory: '/a/b'")
	}

	if matched, _ := ign.MatchDirectory("/a/x/b"); matched == false {
		t.Fatalf("'a/**/b' did not match directory: '/a/x/b'")
	}

	if matched, _ := ign.MatchDirectory("/a/x/y/b"); matched == false {
		t.Fatalf("'a/**/b' did not match directory: '/a/x/y/b'")
	}

	if matched, _ := ign.MatchDirectory("/b/a"); matched == true {
		t.Fatalf("'a/**/b' matched directory: '/b/a'")
	}
	if matched, _ := ign.MatchDirectory("/a"); matched == true {
		t.Fatalf("'a/**/b' matched directory: '/a'")
	}
}

// Directory Mathing with Asterix sub directories

// A leading "**" followed by a slash means match in all directories. For example, "**/foo" matches file or directory "foo" anywhere, the same as pattern "foo". "**/foo/bar" matches file or directory "bar" anywhere that is directly under directory "foo".
func Test_SubDirLeadingAsterix_MathDirectory(t *testing.T) {
	ign := NewIgnore("/sub/directory")
	ign.ReadLine("**/files")

	if matched, _ := ign.MatchDirectory("/sub/directory/sub/dir/filesinhere"); matched == false {
		t.Fatalf("'**/files' did not match directory: '/sub/directory/sub/dir/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/sub/filesinhere"); matched == false {
		t.Fatalf("'**/files' did not match directory: '/sub/directory/sub/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/filesinhere"); matched == false {
		t.Fatalf("'**/files' did not match directory: '/sub/directory/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/corrections/inhere"); matched == true {
		t.Fatalf("'**/files' matched directory: '/sub/directory/corrections/inhere'")
	}
}

// A trailing "/**" matches everything inside. For example, "abc/**" matches all files inside directory "abc", relative to the location of the .gitignore file, with infinite depth.
func Test_SubDirTrailingAsterix_MathDirectory(t *testing.T) {
	ign := NewIgnore("/sub/directory")
	ign.ReadLine("files/**")

	if matched, _ := ign.MatchDirectory("/sub/directory/files/inhere"); matched == false {
		t.Fatalf("'files/**' did not match directory: '/sub/directory/files/inhere'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/files/somefiles/inhere"); matched == false {
		t.Fatalf("'files/**' did not match directory: '/sub/directory/files/somefiles/inhere'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/files"); matched == false {
		t.Fatalf("'files/**' did not match directory: '/sub/directory/files'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/corrections/files"); matched == true {
		t.Fatalf("'files/**' matched directory: '/sub/directory/corrections/files'")
	}
	if matched, _ := ign.MatchDirectory("/sub/directory/corrections"); matched == true {
		t.Fatalf("'files/**' matched directory: '/sub/directory/corrections'")
	}
}

// A slash followed by two consecutive asterisks then a slash matches zero or more directories. For example, "a/**/b" matches "a/b", "a/x/b", "a/x/y/b" and so on.
func Test_SubDirBetweenSlashesAsterix_MathDirectory(t *testing.T) {
	ign := NewIgnore("/sub/directory")
	ign.ReadLine("a/**/b")

	if matched, _ := ign.MatchDirectory("/sub/directory/a/b"); matched == false {
		t.Fatalf("'a/**/b' did not match directory: '/sub/directory/a/b'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/a/x/b"); matched == false {
		t.Fatalf("'a/**/b' did not match directory: '/sub/directory/a/x/b'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/a/x/y/b"); matched == false {
		t.Fatalf("'a/**/b' did not match directory: '/sub/directory/a/x/y/b'")
	}

	if matched, _ := ign.MatchDirectory("/sub/directory/b/a"); matched == true {
		t.Fatalf("'a/**/b' matched directory: '/sub/directory/b/a'")
	}
	if matched, _ := ign.MatchDirectory("/sub/directory/a"); matched == true {
		t.Fatalf("'a/**/b' matched directory: '/sub/directory/a'")
	}
}
