package loaders

import (
	"ele/utils/fileio"
	"testing"
)

func Test_WellContructed_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner("- require: Customers")
	requirements := scanner.Scan(textScanner)

	if len(requirements) != 1 ||
		requirements[0] != "Customers" {
		t.Errorf("Incorrect requirements: %s", requirements)
	}
}

func Test_MultipleRequirements_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner("- require: Customers, Addresses")
	requirements := scanner.Scan(textScanner)

	if len(requirements) != 2 ||
		requirements[0] != "Customers" ||
		requirements[1] != "Addresses" {
		t.Errorf("Incorrect requirements: %s", requirements)
	}
}

func Test_Multiline_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner(
		`
		- This is a multiline File
		- File will be scanned until code starts (No header comments)

		- require: Customers

		CREATE TABLE Address (
			Id "\"")INT
		)
	`)
	requirements := scanner.Scan(textScanner)

	if len(requirements) != 1 ||
		requirements[0] != "Customers" {
		t.Errorf("Incorrect requirements: %s", requirements)
	}
}

func Test_MultiMultipleRequirements_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner(
		`
		- This is a multiline File
		- File will be scanned until code starts (No header comments)

		- require: Customers
		- require: Addresses

		CREATE TABLE Address (
			Id INT
		)
	`)
	requirements := scanner.Scan(textScanner)

	if len(requirements) != 2 ||
		requirements[0] != "Customers" ||
		requirements[1] != "Addresses" {
		t.Errorf("Incorrect requirements: %s", requirements)
	}
}

func Test_IgnoreAfterCodeStarted_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner(
		`
		- This is a multiline File
		- File will be scanned until code starts (No header comments)


		CREATE TABLE Address (
			Id INT
		)

		- require: Customers
	`)
	requirements := scanner.Scan(textScanner)

	if len(requirements) == 0 {
		t.Errorf("Requirements found")
	}
}

func Test_LooselyContructed_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner("- require Customers Addresses")
	requirements := scanner.Scan(textScanner)

	if len(requirements) != 2 ||
		requirements[0] != "Customers" ||
		requirements[1] != "Addresses" {
		t.Errorf("Incorrect requirements: %s", requirements)
	}
}

func Test_WhitespacesContructed_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner("\t\t-require\t Customers\t\t Addresses")
	requirements := scanner.Scan(textScanner)

	if len(requirements) != 2 ||
		requirements[0] != "Customers" ||
		requirements[1] != "Addresses" {
		t.Errorf("Incorrect requirements: %s", requirements)
	}
}

func Test_Quotes_AllRequirements(t *testing.T) {
	scanner := createDependencyScanner(".sql")
	textScanner := fileio.NewTextScanner("-require \"Customer Addresses\"")
	requirements := scanner.Scan(textScanner)

	if len(requirements) != 1 ||
		requirements[0] != "Customer Addresses" {
		t.Errorf("Incorrect requirements: %s", requirements)
	}
}
