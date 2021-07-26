package migrations

import (
	"testing"
)

func Test_CollectionReverse_ReversedCollection(t *testing.T) {
	m1 := Migration{FileName: "F1"}
	m2 := Migration{FileName: "F2"}
	collection := &MigrationCollection{&m1, &m2}

	reverse := collection.Reverse()

	if reverse[0].FileName != "F2" || reverse[1].FileName != "F1" {
		t.Fatalf("Did not reverse collection")
	}
}

func Test_CollectionSprintV_DisplaysVersions(t *testing.T) {
	m1 := Migration{Path: "Customers", Version: "V1", FileName: "F1"}
	m2 := Migration{Path: "Shops", Version: "V2", FileName: "F2"}
	collection := &MigrationCollection{&m1, &m2}

	sequence := collection.SprintV()

	if sequence != "V1->V2" {
		t.Fatalf("Sequence only display versions, sequence: %s", sequence)
	}
}

func Test_CollectionSprintMV_DisplaysVersionsAndModules(t *testing.T) {
	m1 := Migration{Path: "Customers", Version: "V1", FileName: "F1"}
	m2 := Migration{Path: "Shops", Version: "V2", FileName: "F2"}
	collection := &MigrationCollection{&m1, &m2}

	sequence := collection.SprintMV()

	if sequence != "Customers@V1->Shops@V2" {
		t.Fatalf("Sequence only display versions and modules, sequence: %s", sequence)
	}
}

func Test_CollectionSprintFiles_DisplayFiles(t *testing.T) {
	m1 := Migration{Path: "Customers", Version: "V1", FileName: "F1"}
	m2 := Migration{Path: "Shops", Version: "V2", FileName: "F2"}
	collection := &MigrationCollection{&m1, &m2}

	sequence := collection.SprintFiles()

	if sequence != "F1->F2" {
		t.Fatalf("Sequence only display files, sequence: %s", sequence)
	}
}

func Test_CollectionPlan_DisplayPlan(t *testing.T) {
	m1 := Migration{Path: "Customers", Version: "V1", FileName: "F1"}
	m2 := Migration{Path: "Shops", Version: "V2", FileName: "F2"}
	collection := MigrationCollection{&m1, &m2}

	collection.PrintExcecutionPlan()
}
