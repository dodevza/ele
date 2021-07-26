package project

import (
	"ele/migrations"
)

// SearchOptions ...
type SearchOptions struct {
	Search        string
	ExcludeStaged bool
	OnlyStaged    bool
	VersionStart  string
	VersionEnd    string
}

// Search ...
func (project *Project) Search(options *SearchOptions) *migrations.MigrationCollection {

	search := options.Search
	if search == "" {
		search = "*"
	}
	collection := project.
		Query().
		Search(search).
		Limit(options.VersionStart, options.VersionEnd)

	if options.ExcludeStaged {
		exlcuded := project.stagingArea.ExcludeAlreadyStaged(&collection)
		collection = *exlcuded
	}

	if options.OnlyStaged {
		onlyStaged := project.stagingArea.InlcudeOnlyStaged(&collection)
		collection = *onlyStaged
	}

	return &collection
}
