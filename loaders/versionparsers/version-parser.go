package versionparsers

// VersionParser ...
type VersionParser interface {
	Parse(text string) (string, bool)
}
