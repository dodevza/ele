package project

import (
	"ele/loaders"
)

// QueryFromRoot ...
func (project *Project) QueryFromRoot() *loaders.LoaderBuilder {

	return loaders.
		Builder(project.root).
		VersionParser(project.versionParser).
		PathSorter(project.pathSorter).
		Ignore(project.ignore...).
		SetNonVersionTags(project.tags...)
}

// Query ...
func (project *Project) Query() *loaders.LoaderBuilder {

	return project.QueryFromRoot().WorkingDirectory(project.dir)
}
