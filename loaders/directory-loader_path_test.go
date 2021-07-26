package loaders

import (
	"ele/migrations"
	"ele/utils/fileio"
	"testing"
)

func createSequenceWithPath(mr *fileio.MapDirectoryReader, workingDir *fileio.MapDirectoryReader) migrations.MigrationCollection {
	return Builder(mr).WorkingDirectory(workingDir).All()
}

func Test_PathProvided_OnlyPathMigrationsLoaded(t *testing.T) {
	mr := fileio.NewMapReader()
	customers := mr.AddDirectory("Customers")
	customers.
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	mr.
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		Back()

	sequence := createSequenceWithPath(&mr, customers)
	check := sequence.Sprint()

	if check != "Customers" {
		t.Errorf("Only one version sequence:%s", check)
	}
}

func Test_PathProvided_NoMigrationsLoadedUntilPathBegins(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.AddFile("before.sql", "Create schema before")

	customers := mr.AddDirectory("Customers")
	customers.AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	sequence := createSequenceWithPath(&mr, customers)
	check := sequence.SprintMV()

	if check != "Customers@V1.0.0" {
		t.Errorf("Only one version sequence:%s", check)
	}
}
