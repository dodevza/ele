package migrations

import "testing"

func Test_AddOneMigration_SingleOutput(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Customers"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Customers" {
		t.Errorf("Did not print out single output: %s", sequence)
	}
}

func Test_AddTwoMigration_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Customers"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Customers" &&
		sequence != "Customers->Shops" {
		t.Errorf("Did not print both out: %s", sequence)
	}
}

func Test_AddTwoDependentInOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Customers" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddTwoDependentReverseOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Shops"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Customers" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddThreeDependent2On1InOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Products", Requires: []string{"Shops"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Products->Customers" &&
		sequence != "Shops->Customers->Products" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddThreeDependent2On1ReverseOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Products", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Shops"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Products->Customers" &&
		sequence != "Shops->Customers->Products" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddThreeDependent1On2InOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Customers"})
	tree.Add(Migration{Path: "Cart", Requires: []string{"Shops", "Customers"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "Customers->Shops->Cart" &&
		sequence != "Shops->Customers->Cart" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddTreeDependent1On2ReverseOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Cart", Requires: []string{"Shops", "Customers"}})
	tree.Add(Migration{Path: "Customers"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Customers->Cart" &&
		sequence != "Customers->Shops->Cart" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddTreeDependentDuplicateMigration_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Cart", Requires: []string{"Shops", "Customers"}})
	tree.Add(Migration{Path: "Customers"})
	tree.Add(Migration{Path: "Cart"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Customers->Shops->Cart->Cart" &&
		sequence != "Shops->Customers->Cart->Cart" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddTreeDependentDifferentRequiresMigration_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Cart", Requires: []string{"Customers"}})
	tree.Add(Migration{Path: "Customers"})
	tree.Add(Migration{Path: "Cart", Requires: []string{"Shops"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "Customers->Shops->Cart->Cart" &&
		sequence != "Shops->Customers->Cart->Cart" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddSimpleNested_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "CustomerAddresses", Requires: []string{"Customers"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Customers->CustomerAddresses" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddNested_TwoOutcomesOrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Shops"})
	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Address"})
	tree.Add(Migration{Path: "CustomerAddresses", Requires: []string{"Customers", "Address"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Address->Customers->CustomerAddresses" &&
		sequence != "Address->Shops->Customers->CustomerAddresses" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_AddSimpleNestedRecerseOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "CustomerAddresses", Requires: []string{"Customers"}})
	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Shops"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Customers->CustomerAddresses" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_Add5LayersSimpleNestedReversOrder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "SupportedAddressTypes", Requires: []string{"AddressTypes"}})
	tree.Add(Migration{Path: "AddressTypes", Requires: []string{"CustomerAddresses"}})
	tree.Add(Migration{Path: "CustomerAddresses", Requires: []string{"Customers"}})
	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Shops"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->Customers->CustomerAddresses->AddressTypes->SupportedAddressTypes" {
		t.Errorf("Sequence: %s", sequence)
	}
}

// There is an issue where Customers was showing before customer addresses
// Added sorting to the modules, will need to figure out the correct sequence items that could
// cause that particular issue
func Test_Add5Layers4SortVariances_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "AddressType"})
	tree.Add(Migration{Path: "Address", Requires: []string{"AddressType"}})
	tree.Add(Migration{Path: "CustomerAddresses", Requires: []string{"Customers", "Address"}})
	tree.Add(Migration{Path: "Customers", Requires: []string{"Shops"}})
	tree.Add(Migration{Path: "Shops"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Shops->AddressType->Address->Customers->CustomerAddresses" &&
		sequence != "Shops->AddressType->Customers->Address->CustomerAddresses" &&
		sequence != "Shops->Customers->AddressType->Address->CustomerAddresses" &&
		sequence != "AddressType->Shops->Customers->Address->CustomerAddresses" &&
		sequence != "AddressType->Shops->Address->Customers->CustomerAddresses" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_UnSortedModules_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Customers/Addresses"})
	tree.Add(Migration{Path: "Customers"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Customers/Addresses->Customers" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_SortedModules_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "Customers"})
	tree.Add(Migration{Path: "Customers/Addresses"})

	sequence := tree.InOrder().Sprint()

	if sequence != "Customers/Addresses->Customers" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_NestedFolder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "products"})
	tree.Add(Migration{Path: "products/categories"})
	tree.Add(Migration{Path: "products/types"})

	sequence := tree.InOrder().Sprint()

	if sequence != "products/types->products/categories->products" &&
		sequence != "products/categories->products/types->products" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_NestedFolderRequiredSubFolder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "products"})
	tree.Add(Migration{Path: "products/categories", Requires: []string{"products/types"}})
	tree.Add(Migration{Path: "products/types"})

	sequence := tree.InOrder().Sprint()

	if sequence != "products/types->products/categories->products" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_NestedFolderRequiredExternalFolder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "orgs"})
	tree.Add(Migration{Path: "products"})
	tree.Add(Migration{Path: "products/categories", Requires: []string{"products/types"}})
	tree.Add(Migration{Path: "products/types", Requires: []string{"orgs"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "orgs->products/types->products/categories->products" {
		t.Errorf("Sequence: %s", sequence)
	}
}

func Test_NestedFolderRequiredExternalParentFolder_OrderedCorrectly(t *testing.T) {
	tree := EmptyDependencyTree()

	tree.Add(Migration{Path: "orgs"})
	tree.Add(Migration{Path: "products", Requires: []string{"orgs"}})

	tree.Add(Migration{Path: "products/types"})
	tree.Add(Migration{Path: "products/categories", Requires: []string{"products/types"}})

	sequence := tree.InOrder().Sprint()

	if sequence != "orgs->products/types->products/categories->products" {
		t.Errorf("Sequence: %s", sequence)
	}
}
