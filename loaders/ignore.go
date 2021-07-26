package loaders

import (
	"ele/plugins/git"
	"ele/utils/fileio"
	"strings"
)

func chainIgnore(ignoreFiles []string, previous []*git.Ignore, dir fileio.DirectoryReader, relativepath string) []*git.Ignore {
	var ignore *git.Ignore
	for _, igf := range ignoreFiles {
		ignoreFile, err := dir.FileScanner(igf)

		if err == nil {
			if ignore == nil {
				ignore = git.NewIgnore(relativepath)
			}

			for ignoreFile.Scan() {
				ignore.ReadLine(ignoreFile.Text())
			}
			ignoreFile.Close()
		}
	}

	if ignore != nil {
		previous = append(previous, ignore)
	}
	return previous

}

func shouldIgnoreDirectory(chain []*git.Ignore, path string) bool {
	ignoreDir := false
	for _, ign := range chain {
		match, ignore := ign.MatchDirectory(path)
		if match {
			ignoreDir = ignore
		}
	}

	return ignoreDir
}

func shouldIgnoreFile(chain []*git.Ignore, path string) bool {
	if strings.HasSuffix(path, ".gitignore") || strings.HasSuffix(path, ".eleignore") {
		return true
	}
	ignoreFile := false
	for _, ign := range chain {
		match, ignore := ign.MatchFile(path)
		if match {
			ignoreFile = ignore
		}
	}

	return ignoreFile
}
