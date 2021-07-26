package loaders

import (
	"ele/utils/fileio"
	"testing"
)

func Test_NestedDep_InOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		Back().
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back()

	sequence := createSequence(&mr).Sprint()

	if sequence != "Address/Customers->Address" {
		t.Errorf("Nested ordered sequence:%s", sequence)
	}
}

func Test_TripleNestedDep_InOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		Back().
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back().
		AddDirectory("Shops").
		AddDirectory("V1.0.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		Back()

	sequence := createSequence(&mr).Sprint()

	if sequence != "Address/Customers/Shops->Address/Customers->Address" {
		t.Errorf("3 Nested ordered sequence:%s", sequence)
	}
}

func Test_NestedVarientVersionDep_InOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		Back().
		AddDirectory("Customers").
		AddDirectory("V1.1.0").
		AddFile("CreateCustomer.sql", "CREATE TABLE Customers").
		Back().
		AddDirectory("Shops").
		AddDirectory("V1.2.0").
		AddFile("CreateShops.sql", "CREATE TABLE Shops").
		Back()

	sequence := createSequence(&mr).SprintMV()

	if sequence != "Address@V1.0.0->Address/Customers@V1.1.0->Address/Customers/Shops@V1.2.0" {
		t.Errorf("3 Nested ordered sequence:%s", sequence)
	}
}
