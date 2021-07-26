package migrations

import (
	"ele/constants"
	"ele/utils/doc"
	"ele/utils/fileio"
	"fmt"
	"strings"
)

// Migration ...
type Migration struct {
	Version      string
	Path         string
	RelativePath string
	FileName     string

	Requires []string

	VersionDir      fileio.DirectoryReader
	PathFromVersion string
	Dir             fileio.DirectoryReader
}

// GetUniqueKey ...
func (mgr *Migration) GetUniqueKey() string {
	return fmt.Sprintf("%s|%s|%s", mgr.Version, mgr.Path, mgr.FileName)
}

// MigrationCollection ...
type MigrationCollection []*Migration

// Reverse ...
func (collection MigrationCollection) Reverse() MigrationCollection {
	result := make(MigrationCollection, 0)

	for i := len(collection) - 1; i >= 0; i-- {
		result = append(result, collection[i])
	}
	return result
}

// Sprint ...
func (collection MigrationCollection) Sprint() string {
	result := ""
	for _, migration := range collection {
		result += fmt.Sprintf("%s->", migration.Path)
	}

	return strings.TrimRight(result, "->")
}

// SprintV Display Versions
func (collection MigrationCollection) SprintV() string {
	result := ""
	for _, migration := range collection {
		result += fmt.Sprintf("%s->", migration.Version)
	}

	return strings.TrimRight(result, "->")
}

// SprintMV Display Module@Versions
func (collection MigrationCollection) SprintMV() string {
	result := ""
	for _, migration := range collection {
		result += fmt.Sprintf("%s@%s->", migration.Path, migration.Version)
	}

	return strings.TrimRight(result, "->")
}

// SprintFiles Display File
func (collection MigrationCollection) SprintFiles() string {
	result := ""
	for _, migration := range collection {
		result += fmt.Sprintf("%s->", migration.FileName)
	}

	return strings.TrimRight(result, "->")
}

// PrintExcecutionPlan Output execution plan
func (collection MigrationCollection) PrintExcecutionPlan() {
	lastVersion := ""
	lastPath := ""
	tabPrefix := ""
	for _, migration := range collection {
		if lastVersion != migration.Version {
			lastVersion = migration.Version
			lastPath = ""
			tabPrefix = ""
			fmt.Printf("%s\n", migration.Version)
		}
		if lastPath != migration.Path {
			lastPath = migration.Path
			tabPrefix += "\t"
			fmt.Printf("%s%s\n", tabPrefix, migration.Path)
		}

		fmt.Printf("%s%s\n", tabPrefix, migration.FileName)
	}
}

// PrintList ...
func (collection MigrationCollection) PrintList() int {
	var versionRows *doc.Rows = nil
	var repeatRows *doc.Rows = nil

	for _, migration := range collection {

		if migration.Version != constants.REPEATABLE {
			if versionRows == nil {
				versionRows = doc.NewTable().
					Column("Version", 15).Column("Module", 30).Column("File", 35).
					StartRows()
			}

		} else {
			if versionRows != nil {
				versionRows.Divider()
				versionRows = nil
			}
			if repeatRows == nil {
				doc.Line(doc.Title("Repeatable Migrations"), doc.Info("baseline"), doc.Hint("'ele promote <tag name>'"))
				doc.NewLine()
				repeatRows = doc.NewTable().
					Column("Module", 45).Column("File", 35).
					StartRows()
			}

		}

		if versionRows != nil {
			versionRows.Row(migration.Version, migration.Path, migration.FileName)
		}

		if repeatRows != nil {
			repeatRows.Row(migration.Path, migration.FileName)
		}
	}

	if repeatRows != nil {
		repeatRows.Divider()
	}
	count := len(collection)
	if count > 0 {
		doc.Line(doc.Infof("%d Migrations", count))
		doc.NewLine()
	}

	return count
}
