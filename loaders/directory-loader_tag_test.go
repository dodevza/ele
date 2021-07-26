package loaders

import (
	"ele/migrations"
	"ele/utils/fileio"
	"testing"
)

func createSequenceWithTags(mr *fileio.MapDirectoryReader, tags []string) migrations.MigrationCollection {

	return Builder(mr).SetNonVersionTags(tags...).All()
}

func Test_TagsDefined_TreeContainsTagAsVersion(t *testing.T) {

	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("UAT").
		AddFile("CreateAddress.sql", "ALTER TABLE Address")

	sequence := createSequenceWithTags(&mr, []string{"Uat"}).SprintV()

	if sequence != "UAT" {
		t.Errorf("Tag was not used")
	}
}

func Test_TagsAndVersionsDefined_TagsShouldRunAtTheEnd(t *testing.T) {

	mr := fileio.NewMapReader()
	mr.
		AddDirectory("Address").
		AddDirectory("UAT").
		AddFile("CreateAddress.sql", "ALTER TABLE Address").
		Back().
		AddDirectory("V1.0.0").
		AddFile("CreateFile.sql", "CREATE TABLE File").
		Back()

	sequence := createSequenceWithTags(&mr, []string{"Uat"}).SprintV()

	if sequence != "V1.0.0->UAT" {
		t.Errorf("Tag wasn't ordered after version: %s", sequence)
	}
}
