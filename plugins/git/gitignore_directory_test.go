package git

import "testing"

// 1. Matching relative to the directory of Ignore
// If there is a separator at the beginning or middle (or both) of the pattern, then the pattern is relative to the directory level of the particular .gitignore file itself. Otherwise the pattern may also match at any level below the .gitignore level.
//
func Test_SeperatorAtTheBegining_MatchesRelationToIgnore(t *testing.T) {

	ign := NewIgnore("")
	ign.ReadLine("/dir")

	if matched, _ := ign.MatchDirectory("/directions"); matched == false {
		t.Fatalf("'/dir' did not match directory: '/directions'")
	}

	if matched, _ := ign.MatchDirectory("/corrections"); matched == true {
		t.Fatalf("'/dir' matched directory: '/corrections'")
	}
}

func Test_SeperatorInTheMiddle_MatchesRelationToIgnore(t *testing.T) {

	ign := NewIgnore("")
	ign.ReadLine("dir/files")

	if matched, _ := ign.MatchDirectory("/dir/filesinhere"); matched == false {
		t.Fatalf("'dir/files' did not match directory: '/dir/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/corrections/filesinhere"); matched == true {
		t.Fatalf("'dir/files' matched directory: '/corrections/files'")
	}
}

func Test_NestedSeperatorAtTheBegining_MatchesRelationToIgnore(t *testing.T) {

	ign := NewIgnore("/sub")
	ign.ReadLine("/dir")

	if matched, _ := ign.MatchDirectory("/sub/directions"); matched == false {
		t.Fatalf("'/dir' did not match directory: '/sub/directions'")
	}

	if matched, _ := ign.MatchDirectory("/sub/corrections"); matched == true {
		t.Fatalf("'/dir' matched directory: '/sub/corrections'")
	}
}

func Test_NestedSeperatorInTheMiddle_MatchesRelationToIgnore(t *testing.T) {

	ign := NewIgnore("/sub")
	ign.ReadLine("dir/files")

	if matched, _ := ign.MatchDirectory("/sub/dir/filesinhere"); matched == false {
		t.Fatalf("'dir/files' did not match directory: '/sub/dir/filesinhere'")
	}

	if matched, _ := ign.MatchDirectory("/corrections/filesinhere"); matched == true {
		t.Fatalf("'dir/files' matched directory: '/sub/corrections/files'")
	}
}

func Test_SepAtEnd_MatchesRelationToRoot(t *testing.T) {

	ign := NewIgnore("")
	ign.ReadLine("bin/")

	if matched, _ := ign.MatchDirectory("/bin/files"); matched == false {
		t.Fatalf("'bin/' did not match directory: '/bin/files'")
	}
}

func Test_SepAtEnd_MatchesNestedRelations(t *testing.T) {

	ign := NewIgnore("")
	ign.ReadLine("bin/")

	if matched, _ := ign.MatchDirectory("/files/bin"); matched == false {
		t.Fatalf("'bin/' did not match directory: '/files/bin'")
	}
}
