package plugins

import (
	"ele/commands/sizes"
	"ele/environments"
	providers "ele/plugins/sql-providers"
	"ele/project"
	"ele/utils/doc"
	"fmt"
	"time"
)

// WaitDB ...
type WaitDB struct {
	project     *project.Project
	environment string
}

// WaitDB ...
func (cmd *WaitDB) WaitDB(environment string, timeOutSec int) error {

	return cmd.ExecuteOnEnvironments(environment, timeOutSec)
}

// Execute ...
func (cmd *WaitDB) Execute(env *environments.Environment, timeOutSec int) error {

	doc.Progress(doc.Info(doc.FitLeft(fmt.Sprintf("%s", env.Database.DBName), sizes.PROGRESS)), doc.Env(env.Name))

	provider, ok := providers.Get(env.Database.Driver)

	if !ok {
		doc.Complete(doc.Error("Error"))
		return fmt.Errorf("No sql provider registered for driver: %s", env.Database.Driver)
	}

	db := env.Database.Clone()

	canConnect := false
	count := 0
	timeOut := time.Second * time.Duration(timeOutSec)

	ts := time.Now()

	var err error
	for canConnect == false && time.Since(ts) < timeOut && err == nil {

		canConnect, err = provider.CanConnect(&db)
		if !canConnect && err == nil {
			doc.Cell("*")
			time.Sleep(10 * time.Second)
		}
		count++
	}
	if canConnect == true {
		doc.Complete(doc.Info("Done"))
		doc.NewLine()
		return nil
	}

	doc.Complete(doc.Info("Error"))
	if err != nil {
		return fmt.Errorf("Able to reach server on %s, returned error:  %s", db.DBName, err)
	}
	return fmt.Errorf("Timout exceeded trying to connect to %s", db.DBName)
}

// ExecuteOnEnvironments ...
func (cmd *WaitDB) ExecuteOnEnvironments(environment string, timeOutSec int) error {

	doc.Heading(doc.Title("Waiting for Database(s):"))
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
		err := cmd.Execute(env, timeOutSec)
		if err != nil {
			return err
		}
	}

	return nil

}

// NewWaitDB ...
func NewWaitDB(project *project.Project) *WaitDB {
	environment := project.Environments.ActiveName()
	return &WaitDB{
		project:     project,
		environment: environment,
	}
}
