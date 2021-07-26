package loaders

import (
	"ele/migrations"
	"ele/utils/fileio"
	"testing"
)

func createSequenceWithHooks(mr *fileio.MapDirectoryReader) migrations.MigrationCollection {
	return Builder(mr).All()
}

func Test_DirectoryWithAfterHooksAtFileLevel_AfterHooksRunLast(t *testing.T) {

	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("create.sql", "ALTER TABLE Address").
		AddFile("create.after.sql", "Initialise Table").
		Back()

	sequence := createSequenceWithHooks(&mr).SprintFiles()

	if sequence != "create.sql->create.after.sql" {
		t.Errorf("After hook wan't run last, sequence: %s", sequence)
	}
}
