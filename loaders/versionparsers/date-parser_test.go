package versionparsers

import (
	"testing"
)

func Test_Date_IsMatch(t *testing.T) {

	parser := NewDateParser("YYYY-MM-DD")

	_, isVersion := parser.Parse("2018-02-01")

	if isVersion == false {
		t.Fatalf("Not a version")
	}
}

func Test_DateMericaFormat_IsMatch(t *testing.T) {

	parser := NewDateParser("MM.DD.YYYY")

	version, isVersion := parser.Parse("02.01.2018")

	if isVersion == false {
		t.Fatalf("Not a version")
	}

	if version != "2018-02-01" {
		t.Fatalf("Did not convert to alpha numerical sortable format")
	}
}

func Test_DatePrefixFormat_IsMatch(t *testing.T) {

	parser := NewDateParser("VYYYY.DD.MM")

	version, isVersion := parser.Parse("V2018.01.02")

	if isVersion == false {
		t.Fatalf("Not a version")
	}

	if version != "2018-02-01" {
		t.Fatalf("Did not convert to alpha numerical sortable format")
	}
}
