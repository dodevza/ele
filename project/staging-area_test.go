package project

import (
	"ele/migrations"
	"ele/utils/fileio"
	"testing"
)

func createProjectStagingArea(mr *fileio.MapDirectoryReader) *StagingArea {
	mr.AddFile("ele.yml", "")

	return NewStagingArea(mr)
}

func Test_StageMigration_MigrationStaged(t *testing.T) {
	mr := fileio.NewMapReader()
	area := createProjectStagingArea(&mr)

	migration := migrations.Migration{Version: "1", Path: "Customers", FileName: "MyFile"}

	collection := &migrations.MigrationCollection{&migration}

	area.Stage(collection)

	migmap := area.getMigrationMap()

	if migmap[migration.GetUniqueKey()] == false {
		t.Fatalf("Did not add migration to staged")
	}
}

func Test_StageMultpleMigration_AllMigrationStaged(t *testing.T) {
	mr := fileio.NewMapReader()
	area := createProjectStagingArea(&mr)

	mig1 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MyFile"}
	mig2 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MySecond"}

	collection := migrations.MigrationCollection{&mig1, &mig2}

	area.Stage(&collection)

	migmap := area.getMigrationMap()

	if len(migmap) != 2 {
		t.Fatalf("Did not add migration to staged, 1: %t, 2: %t", migmap[mig1.GetUniqueKey()], migmap[mig2.GetUniqueKey()])
	}
}

func Test_UnstageStaged_AllMigrationUnStaged(t *testing.T) {
	mr := fileio.NewMapReader()
	area := createProjectStagingArea(&mr)

	mig1 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MyFile"}
	mig2 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MySecond"}

	collection := migrations.MigrationCollection{&mig1, &mig2}

	area.Stage(&collection)

	area.Unstage(&collection)

	migmap := area.getMigrationMap()

	count := 0
	for _, value := range migmap {
		if value {
			count++
		}
	}
	if count != 0 {
		t.Fatalf("Not all migrations did unstage")
	}
}

func Test_StageUnstaged_AllMigrationUnStaged(t *testing.T) {
	mr := fileio.NewMapReader()
	area := createProjectStagingArea(&mr)

	mig1 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MyFile"}
	mig2 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MySecond"}

	collection := migrations.MigrationCollection{&mig1, &mig2}

	area.Stage(&collection)

	area.Unstage(&collection)

	area.Stage(&collection)

	migmap := area.getMigrationMap()

	count := 0
	for _, value := range migmap {
		if value {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("Not all unstaged migrations did stage again")
	}
}

func Test_UndoStaged_AllMigrationUnStaged(t *testing.T) {
	mr := fileio.NewMapReader()
	area := createProjectStagingArea(&mr)

	mig1 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MyFile"}
	mig2 := migrations.Migration{Version: "1", Path: "Customers", FileName: "MySecond"}

	collection := migrations.MigrationCollection{&mig1, &mig2}

	area.Stage(&collection)

	area.Undo()

	migmap := area.getMigrationMap()

	if len(migmap) != 0 {
		t.Fatalf("Not all migrations did undo")
	}
}
