package project

import (
	"ele/utils/fileio"
	"testing"
)

func Test_SearchMigration_ReturnAllMigrations(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)

	collection := *project.Search(&SearchOptions{Search: "*"})
	if len(collection) != 2 {
		t.Fatalf("Did not find exact migrations, found: %d", len(collection))
	}
}

func Test_SearchStagedMigration_ReturnUnStagedMigrations(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")

	mr.
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("createshops.sql", "CREATE TABLE Shops").
		AddFile("createshops.rollback.sql", "INSERT INTO SHOPS")

	project := createTestProject(&mr)

	project.Stage("*Customers*")

	collection := *project.Search(&SearchOptions{Search: "*", ExcludeStaged: true})
	if len(collection) != 2 {
		t.Fatalf("Did not find migrations exluding staged, found: %d", len(collection))
	}
}
