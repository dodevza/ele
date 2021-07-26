package loaders

import (
	"ele/utils/fileio"
	"testing"
)

func Test_FolderVersioning_MigrationsHasVersions(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		Back().
		AddDirectory("Customers").
		AddDirectory("V1.0.1").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	sequence := createSequence(&mr).SprintV()

	if sequence != "V1.0.0->V1.0.1" {
		t.Errorf("Versions not found for migrations:%s", sequence)
	}
}

func Test_MultipleVersionsInHierarchy_HaveFirstVersionAsVersion(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddDirectory("Customers").
		AddDirectory("V1.0.BAD").
		AddFile("CreateAddress.sql", "ALTER TABLE Address")

	sequence := createSequence(&mr).SprintV()

	if sequence != "V1.0.0" {
		t.Errorf("First version in hierachy not used, sequence:%s", sequence)
	}
}

func Test_MigrationsInSubFolderUnderVersion_MigrationLoadedWithVersion(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddDirectory("Customers").
		AddFile("CreateAddress.sql", "ALTER TABLE Address")

	sequence := createSequence(&mr).SprintV()

	if sequence != "V1.0.0" {
		t.Errorf("Migration not found, sequence:%s", sequence)
	}
}

func Test_NoVersionSpecified_MigrationLoadedWithoutVersion(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("Customers").
		AddFile("CreateAddress.sql", "ALTER TABLE Address")

	sequence := createSequence(&mr).SprintMV()

	if sequence != "Address/Customers@~" {
		t.Errorf("Migration not found, sequence:%s", sequence)
	}
}
