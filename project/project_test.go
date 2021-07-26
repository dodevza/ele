package project

import (
	"ele/utils/fileio"
)

func createTestProject(mr *fileio.MapDirectoryReader) *Project {
	mr.AddFile("ele.yml", "")

	return NewProject(mr)
}
