package project

import (
	"ele/config"
	"ele/environments"
	"ele/loaders"
	"ele/loaders/versionparsers"
	"ele/utils/fileio"
	"log"
	"strings"
)

// Project ...
type Project struct {
	root          fileio.DirectoryReader
	dir           fileio.DirectoryReader
	InProject     bool
	stagingArea   *StagingArea
	config        *config.AppConfig
	versionParser versionparsers.VersionParser
	pathSorter    *loaders.PathSorter
	Environments  *environments.EnvironmentState
	tagMap        map[string]bool
	tags          []string
	ignore        []string
	ActiveStart   string
	ActiveEnd     string
}

// Stage ...
func (project *Project) Stage(search string) int {
	options := SearchOptions{Search: search, VersionStart: project.ActiveStart, VersionEnd: project.ActiveEnd}

	collections := project.Search(&options)
	return project.stagingArea.Stage(collections)
}

// Unstage ...
func (project *Project) Unstage(search string) int {
	options := SearchOptions{Search: search, VersionStart: project.ActiveStart, VersionEnd: project.ActiveEnd}

	collections := project.Search(&options)
	return project.stagingArea.Unstage(collections)
}

// Undo ...
func (project *Project) Undo() {
	project.stagingArea.Undo()
}

// Redo ...
func (project *Project) Redo() {
	log.Fatalf("Not implemented")
}

// Clear ...
func (project *Project) Clear() {
	project.stagingArea.Clear()
}

// Init ...
func (project *Project) Init() {
	project.SaveConfig()
	appender := project.root.FileAppender(".eleignore")
	defer appender.Close()
	appender.Text("")
	appender.Text("*.*")
	appender.Text("!*.sql")

	project.ActivateTags("", "")
}

// IsValidTag ...
func (project *Project) IsValidTag(tag string) bool {
	if project.tagMap[tag] {
		return true
	}
	_, isTag := project.versionParser.Parse(tag)
	return isTag
}

// Config ...
func (project *Project) Config() *config.AppConfig {
	conf := &config.AppConfig{}
	conf = conf.Assign(project.config)
	return conf
}

// Path ...
func (project *Project) Path() string {
	return project.root.Path()
}

// SaveConfig ...
func (project *Project) SaveConfig() {
	yml := project.config.Subtract(config.Defaults()).ToYaml()
	err := project.root.MoveFile("ele.yml", "ele.yml_")
	if strings.HasPrefix(yml, "{}") == false {
		appender := project.root.FileAppender("ele.yml")
		appender.Text(yml)
		appender.Close()
	}

	if err == nil {
		project.root.RemoveFile("ele.yml_")
	}

}

// NewProject ...
func NewProject(dir fileio.DirectoryReader) *Project {
	conf := config.Defaults()
	root, found := fileio.FindRoot(dir)
	foundConfig := false
	if !found {
		root = dir
	} else {
		fileScanner, err := root.FileScanner("ele.yml")
		if err == nil {
			foundConfig = false
			yml := fileio.ReadAll(fileScanner)

			projectConf := config.Parse(yml)
			conf = conf.Assign(projectConf)

		}

	}
	ignore := []string{"ele.yml", ".ele/", ".git/"}

	stagingArea := NewStagingArea(dir)
	versionParser := versionparsers.BuildFromConfig(conf)

	tagMap := make(map[string]bool)
	for _, t := range conf.Tags {
		tagMap[strings.ToUpper(t)] = true
	}

	pathSorter := loaders.NewPathSorter(conf.Hooks)
	env := environments.Load(root)
	project := Project{root: root, dir: dir, Environments: env, stagingArea: stagingArea, config: conf, versionParser: versionParser, pathSorter: pathSorter, tags: conf.Tags, tagMap: tagMap, ignore: ignore, InProject: found}

	if foundConfig {
		project.loadActiveTags()
	}
	return &project
}

// New ...
func New(path string) *Project {
	dir := fileio.NewIODirectoryReader(path)
	return NewProject(dir)
}
