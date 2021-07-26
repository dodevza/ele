package environments

import (
	"ele/config"
	"fmt"
	"strings"
)

// Init ...
func (state *EnvironmentState) Init(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("Did not provide a name")
	}
	state.lazyLoad()
	key := strings.ToLower(name)
	_, foundFile := state.Map[key]
	if foundFile {
		return fmt.Errorf("Environment allready exist")
	}

	_, foundGroup := state.Groups[key]

	if foundGroup {
		return fmt.Errorf("Environment Group allready exist")
	}

	conf := EnvironmentFile{}
	conf.Environments = make([]*Environment, 0)
	db := config.DatabaseConfig{Driver: "postgres", DBName: name, User: name, Password: "apassword", SSLMode: "disable", Host: "127.0.0.1", Port: 5432}
	env := Environment{Name: name, Database: &db}
	conf.Environments = append(conf.Environments, &env)

	yml := conf.ToYaml()

	filename := key + ".ele.yml"
	err := state.Dir.MoveFile(filename, filename+"_")

	fa := state.Dir.FileAppender(filename)
	fa.Text(yml)
	fa.Close()
	if err == nil {
		state.Dir.RemoveFile(filename + "_")
	}

	return nil
}
