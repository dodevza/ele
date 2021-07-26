package loaders

import (
	"ele/migrations"
	"ele/utils/fileio"
	"strings"
	"testing"
)

func createSequenceOnlyInclusions(mr *fileio.MapDirectoryReader, included []string) migrations.MigrationCollection {

	return Builder(mr).Search(included...).All()
}

func Test_IncludeContainsFileWildcards_OnlyHaveIncludedFiles(t *testing.T) {

	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		AddFile("CreateAddress.rollback.sql", "ALTER TABLE Address").
		Back().
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		AddFile("CreateCustomer.rollback.sql", "DROP TABLE Customers").
		Back().
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		AddFile("CreateShops.rollback.sql", "DROP TABLE Shops").
		Back()

	sequence := createSequenceOnlyInclusions(&mr, []string{"*.rollback.sql"})

	anyOther := false

	for _, mig := range sequence {
		if strings.HasSuffix(mig.FileName, "rollback.sql") == false {
			anyOther = true
			break
		}
	}

	if anyOther {
		t.Errorf("Non included items found in migrations found")
	}
}

func Test_IncluedContainsFolderWildcards_OnlyContainsFilesForThatFolder(t *testing.T) {

	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		AddFile("CreateAddress.rollback.sql", "ALTER TABLE Address").
		Back().
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		AddFile("CreateCustomer.rollback.sql", "DROP TABLE Customers").
		Back().
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		AddFile("CreateShops.rollback.sql", "DROP TABLE Shops").
		Back()

	sequence := createSequenceOnlyInclusions(&mr, []string{"*Customers*"})

	anyOther := false

	example := ""
	for _, mig := range sequence {
		if strings.Contains(mig.Path, "Customers") == false {
			anyOther = true
			example = mig.Path
			break
		}
	}

	if anyOther {
		t.Errorf("Non included items found in migrations found: example: %s", example)
	}
}
