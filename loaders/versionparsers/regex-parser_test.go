package versionparsers

import "testing"

func Test_InitRegexVersion_IsMatch(t *testing.T) {
	parser := NewRegexParser("(?i)^_init")

	_, isVersion := parser.Parse("_init")

	if isVersion == false {
		t.Fatalf("Did not detect as version")
	}
}
