package environments

import (
	"bytes"
	"ele/loaders"
	"ele/utils/fileio"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/imdatngo/mergo"
)

// Load ...
func Load(root fileio.DirectoryReader) *EnvironmentState {

	state := newState(root)

	state.LoadActiveState()
	return state
}

func (state *EnvironmentState) lazyLoad() {
	if state.loaded {
		return
	}
	state.loaded = true

	state.loadSecrets()

	files := loaders.
		Builder(state.Dir).
		IgnoreFile(".gitignore").
		Search("*ele.yml").All()

	template := template.New("Variables").Option("missingkey=error")

	for _, file := range files {
		scn, err := file.Dir.FileScanner(file.FileName)
		if err != nil {
			log.Printf("Loading env file: %s, %s", file, err)
			continue
		}
		yml := fileio.ReadAll(scn)

		variables := ParseVariables(yml)

		mergo.Merge(variables, state.secrets)

		template, _ = template.Parse(yml)
		var b bytes.Buffer
		err = template.Execute(&b, variables)

		if err != nil {
			log.Fatalf("Loading env file: %s, %s", file.FileName, "Template error")
		}

		envfile := Parse(b.String())

		if envfile == nil {
			log.Printf("Loading env file: %s, %s", file.FileName, "No environments loaded")
			continue
		}

		for _, env := range envfile.Environments {
			key := strings.ToLower(env.Name)
			state.Map[key] = env
			state.variables[key] = variables
		}

		for _, grp := range envfile.Groups {
			key := strings.ToLower(grp.Name)
			mp, exist := state.Groups[key]
			if !exist {
				mp = make(map[string]string)
			}

			for _, env := range grp.Environments {
				envKey := strings.ToLower(env)
				mp[envKey] = env
			}
			state.Groups[key] = mp
		}

	}
}

func (state *EnvironmentState) loadSecrets() {

	files := loaders.
		Builder(state.Dir).
		WorkingDirectory(state.SettingsDir).
		IgnoreFile(".gitignore").
		Search("secrets.yml").All()
	sMap := VariableMap{}
	state.secrets = &sMap
	for _, file := range files {
		scn, err := file.Dir.FileScanner(file.FileName)
		if err != nil {
			log.Printf("Loading secret file: %s, %s", file, err)
			continue
		}
		yml := fileio.ReadAll(scn)
		secretfile := ParseSecrets(yml)
		if secretfile == nil {
			log.Printf("Loading secret file: %s, %s", file, "No environments loaded")
			continue
		}
		state.secrets = secretfile
	}
	includeEnvironmentVariables(state.secrets)
}

func includeEnvironmentVariables(varmap *VariableMap) {

	mp := *varmap

	envmap := make(VariableMap, 0)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}

		envmap[pair[0]] = pair[1]
	}

	mp["env"] = &envmap
}
