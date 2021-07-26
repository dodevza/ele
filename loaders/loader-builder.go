package loaders

import (
	"ele/config"
	"ele/loaders/versionparsers"
	"ele/migrations"
	"ele/plugins/git"
	"ele/utils/fileio"
	"ele/utils/search"
	"strings"
)

// LoaderBuilder ...
type LoaderBuilder struct {
	searchIncluded []string
	searchExcluded []string
	excludedAll    bool
	ignored        []string
	ignoredFiles   []string
	tags           []string
	root           fileio.DirectoryReader
	pathFromRoot   string
	versionParser  versionparsers.VersionParser
	pathSorter     *PathSorter
	hooks          *config.HooksConfig
	versionPrefix  string
}

// Builder Create a loader builder
func Builder(root fileio.DirectoryReader) *LoaderBuilder {
	b := LoaderBuilder{root: root, tags: []string{}, searchIncluded: []string{}, searchExcluded: []string{}, ignored: []string{}, ignoredFiles: []string{}}
	return &b
}

// Tags ...
func (b *LoaderBuilder) Tags() *migrations.VersionTree {
	loader := b.create()

	walkpath := make([]string, 0)
	if b.pathFromRoot != "" {
		walkpath = strings.Split(b.pathFromRoot, "/")
	}

	tree := loader.LoadPath(walkpath)
	return &tree
}

// Limit ...
func (b *LoaderBuilder) Limit(startTag string, endTag string) migrations.MigrationCollection {
	collection := b.Tags().Forward(startTag, endTag)
	return collection
}

// All ...
func (b *LoaderBuilder) All() migrations.MigrationCollection {
	collection := b.Tags().Forward("", "")
	return collection
}

func (b *LoaderBuilder) create() *DirectoryLoader {
	ignore := git.NewIgnore("")
	ignore.Include(b.ignored...)

	s := search.NewCollection()

	for _, srch := range b.searchIncluded {
		s.Include(srch)
	}
	if len(b.searchIncluded) == 0 {
		s.Include("*")
	}
	for _, srch := range b.searchExcluded {
		s.Include("!" + srch)
	}

	if len(b.ignoredFiles) == 0 {
		b.ignoredFiles = []string{".gitignore", ".eleignore"}
	}

	hooks := b.hooks
	if hooks == nil {
		hooks = config.Defaults().Hooks
	}
	pathSorter := b.pathSorter
	if pathSorter == nil {
		pathSorter = NewPathSorter(hooks)
	}
	versionPrefix := b.versionPrefix
	if versionPrefix == "" {
		versionPrefix = "V"
	}
	versionParser := b.versionParser
	if versionParser == nil {
		versionParser = versionparsers.NewSemanticParser(versionPrefix)
	}
	loader := NewDirectoryLoader(b.root, versionParser, pathSorter, s, ignore, b.ignoredFiles, b.tags)
	return &loader
}

// WorkingDirectory ...
func (b *LoaderBuilder) WorkingDirectory(dir fileio.DirectoryReader) *LoaderBuilder {
	b.pathFromRoot = fileio.PathFromRoot(b.root.Path(), dir.Path())
	return b
}

// SetNonVersionTags ...
func (b *LoaderBuilder) SetNonVersionTags(tags ...string) *LoaderBuilder {
	for _, t := range tags {
		b.tags = append(b.tags, t)
	}
	return b
}

// Search ...
func (b *LoaderBuilder) Search(searches ...string) *LoaderBuilder {

	for _, s := range searches {
		b.searchIncluded = append(b.searchIncluded, s)
	}
	return b
}

func (b *LoaderBuilder) exclusivelySearchForSearchTerms() {
	if b.excludedAll == false {
		b.ExcludedFromSearch("*.*") // Don't Search for anything
		b.excludedAll = true
	}
}

// ExcludedFromSearch ...
func (b *LoaderBuilder) ExcludedFromSearch(searchExcluded ...string) *LoaderBuilder {
	for _, s := range searchExcluded {
		b.searchExcluded = append(b.searchExcluded, s)
	}
	return b
}

// Ignore ...
func (b *LoaderBuilder) Ignore(ignored ...string) *LoaderBuilder {
	for _, ign := range ignored {
		b.ignored = append(b.ignored, ign)
	}
	return b
}

// IgnoreFile ...
func (b *LoaderBuilder) IgnoreFile(ignoredFile ...string) *LoaderBuilder {
	for _, ign := range ignoredFile {
		b.ignoredFiles = append(b.ignoredFiles, ign)
	}
	return b
}

// VersionParser ...
func (b *LoaderBuilder) VersionParser(vp versionparsers.VersionParser) *LoaderBuilder {
	b.versionParser = vp
	return b
}

// PathSorter ...
func (b *LoaderBuilder) PathSorter(ps *PathSorter) *LoaderBuilder {
	b.pathSorter = ps
	return b
}

// Hooks ...
func (b *LoaderBuilder) Hooks(hooks *config.HooksConfig) *LoaderBuilder {
	b.hooks = hooks
	return b
}

// VersionPrefix ...
func (b *LoaderBuilder) VersionPrefix(prefix string) *LoaderBuilder {
	b.versionPrefix = prefix
	return b
}
