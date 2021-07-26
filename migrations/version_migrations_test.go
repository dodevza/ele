package migrations

import "testing"

func Test_OneMigration_SingleMigration(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.0.0"})

	sequence := list.Forward("", "").Sprint()

	if sequence != "Customers" {
		t.Errorf("No migration sequence:%s", sequence)
	}
}

func Test_DependentSameVersion_InOrder(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.0.0", Requires: []string{"Shops"}})
	list.Add(Migration{Path: "Shops", Version: "1.0.0"})

	sequence := list.Forward("", "").Sprint()

	if sequence != "Shops->Customers" {
		t.Errorf("Out of sequence:%s", sequence)
	}
}

func Test_IndependentDiffVersion_InOrder(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.1.0"})
	list.Add(Migration{Path: "Shops", Version: "1.0.0"})

	sequence := list.Forward("", "").Sprint()

	if sequence != "Shops->Customers" {
		t.Errorf("Out of sequence:%s", sequence)
	}
}
