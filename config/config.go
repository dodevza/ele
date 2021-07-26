package config

import (
	"ele/constants"
	"log"

	"gopkg.in/yaml.v3"
)

// AppConfig ...
type AppConfig struct {
	DefaultTag     string           `yaml:"default_tag,omitempty"`
	Hooks          *HooksConfig     `yaml:"hooks,omitempty"`
	Tags           []string         `yaml:"tags,omitempty"`
	VersionParsers []*VersionConfig `yaml:"version_parsers,omitempty"`
}

// HooksConfig ...
type HooksConfig struct {
	Before   []string `yaml:"before,omitempty"`
	After    []string `yaml:"after,omitempty"`
	Rollback []string `yaml:"rollback,omitempty"`
}

// DatabaseConfig ...
type DatabaseConfig struct {
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	DBName   string `yaml:"database,omitempty"`
	SSLMode  string `yaml:"ssl_mode,omitempty"`
	Driver   string `yaml:"driver,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
}

// VersionConfig ...
type VersionConfig struct {
	Type      string `yaml:"type,omitempty"`
	Parameter string `yaml:"param,omitempty"`
}

// Clone ...
func (dbconf DatabaseConfig) Clone() DatabaseConfig {
	// Use Stack to copy config
	return dbconf
}

// Defaults ...
func Defaults() *AppConfig {
	tags := []string{}
	hooks := HooksConfig{Before: []string{"before.", "bf."}, After: []string{"after.", "af."}, Rollback: []string{"*.rollback.*", "*.rb.*"}}
	//database := DatabaseConfig{User: "postgres", DBName: "test", Password: "test9*8&PSQL", SSLMode: "disable", Driver: "postgres"}
	versionParsers := []*VersionConfig{}
	versionParsers = append(versionParsers, &VersionConfig{Type: constants.REGEX, Parameter: "(?i)^_"})
	versionParsers = append(versionParsers, &VersionConfig{Type: constants.DATE, Parameter: "YYYY-MM-DD"})
	versionParsers = append(versionParsers, &VersionConfig{Type: constants.SEMANTIC, Parameter: "V"})
	return &AppConfig{Hooks: &hooks, Tags: tags, VersionParsers: versionParsers, DefaultTag: constants.REPEATABLE}
}

// Parse ...
func Parse(yamlText string) *AppConfig {
	conf := AppConfig{}

	err := yaml.Unmarshal([]byte(yamlText), &conf)

	if err != nil {
		log.Fatalf("Could parse yaml config: %s", err)
	}

	return &conf
}

// ToYaml ...
func (conf *AppConfig) ToYaml() string {
	byteString, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatalf("Could not convert to yaml: %s", err)
	}

	return string(byteString)
}

// Assign ...
func (conf *AppConfig) Assign(from *AppConfig) *AppConfig {
	result := AppConfig{}
	result.DefaultTag = conf.DefaultTag
	if from == nil {
		from = &AppConfig{}
	}
	if from.DefaultTag != "" {
		result.DefaultTag = from.DefaultTag
	}
	result.Hooks = conf.mergeHooksConfig(from.Hooks)
	result.Tags = conf.mergeTags(from.Tags)
	result.VersionParsers = conf.mergeVersionParsers(from.VersionParsers)
	return &result
}

// Subtract Remove properties that is the same as deafult
func (conf *AppConfig) Subtract(from *AppConfig) *AppConfig {
	result := conf.Assign(nil)

	if from == nil {
		return result
	}

	if conf.Hooks != nil && from.Hooks != nil {
		if IsTheSameArray(conf.Hooks.After, from.Hooks.After) {
			result.Hooks.After = nil
		}
		if IsTheSameArray(conf.Hooks.Before, from.Hooks.Before) {
			result.Hooks.Before = nil
		}
		if IsTheSameArray(conf.Hooks.Rollback, from.Hooks.Rollback) {
			result.Hooks.Rollback = nil
		}

	}

	h := result.Hooks
	if h != nil && h.After == nil && h.Before == nil && h.Rollback == nil {
		result.Hooks = nil
	}

	if IsTheSameArray(conf.Tags, from.Tags) {
		result.Tags = nil
	}

	if IsTheVersionParsers(conf.VersionParsers, from.VersionParsers) {
		result.VersionParsers = nil
	}

	if conf.DefaultTag == from.DefaultTag {
		result.DefaultTag = ""
	}

	return result
}

// IsTheSameArray ...
func IsTheSameArray(a []string, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	if b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for index, avalue := range a {
		if avalue != b[index] {
			return false
		}
	}

	return true
}

// IsTheVersionParsers ...
func IsTheVersionParsers(a []*VersionConfig, b []*VersionConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	if b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}

	for index, avalue := range a {
		if *avalue != *b[index] {
			return false
		}
	}

	return true
}

func (conf *AppConfig) mergeTags(from []string) []string {
	result := make([]string, 0)
	if from == nil && conf.Tags == nil {
		return nil
	}

	var tagsToUse []string
	if from != nil {
		tagsToUse = from
	} else {
		tagsToUse = conf.Tags
	}

	for _, t := range tagsToUse {
		result = append(result, t)
	}

	return result
}

func (conf *AppConfig) mergeVersionParsers(from []*VersionConfig) []*VersionConfig {
	result := make([]*VersionConfig, 0)
	if from == nil && conf.VersionParsers == nil {
		return nil
	}

	var parsersToUse []*VersionConfig
	if from != nil {
		parsersToUse = from
	} else {
		parsersToUse = conf.VersionParsers
	}

	for _, t := range parsersToUse {
		result = append(result, t)
	}

	return result
}

func (conf *AppConfig) mergeHooksConfig(from *HooksConfig) *HooksConfig {
	result := HooksConfig{}

	if from == nil && conf.Hooks == nil {
		return nil
	}

	if from != nil {
		result.After = from.After
		result.Before = from.Before
		result.Rollback = from.Rollback
	}

	if conf.Hooks != nil {
		if len(result.After) == 0 {
			result.After = conf.Hooks.After
		}
		if len(result.Before) == 0 {
			result.Before = conf.Hooks.Before
		}
		if len(result.Rollback) == 0 {
			result.Rollback = conf.Hooks.Rollback
		}
	}

	return &result
}
