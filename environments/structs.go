package environments

import (
	"ele/config"
	"log"

	"gopkg.in/yaml.v3"
)

// VariableMap ...
type VariableMap map[string]interface{}

type variablesSection struct {
	Variables VariableMap `yaml:"variables,omitempty"`
}

// Environment ...
type Environment struct {
	Name     string                 `yaml:"name,omitempty"`
	Database *config.DatabaseConfig `yaml:"database,omitempty"`
	TagStart string                 `yaml:"tag_start,omitempty"`
	TagEnd   string                 `yaml:"tag_end,omitempty"`
	Tag      string                 `yaml:"tag,omitempty"`
}

// EnvironmentCollection ...
type EnvironmentCollection []*Environment

// EnvironmentGroup ...
type EnvironmentGroup struct {
	Name         string   `yaml:"name,omitempty"`
	Environments []string `yaml:"environments,omitempty"`
}

// EnvironmentFile ...
type EnvironmentFile struct {
	Variables    *VariableMap        `yaml:"variables,omitempty"`
	Environments []*Environment      `yaml:"environments,omitempty"`
	Groups       []*EnvironmentGroup `yaml:"groups,omitempty"`
}

// Parse ...
func Parse(yamlText string) *EnvironmentFile {
	conf := EnvironmentFile{}

	err := yaml.Unmarshal([]byte(yamlText), &conf)

	if err != nil {
		log.Fatalf("Could parse environment yaml config: %s", err)
	}

	return &conf
}

// ParseSecrets ...
func ParseSecrets(yamlText string) *VariableMap {
	conf := VariableMap{}

	err := yaml.Unmarshal([]byte(yamlText), &conf)

	if err != nil {
		log.Fatalf("Could parse secret yaml config: %s", err)
	}

	return &conf
}

// ParseVariables ...
func ParseVariables(yamlText string) *VariableMap {
	conf := variablesSection{}

	err := yaml.Unmarshal([]byte(yamlText), &conf)

	if err != nil {
		log.Fatalf("Could parse secret yaml config: %s", err)
	}

	return &conf.Variables
}

// ToYaml ...
func (conf *EnvironmentFile) ToYaml() string {
	byteString, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatalf("Could not convert to yaml: %s", err)
	}

	return string(byteString)
}
