package loaders

import (
	"ele/migrations"
	"ele/utils/fileio"
	"testing"
)

func createSequence(mr *fileio.MapDirectoryReader) migrations.MigrationCollection {

	return Builder(mr).All()
}

func Test_OnlyOneVersion_Loaded(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		AddFile("UpdateCustomer.sql", "ALTER TABLE Customers").
		Back()

	sequence := createSequence(&mr).SprintV()

	if sequence != "V1.0.0->V1.0.0" {
		t.Errorf("Only one version sequence:%s", sequence)
	}
}

func Test_MultipleVersions_InOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.1.0").
		AddFile("UpdateCustomer.sql", "ALTER TABLE Customers").
		Back().
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	sequence := createSequence(&mr).SprintV()

	if sequence != "V1.0.0->V1.1.0" {
		t.Errorf("In Order sequence:%s", sequence)
	}
}

func Test_SingleModule_Loaded(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	sequence := createSequence(&mr).Sprint()

	if sequence != "Customers" {
		t.Errorf("Loaded sequence:%s", sequence)
	}
}

func Test_MultipeModules_InOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	mr.
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		Back()

	sequence := createSequence(&mr).Sprint()

	if sequence != "Customers->Shops" {
		t.Errorf("Loaded sequence:%s", sequence)
	}
}

func Test_MultipeModulesDiffVersions_InOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.1").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	mr.
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		Back()

	sequence := createSequence(&mr).Sprint()

	if sequence != "Shops->Customers" {
		t.Errorf("Loaded sequence:%s", sequence)
	}
}

func Test_DuplicateModulesDiffVersions_InOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.1").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	mr.
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		Back().
		AddDirectory("V1.1.0").
		AddFile("UpdateShops.sql", "ALTER TABLE Shops").
		Back()

	sequence := createSequence(&mr).SprintMV()

	if sequence != "Shops@V1.0.0->Customers@V1.0.1->Shops@V1.1.0" {
		t.Errorf("InOrder sequence:%s", sequence)
	}

}
