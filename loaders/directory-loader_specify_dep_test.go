package loaders

import (
	"ele/utils/fileio"
	"testing"
)

func Test_DependencySpecifiedInFile_InOrder(t *testing.T) {

	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("V1.0.0").
		AddFile("CreateAddress.sql", "Create Table Address").
		Back()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("CreateCustomer.sql", "-- require: Address").
		Back()

	sequence := createSequence(&mr).Sprint()

	if sequence != "Address->Customers" {
		t.Errorf("Sequence: %s", sequence)
	}
}
