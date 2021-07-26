package loaders

import (
	"ele/constants"
	"ele/loaders/versionparsers"
	"ele/migrations"
	"ele/plugins/git"
	"ele/utils/fileio"
	"ele/utils/search"
	"strings"
)

// DirectoryLoader ...
type DirectoryLoader struct {
	reader        *fileio.DirectoryReader
	versionParser versionparsers.VersionParser
	fileSort      *PathSorter
	search        *search.Path
	ignored       *git.Ignore
	ignoreFiles   []string
	tagMap        map[string]bool
	tags          []string
}

// Load ...
func (loader *DirectoryLoader) Load() migrations.VersionTree {
	return loader.LoadPath([]string{})
}

// LoadPath ...
func (loader *DirectoryLoader) LoadPath(walkpath []string) migrations.VersionTree {
	reader := *loader.reader
	tree := migrations.EmptyVersionTreeWithTags(loader.tags)
	ignoreChain := make([]*git.Ignore, 0)
	if loader.ignored != nil {
		ignoreChain = append(ignoreChain, loader.ignored)
	}
	loader.walk(walkpath, "", "", "", "", "", reader, "", reader, &tree, ignoreChain)
	return tree
}

func (loader *DirectoryLoader) searchIncludesFile(path string) bool {
	matched, include := loader.search.Match(path)
	if matched {
		return include
	}

	return false
}

func (loader *DirectoryLoader) walk(
	path []string,
	prev string,
	module string,
	relative string,
	name string,
	firstVersion string,
	versionDir fileio.DirectoryReader,
	pathFromVersion string,
	reader fileio.DirectoryReader,
	tree *migrations.VersionTree,
	ignoreChain []*git.Ignore) {

	ignoreChain = chainIgnore(loader.ignoreFiles, ignoreChain, reader, relative)
	if shouldIgnoreDirectory(ignoreChain, relative) {
		return
	}

	newModule := module
	version, isVersion := loader.versionParser.Parse(name)
	lowerVersion := strings.ToUpper(version)
	if !isVersion && loader.tagMap[lowerVersion] {
		version = lowerVersion
		isVersion = true
	}
	foundVersion := false
	if isVersion {
		if len(firstVersion) == 0 {
			firstVersion = version
			versionDir = reader
			foundVersion = true
			pathFromVersion = ""
			module = prev
		}
	} else if len(name) > 0 {
		if len(newModule) > 0 {
			newModule += "/" + name
		} else {
			newModule = name
		}
	}

	if !foundVersion {
		if pathFromVersion != "" {
			pathFromVersion = pathFromVersion + "/" + name
		} else {
			pathFromVersion = name
		}

	}

	// Only start loading files when we reach the end of the walkpath
	if len(path) == 0 {

		includedFiles := make([]string, 0)
		for _, file := range reader.Files() {
			if !loader.searchIncludesFile(relative + "/" + file) {
				continue
			}
			if shouldIgnoreFile(ignoreChain, relative+"/"+file) {
				continue
			}

			includedFiles = append(includedFiles, file)
		}

		sortedFiles := loader.fileSort.Sort(includedFiles...)
		for _, file := range sortedFiles {
			migration := createMigrationFromFile(reader, firstVersion, newModule, relative, file, module, versionDir, pathFromVersion)
			tree.Add(*migration)
		}
	}

	nextdir := ""
	if len(path) > 0 {
		nextdir = path[0]
		path = path[1:]
	}

	for _, dir := range reader.SubDirectories() {
		if nextdir != "" && dir != nextdir {
			continue
		}
		subDir := reader.GetReader(dir)
		loader.walk(path, module, newModule, relative+"/"+dir, dir, firstVersion, versionDir, pathFromVersion, subDir, tree, ignoreChain)
	}
}

func createMigrationFromFile(
	reader fileio.DirectoryReader,
	version string,
	module string,
	relative string,
	filename string,
	prevPath string,
	versionDir fileio.DirectoryReader,
	pathFromVersion string) *migrations.Migration {

	if version == "" {
		version = constants.REPEATABLE
	}
	dep := createDependencyScanner(filename)

	fileScanner, _ := reader.FileScanner(filename)

	requires := make([]string, 0)
	// if module != prevPath && len(prevPath) > 0 {
	// 	requires = append(requires, prevPath)
	// }

	fileRequires := dep.Scan(fileScanner)

	requires = append(requires, fileRequires...)

	return &migrations.Migration{Version: version, Path: module, RelativePath: relative, FileName: filename, Requires: requires, VersionDir: versionDir, PathFromVersion: pathFromVersion, Dir: reader}
}

// NewDirectoryLoader ...
func NewDirectoryLoader(
	reader fileio.DirectoryReader,
	versionParser versionparsers.VersionParser,
	fileSorter *PathSorter,
	search *search.Path,
	ignored *git.Ignore,
	ignoreFiles []string,
	tags []string) DirectoryLoader {

	tagMap := make(map[string]bool, 0)
	for _, tag := range tags {
		key := strings.ToUpper(tag)
		tagMap[key] = true
	}
	return DirectoryLoader{
		reader:        &reader,
		versionParser: versionParser,
		fileSort:      fileSorter,
		search:        search,
		ignored:       ignored,
		ignoreFiles:   ignoreFiles,
		tagMap:        tagMap,
		tags:          tags,
	}
}
