package loaders

import (
	"ele/migrations"
	"ele/utils/fileio"
	"strings"
	"testing"
)

func createSequenceWithExclusions(mr *fileio.MapDirectoryReader, excluded []string) migrations.MigrationCollection {
	return Builder(mr).Search("*").ExcludedFromSearch(excluded...).All()
}

func Test_ExcludeContainsFileWildcards_RemoveFilesFromTree(t *testing.T) {

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

	sequence := createSequenceWithExclusions(&mr, []string{"*.rollback.*"})

	anyRollbacks := false

	for _, mig := range sequence {
		if strings.HasSuffix(mig.FileName, "rollback.sql") {
			anyRollbacks = true
			break
		}
	}

	if anyRollbacks {
		t.Errorf("Excluded items found in migrations found")
	}
}

func Test_ExcludeContainsFolderWildcards_RemoveFolderFromTree(t *testing.T) {

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

	sequence := createSequenceWithExclusions(&mr, []string{"Customers/"})

	anyCustomers := false

	example := ""
	for _, mig := range sequence {
		if strings.HasSuffix(mig.Path, "Customers") {
			anyCustomers = true
			example = mig.Path
			break
		}
	}

	if len(sequence) == 0 || anyCustomers {
		t.Errorf("Excluded items found in migrations found: example: %s", example)
	}
}
