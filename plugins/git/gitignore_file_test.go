package git

import "testing"

func Test_FileName_BinFilesExcluded(t *testing.T) {
	ign := NewIgnore("")
	ign.ReadLine("file.dll")

	if matched, _ := ign.MatchFile("/bin/file.dll"); matched == false {
		t.Fatalf("'file.dll' didn't match directory: '/bin/file.dll'")
	}
}

func Test_IgnoreInBinFileName_BinFilesExcluded(t *testing.T) {
	ign := NewIgnore("/bin")
	ign.ReadLine("file.dll")

	if matched, _ := ign.MatchFile("/bin/file.dll"); matched == false {
		t.Fatalf("'file.dll' didn't match directory: '/bin/file.dll'")
	}
}

func Test_WildCardMathingPattern_IncludedFile(t *testing.T) {
	ign := NewIgnore("/bin")
	ign.ReadLine("*.dll")

	if matched, _ := ign.MatchFile("/bin/file.dll"); matched == false {
		t.Fatalf("'file.dll' didn't match  directory: '/bin/file.dll'")
	}
}

func Test_WildCardAtTheEndMathingPattern_IncludedFile(t *testing.T) {
	ign := NewIgnore("/bin")
	ign.ReadLine("file.*")

	if matched, _ := ign.MatchFile("/bin/file.dll"); matched == false {
		t.Fatalf("'file.dll' didn't match  directory: '/bin/file.dll'")
	}
}

func Test_WildCardAllMathingPattern_IncludedFile(t *testing.T) {
	ign := NewIgnore("/bin")
	ign.ReadLine("*.*")

	if matched, _ := ign.MatchFile("/bin/file.dll"); matched == false {
		t.Fatalf("'file.dll' didn't match  directory: '/bin/file.dll'")
	}
}

func Test_SingleWildCard_IncludedFile(t *testing.T) {
	ign := NewIgnore("/bin")
	ign.ReadLine("*")

	if matched, _ := ign.MatchFile("/bin/file.dll"); matched == false {
		t.Fatalf("'file.dll' didn't match  directory: '/bin/file.dll'")
	}
}

func Test_IgnoreInBinFileName_FilesLikeFilenameShouldNotMatch(t *testing.T) {
	ign := NewIgnore("/bin")
	ign.ReadLine("file.dll")

	if matched, _ := ign.MatchFile("/bin/afile.dll"); matched == true {
		t.Fatalf("'file.dll' matched directory: '/bin/afile.dll'")
	}
}

func Test_FileInRoot_FileInRootIncluded(t *testing.T) {
	ign := NewIgnore("")
	ign.ReadLine("file.dll")

	if matched, _ := ign.MatchFile("/file.dll"); matched == false {
		t.Fatalf("'file.dll' did not match directory: '/file.dll'")
	}
}

func Test_SimularFileInRoot_FileInRootIncluded(t *testing.T) {
	ign := NewIgnore("")
	ign.ReadLine("file.dll")

	if matched, _ := ign.MatchFile("/afile.dll"); matched == true {
		t.Fatalf("'file.dll' matched directory: '/afile.dll'")
	}
}

func Test_NegateWildcardAllFiles_FileInRootIncluded(t *testing.T) {
	ign := NewIgnore("")
	ign.ReadLine("*.*")
	ign.ReadLine("!*.dll")

	if matched, _ := ign.MatchFile("/file.dll"); matched == false {
		t.Fatalf("'file.dll' didn't match directory: '/file.dll'")
	}
	if _, ignore := ign.MatchFile("/file.dll"); ignore == true {
		t.Fatalf("'file.dll' ignored file: '/file.dll'")
	}
}
