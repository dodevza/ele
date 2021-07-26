package loaders

import (
	"ele/utils/fileio"
	"testing"
)

func Test_WildcardAllFiles_IgnoreAllFiles(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddFile(".eleignore", "*.*").
		AddFile("package.json", "{ \"Name\": \"My Package\" }").
		AddFile("package.sql", "Some sql script in the root").
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("help.md", "### Help").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		AddFile("CreateAddress.rollback.sql", "ALTER TABLE Address").
		Back()

	sequence := createSequence(&mr)

	if len(sequence) != 0 {
		t.Errorf("Not all the files was ignored")
	}
}

func Test_AnyV1Direcotries_IgnoreOnlyV1Files(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddFile(".eleignore", "V1.0.0/").
		AddFile("package.json", "{ \"Name\": \"My Package\" }").
		AddFile("package.sql", "Some sql script in the root").
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("help.md", "### Help").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		AddFile("CreateAddress.rollback.sql", "ALTER TABLE Address").
		Back()

	sequence := createSequence(&mr)

	if len(sequence) != 2 {
		t.Errorf("Not all the files was ignored, files: %s", sequence.SprintFiles())
	}
}

func Test_IngoreFiles_IgnorefilesShouldAutomaticallyBeExcluded(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddFile(".eleignore", "*.sql").
		AddDirectory("Address").
		AddFile(".eleignore", "*.sql").
		AddFile(".gitignore", "*.sql")

	sequence := createSequence(&mr)

	if len(sequence) != 0 {
		t.Errorf("Not all the files was ignored, files: %s", sequence.SprintFiles())
	}
}

func Test_NestedIgnoreFile_IgnoreOnlyFromIgnoreFile(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddFile("package.json", "{ \"Name\": \"My Package\" }").
		AddFile("package.sql", "Some sql script in the root").
		AddDirectory("Address").
		AddFile(".eleignore", "*.sql").
		AddDirectory("V1.0.0").
		AddFile("help.md", "### Help").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		AddFile("CreateAddress.rollback.sql", "ALTER TABLE Address").
		Back().Back().
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("help.md", "### Help").
		AddFile("CreateCustomers.sql", "ALTER TABLE Customers").
		AddFile("CreateCustomers.rollback.sql", "ALTER TABLE Customers").
		Back().Back()

	sequence := createSequence(&mr)

	if len(sequence) != 6 {
		t.Errorf("Not all the files was ignored, files: %s", sequence.SprintFiles())
	}
}

func Test_NegateAllFilesOnlySqlFiles_DontIngoreSqlFiles(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddFile("package.json", "{ \"Name\": \"My Package\" }").
		AddFile("package.sql", "Some sql script in the root").
		AddFile(".eleignore", "*.*\n!*.sql").
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("help.md", "### Help").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		AddFile("CreateAddress.rollback.sql", "ALTER TABLE Address").
		Back().Back().
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("help.md", "### Help").
		AddFile("CreateCustomers.sql", "ALTER TABLE Customers").
		AddFile("CreateCustomers.rollback.sql", "ALTER TABLE Customers").
		Back().Back()

	sequence := createSequence(&mr)

	if len(sequence) != 5 {
		t.Errorf("Sql files count missmatched, files (%d): %s", len(sequence), sequence.SprintFiles())
	}
}
