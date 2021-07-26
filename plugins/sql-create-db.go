package plugins

import (
	"ele/commands/sizes"
	"ele/environments"
	providers "ele/plugins/sql-providers"
	"ele/project"
	"ele/utils/doc"
	"fmt"
)

// CreateDB ...
type CreateDB struct {
	project     *project.Project
	environment string
}

// CreateDB ...
func (cmd *CreateDB) CreateDB(name string) error {

	return cmd.ExecuteOnEnvironments(name)
}

// Execute ...
func (cmd *CreateDB) Execute(env *environments.Environment) error {

	doc.Progress(doc.Info(doc.FitLeft(fmt.Sprintf("%s", env.Database.DBName), sizes.PROGRESS)), doc.Env(env.Name))

	provider, ok := providers.Get(env.Database.Driver)

	if !ok {
		doc.Complete(doc.Error("Error"))
		return fmt.Errorf("No sql provider registered for driver: %s", env.Database.Driver)
	}

	db := env.Database.Clone()

	opt := providers.CreateOptions{Name: env.Database.DBName, Database: &db}
	provider.CreateDB(opt)

	doc.Complete(doc.Info("Done"))
	doc.NewLine()
	return nil
}

// ExecuteOnEnvironments ...
func (cmd *CreateDB) ExecuteOnEnvironments(environment string) error {

	doc.Heading(doc.Title("Creating Database(s):"))
	envName := cmd.environment
	if len(environment) != 0 {
		envName = environment
	}
	envs := cmd.project.Environments.ByName(envName)

	if len(envs) == 0 {
		return fmt.Errorf("No environments found to run for %s", cmd.environment)
	}
	for _, env := range envs {
		if env.Database == nil {
			return fmt.Errorf("No database settings found to for environment: %s", env.Name)
		}
	}

	for _, env := range envs {
		err := cmd.Execute(env)
		if err != nil {
			return err
		}
	}

	return nil

}

// NewCreateDB ...
func NewCreateDB(project *project.Project) *CreateDB {
	environment := project.Environments.ActiveName()
	return &CreateDB{
		project:     project,
		environment: environment,
	}
}
