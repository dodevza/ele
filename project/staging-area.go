package project

import (
	"ele/migrations"
	"ele/utils/fileio"
	"strconv"
	"strings"
)

// StagingArea ...
type StagingArea struct {
	root   fileio.DirectoryReader
	staged fileio.DirectoryReader
	reader fileio.DirectoryReader
}

// Stage ...
func (area *StagingArea) Stage(collection *migrations.MigrationCollection) int {
	excluded := *area.ExcludeAlreadyStaged(collection)
	if len(excluded) == 0 {
		return 0
	}
	count := len(area.staged.Files()) + 1
	appender := area.staged.FileAppender(strconv.Itoa(count))
	for _, migration := range excluded {
		appender.Text(migration.GetUniqueKey())
	}
	appender.Close()

	return len(excluded)
}

// Unstage ...
func (area *StagingArea) Unstage(collection *migrations.MigrationCollection) int {
	included := *area.InlcudeOnlyStaged(collection)
	if len(included) == 0 {
		return 0
	}
	count := len(area.staged.Files()) + 1
	appender := area.staged.FileAppender(strconv.Itoa(count))
	for _, migration := range included {
		appender.Text("!" + migration.GetUniqueKey())
	}
	appender.Close()

	return len(included)
}

// Undo ...
func (area *StagingArea) Undo() {

	count := len(area.staged.Files())
	area.staged.RemoveFile(strconv.Itoa(count))
}

// Clear  ...
func (area *StagingArea) Clear() {
	for _, file := range area.staged.Files() {
		area.staged.RemoveFile(file)
	}
}

// ExcludeAlreadyStaged ...
func (area *StagingArea) ExcludeAlreadyStaged(collection *migrations.MigrationCollection) *migrations.MigrationCollection {
	migrationMap := area.getMigrationMap()
	items := make(migrations.MigrationCollection, 0)
	for _, item := range *collection {
		key := item.GetUniqueKey()
		found := migrationMap[key]
		if !found {
			items = append(items, item)
		}
	}

	return &items
}

// InlcudeOnlyStaged ...
func (area *StagingArea) InlcudeOnlyStaged(collection *migrations.MigrationCollection) *migrations.MigrationCollection {
	migrationMap := area.getMigrationMap()
	items := make(migrations.MigrationCollection, 0)
	for _, item := range *collection {
		key := item.GetUniqueKey()
		found := migrationMap[key]
		if found {
			items = append(items, item)
		}
	}

	return &items
}

func (area *StagingArea) getMigrationMap() map[string]bool {
	dir := area.staged
	mp := make(map[string]bool, 0)

	for _, file := range dir.Files() {
		scanner, _ := dir.FileScanner(file)

		for scanner.Scan() {
			key := scanner.Text()
			staged := true
			if strings.HasPrefix(key, "!") {
				key = key[1:]
				staged = false
			}
			mp[key] = staged
		}
	}

	return mp
}

// NewStagingArea ...
func NewStagingArea(reader fileio.DirectoryReader) *StagingArea {
	root, found := fileio.FindRoot(reader)
	var staged fileio.DirectoryReader
	if !found {
		root = reader
		staged = fileio.NewInMemoryReader() // In Memory don't want to persist if we don't have a project
	} else {
		staged = root.CreateDirectory(".ele").CreateDirectory("staged")
	}

	return &StagingArea{root: root, staged: staged}
}
