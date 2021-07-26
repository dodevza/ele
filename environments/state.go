package environments

import (
	"ele/utils/fileio"
	"strings"
)

// EnvironmentState ...
type EnvironmentState struct {
	Dir         fileio.DirectoryReader
	SettingsDir fileio.DirectoryReader
	Groups      map[string]map[string]string
	Map         map[string]*Environment
	variables   map[string]*VariableMap
	secrets     *VariableMap
	active      string
	loaded      bool
}

// All ...
func (state *EnvironmentState) All() EnvironmentCollection {
	state.lazyLoad()
	environments := make(EnvironmentCollection, 0)

	for _, env := range state.Map {
		environments = append(environments, env)
	}

	return environments

}

// ActiveName ...
func (state *EnvironmentState) ActiveName() string {
	return state.active
}

// Active ...
func (state *EnvironmentState) Active() EnvironmentCollection {
	return state.ByName(state.active)
}

// ByName ...
func (state *EnvironmentState) ByName(name string) EnvironmentCollection {
	state.lazyLoad()
	environments := make(EnvironmentCollection, 0)
	key := strings.ToLower(name)
	env, found := state.Map[key]
	if found {
		environments = append(environments, env)
		return environments
	}

	environmentNames, foundgroup := state.Groups[key]

	if !foundgroup {
		return environments
	}

	for _, name := range environmentNames {
		ekey := strings.ToLower(name)
		env, found = state.Map[ekey]
		if found {
			environments = append(environments, env)
		}
	}

	return environments

}

// SaveActiveState ...
func (state *EnvironmentState) SaveActiveState() error {
	settingsDir := state.SettingsDir
	moveErr := settingsDir.MoveFile("env", "env_")
	fa := settingsDir.FileAppender("env")
	fa.Text(state.active)
	fa.Close()

	if moveErr == nil {
		settingsDir.RemoveFile("env_")
	}
	return nil
}

// LoadActiveState ...
func (state *EnvironmentState) LoadActiveState() error {
	settingsDir := state.SettingsDir
	fs, err := settingsDir.FileScanner("env")

	if err != nil {
		return err
	}

	env := ""
	if fs.Scan() {
		env = fs.Text()
	}

	state.active = env

	return nil
}

func newState(dir fileio.DirectoryReader) *EnvironmentState {
	root, found := fileio.FindRoot(dir)
	var settingsDir fileio.DirectoryReader
	if !found {
		root = dir
		settingsDir = fileio.NewInMemoryReader() // In Memory don't want to persist if we don't have a project
	} else {
		settingsDir = root.CreateDirectory(".ele")
	}

	groups := make(map[string]map[string]string, 0)
	mp := make(map[string]*Environment, 0)
	vars := make(map[string]*VariableMap, 0)

	return &EnvironmentState{Dir: root, SettingsDir: settingsDir, Groups: groups, Map: mp, variables: vars, loaded: false}
}
