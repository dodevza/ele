package migrations

import "testing"

func Test_ExistingVersionRange_InVersionRange(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.1.0"})
	list.Add(Migration{Path: "Shops", Version: "1.0.0"})
	list.Add(Migration{Path: "CustomerShops", Version: "1.2.0"})
	list.Add(Migration{Path: "Addresses", Version: "1.4.0"})
	list.Add(Migration{Path: "CustomerAddresses", Version: "1.3.0"})

	sequence := list.Forward("1.2.0", "1.3.0").Sprint()

	if sequence != "CustomerShops->CustomerAddresses" {
		t.Errorf("In version range sequence:%s", sequence)
	}
}

func Test_SlighlyHigherVersionRange_InVersionRange(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.1.0"})
	list.Add(Migration{Path: "Shops", Version: "1.0.0"})
	list.Add(Migration{Path: "CustomerShops", Version: "1.2.0"})
	list.Add(Migration{Path: "Addresses", Version: "1.4.0"})
	list.Add(Migration{Path: "CustomerAddresses", Version: "1.3.0"})

	sequence := list.Forward("1.2.0", "1.3.1").Sprint()

	if sequence != "CustomerShops->CustomerAddresses" {
		t.Errorf("In version range sequence:%s", sequence)
	}
}

func Test_SlighlyLowerVersionRange_InVersionRange(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.1.0"})
	list.Add(Migration{Path: "Shops", Version: "1.0.0"})
	list.Add(Migration{Path: "CustomerShops", Version: "1.2.0"})
	list.Add(Migration{Path: "Addresses", Version: "1.4.0"})
	list.Add(Migration{Path: "CustomerAddresses", Version: "1.3.0"})

	sequence := list.Forward("1.1.9", "1.3.0").Sprint()

	if sequence != "CustomerShops->CustomerAddresses" {
		t.Errorf("In version range sequence:%s", sequence)
	}
}

func Test_OpenEndVersionRange_InVersionRange(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.1.0"})
	list.Add(Migration{Path: "Shops", Version: "1.0.0"})
	list.Add(Migration{Path: "CustomerShops", Version: "1.2.0"})
	list.Add(Migration{Path: "Addresses", Version: "1.4.0"})
	list.Add(Migration{Path: "CustomerAddresses", Version: "1.3.0"})

	sequence := list.Forward("1.2.0", "").Sprint()

	if sequence != "CustomerShops->CustomerAddresses->Addresses" {
		t.Errorf("In version range sequence:%s", sequence)
	}
}

func Test_OpenStartVersionRange_InVersionRange(t *testing.T) {
	list := EmptyVersionTree()

	list.Add(Migration{Path: "Customers", Version: "1.1.0"})
	list.Add(Migration{Path: "Shops", Version: "1.0.0"})
	list.Add(Migration{Path: "CustomerShops", Version: "1.2.0"})
	list.Add(Migration{Path: "Addresses", Version: "1.4.0"})
	list.Add(Migration{Path: "CustomerAddresses", Version: "1.3.0"})

	sequence := list.Forward("", "1.2.0").Sprint()

	if sequence != "Shops->Customers->CustomerShops" {
		t.Errorf("In version range sequence:%s", sequence)
	}
}
