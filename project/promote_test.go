package project

import (
	"ele/utils/fileio"
	"testing"
)

func Test_PromoteStaged_MigrationsMovedToTargetFolders(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)

	project.Stage("*")

	project.Promote(&PromoteOptions{Target: "V1.0.1"})

	movedMirgations := project.Search(&SearchOptions{Search: "*"}).SprintV()

	if movedMirgations != "V1.0.1->V1.0.1" {
		t.Fatalf("Did not move staged migrations, versions: %s", movedMirgations)
	}
}

func Test_VersionNestedModules_MigrationsMovedToTargetFolders(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddDirectory("Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS").
		Back().
		AddDirectory("Shops").
		AddFile("createshops.rollback.sql", "INSERT INTO SHOPS")
	project := createTestProject(&mr)

	project.Stage("*")

	project.Promote(&PromoteOptions{Target: "V1.0.1"})

	movedMirgations := project.Search(&SearchOptions{Search: "*"}).SprintV()

	if movedMirgations != "V1.0.1->V1.0.1" {
		t.Fatalf("Did not move staged migrations, versions: %s", movedMirgations)
	}
}

func Test_NoVersionNestedModules_NewTargetFolderCreatedWithNestedMigrations(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS").
		Back().
		AddDirectory("Shops").
		AddFile("createshops.rollback.sql", "INSERT INTO SHOPS")
	project := createTestProject(&mr)

	project.Stage("*")

	project.Promote(&PromoteOptions{Target: "V1.0.1"})

	movedMirgations := project.Search(&SearchOptions{Search: "*"}).SprintV()

	if movedMirgations != "V1.0.1->V1.0.1" {
		t.Fatalf("Did not move staged migrations, versions: %s", movedMirgations)
	}
}

func Test_PromoteUnStaged_MigrationsMovedToTargetFolders(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)
	project.ActivateTags("V1.0.0", "V1.0.0")

	project.Promote(&PromoteOptions{Target: "V1.0.1"})

	movedMirgations := project.Search(&SearchOptions{Search: "*"}).SprintV()

	if movedMirgations != "V1.0.1->V1.0.1" {
		t.Fatalf("Did not move staged migrations, versions: %s", movedMirgations)
	}
}

func Test_PromoteVersionThatDoentExist_MigrationsStayAsItWas(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)
	project.ActivateTags("V1.0.1", "V1.0.1")

	project.Promote(&PromoteOptions{Target: "V2.0.1"})

	movedMirgations := project.Search(&SearchOptions{Search: "*"}).SprintV()

	if movedMirgations != "V1.0.0->V1.0.0" {
		t.Fatalf("Migrations did move, versions: %s", movedMirgations)
	}
}

func Test_PromoteNonVersionedUnStaged_NonVersionMigrationsMoved(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Customers").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)

	project.Promote(&PromoteOptions{Target: "V2.0.1"})

	movedMirgations := project.Search(&SearchOptions{Search: "*"}).SprintV()

	if movedMirgations != "V2.0.1->V2.0.1" {
		t.Fatalf("Migrations did not move, versions: %s", movedMirgations)
	}
}

func Test_PromoteNonVersionedUnStagedInRoot_NonVersionMigrationsMoved(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)

	project.Promote(&PromoteOptions{Target: "V2.0.1"})

	movedMirgations := project.Search(&SearchOptions{Search: "*"}).SprintV()

	if movedMirgations != "V2.0.1->V2.0.1" {
		t.Fatalf("Migrations did not move, versions: %s", movedMirgations)
	}
}

func Test_PromoteNonVersionedUnStagedInRoot_RemovedV1Directory(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)
	project.ActivateTags("V1.0.0", "V1.0.0")

	project.Promote(&PromoteOptions{Target: "V2.0.1"})

	foundOldVersion := false
	foundNewVersion := false
	for _, d := range project.dir.SubDirectories() {
		if d == "V1.0.0" {
			foundOldVersion = true
		}

		if d == "V2.0.1" {
			foundNewVersion = true
			break
		}
	}
	if foundOldVersion {
		t.Fatalf("Did not remove previous version dir")
	}

	if !foundNewVersion {
		t.Fatalf("New version folder wasn't created")
	}
}

func Test_PromoteInvalidTarget_ReturnError(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)
	project.ActivateTags("V1.0.0", "V1.0.0")
	error := project.Promote(&PromoteOptions{Target: "ANINVALIDTAG"})

	if error == nil {
		t.Fatalf("No error provided")
	}
}

func Test_PromoteWithoutTarget_ReturnError(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)
	project.ActivateTags("V1.0.0", "V1.0.0")
	error := project.Promote(&PromoteOptions{})

	if error == nil {
		t.Fatalf("No error provided")
	}
}

func Test_PromoteWithoutValidTargetWithCreateTarget_ReturnError(t *testing.T) {
	mr := fileio.NewMapReader()
	mr.
		AddDirectory("V1.0.0").
		AddFile("createcustomers.sql", "CREATE TABLE Customers").
		AddFile("createcustomers.rollback.sql", "INSERT INTO CUSTOMERS")
	project := createTestProject(&mr)
	project.ActivateTags("V1.0.0", "V1.0.0")

	error := project.Promote(&PromoteOptions{Target: "ANINVALIDTAG", CreateTag: true})

	if error != nil {
		t.Fatalf("No error exepected, error: %s", error)
	}

	if project.IsValidTag("ANINVALIDTAG") == false {
		t.Fatalf("Target wasn't added to project")
	}
}
