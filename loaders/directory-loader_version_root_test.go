package loaders

import (
	"ele/utils/fileio"
	"testing"
)

func Test_NoVersion_MigrationsInNestedOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("orgs").
		AddFile("CreateOrganisations.sql", "ALTER TABLE Organisations").
		Back().
		AddDirectory("products").
		AddFile("CreateProducts.sql", "-- require: orgs").
		AddDirectory("types").
		AddFile("CreateProductTypes.sql", "Create Product Table").
		Back().
		AddDirectory("categories").
		AddFile("CreateProductCategories.sql", "-- require: products/types").
		Back()

	sequence := createSequence(&mr)

	moduleOrder := sequence.SprintMV()
	if moduleOrder != "orgs@~->products/types@~->products/categories@~->products@~" {
		t.Errorf("Modules in incorrect order:%s", moduleOrder)
	}
}

func Test_NoRootVersion_MigrationsInNestedOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("orgs").
		AddDirectory("V1.0.0").
		AddFile("CreateOrganisations.sql", "ALTER TABLE Organisations").
		Back().
		Back().
		AddDirectory("products").
		AddDirectory("V1.0.0").
		AddFile("CreateProducts.sql", "-- require: orgs").
		Back().
		AddDirectory("types").
		AddDirectory("V1.0.0").
		AddFile("CreateProductTypes.sql", "Create products").
		Back()

	sequence := createSequence(&mr)

	moduleOrder := sequence.SprintMV()
	if moduleOrder != "orgs@V1.0.0->products/types@V1.0.0->products@V1.0.0" {
		t.Errorf("Modules in incorrect order:%s", moduleOrder)
	}
}

func Test_RootVersion_MigrationsInNestedOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddDirectory("orgs").
		AddFile("CreateOrganisations.sql", "ALTER TABLE Organisations").
		Back().
		AddDirectory("products").
		AddFile("CreateProducts.sql", "ALTER TABLE Products").
		AddDirectory("types").
		AddFile("CreateProductTypes.sql", "-- require: orgs").
		Back()

	sequence := createSequence(&mr)

	moduleOrder := sequence.SprintMV()
	// Indiferent about the execution order if orgs is required on sub level
	if moduleOrder != "orgs@V1.0.0->products/types@V1.0.0->products@V1.0.0" {
		t.Errorf("Modules in incorrect order:%s", moduleOrder)
	}
}

func Test_RequiresInBothSubAndModule_MigrationsInNestedOrder(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddDirectory("orgs").
		AddFile("CreateOrganisations.sql", "ALTER TABLE Organisations").
		Back().
		AddDirectory("products").
		AddFile("CreateProducts.sql", "-- require: orgs").
		AddDirectory("types").
		AddFile("CreateProductTypes.sql", "-- require: orgs").
		Back()

	sequence := createSequence(&mr)

	moduleOrder := sequence.SprintMV()
	if moduleOrder != "orgs@V1.0.0->products/types@V1.0.0->products@V1.0.0" {
		t.Errorf("Modules in incorrect order:%s", moduleOrder)
	}
}

func Test_Root_Version_Require_RequireCarriedOverToSubFolders(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddDirectory("orgs").
		AddFile("CreateOrganisations.sql", "ALTER TABLE Organisations").
		Back().
		AddDirectory("products").
		AddFile("CreateProducts.sql", "-- require: orgs").
		AddDirectory("types").
		AddFile("CreateProductTypes.sql", "CREATE TABLE ProductTypes").
		Back()

	sequence := createSequence(&mr)

	moduleOrder := sequence.SprintMV()
	if moduleOrder != "orgs@V1.0.0->products/types@V1.0.0->products@V1.0.0" {
		t.Errorf("Modules in incorrect order:%s", moduleOrder)
	}
}
