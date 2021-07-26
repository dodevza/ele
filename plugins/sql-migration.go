package plugins

import (
	"database/sql"
	"ele/environments"
	"ele/migrations"
	providers "ele/plugins/sql-providers"
	"ele/project"
	"ele/utils/doc"
	"fmt"
	"io/ioutil"
	"os"
)

// TransactionLevel ...
type TransactionLevel string

const (
	// File will commit a transaction after each file have been appliend
	File TransactionLevel = "file"
	// Version Transaction will comit after each version have been applied
	Version = "version"
	// All will commit the transaction after all migrations have been applied
	All = "all"
)

// SQLMigration ...
type SQLMigration struct {
	environment string
	txlevel     TransactionLevel
	project     *project.Project
}

// Execute ...
func (cmd *SQLMigration) Execute(env *environments.Environment, migrations migrations.MigrationCollection) error {

	provider, found := providers.Get(env.Database.Driver)
	if found == false {
		doc.Errorf("Driver %s is not registered", env.Database.Driver)
		os.Exit(1)
	}
	connectionString := provider.ToConnectionString(env.Database)
	db, err := sql.Open(env.Database.Driver, connectionString)
	if err != nil {
		return err
	}

	path := cmd.project.Path()
	for _, migration := range migrations {
		doc.Progress(doc.Info(doc.FitLeft(migration.Version, 10)), doc.Info(doc.FitLeft(migration.Path, 35)), doc.Info(doc.FitLeft(migration.FileName, 35)))
		contentBytes, err := ioutil.ReadFile(path + "/" + migration.RelativePath + "/" + migration.FileName)
		if err != nil {
			return err
		}
		sql := string(contentBytes)

		tx, txerr := db.Begin()
		if txerr != nil {
			return txerr
		}

		_, sqlErr := db.Exec(sql)
		if sqlErr != nil {
			rberror := tx.Rollback()
			if rberror != nil {
				return rberror
			}
			doc.Complete(doc.Error("Error"))
			return fmt.Errorf("%s: %s", migration.FileName, sqlErr)
		}

		commitErr := tx.Commit()
		if commitErr != nil {
			doc.Complete(doc.Error("Error"))
			return commitErr
		}

		doc.Complete(doc.Info("Complete"))
	}

	db.Close()
	doc.NewLine()
	return nil
}

// ExecuteEnvironments ...
func (cmd *SQLMigration) ExecuteEnvironments(environment string, migrations migrations.MigrationCollection) error {

	envs := cmd.project.Environments.ByName(environment)

	if len(envs) == 0 {
		return fmt.Errorf("No environments found to run for %s", environment)
	}
	for _, env := range envs {
		err := cmd.Execute(env, migrations)
		if err != nil {
			return err
		}
	}

	return nil

}

// Migrate ...
func (cmd *SQLMigration) Migrate(options *migrations.Options) error {
	migrations := cmd.MigrationPlan(options)

	envName := cmd.environment
	if options.Environment != "" {
		envName = options.Environment
	}
	cmd.PrintStatus(len(migrations), envName, options.Offset, options.Limit)
	doc.NewLine()
	return cmd.ExecuteEnvironments(envName, migrations)
}

// MigrationPlan ...
func (cmd *SQLMigration) MigrationPlan(options *migrations.Options) migrations.MigrationCollection {

	return cmd.project.Query().
		ExcludedFromSearch(cmd.project.Config().Hooks.Rollback...).
		Limit(options.Offset, options.Limit)
}

// Rollback ...
func (cmd *SQLMigration) Rollback(options *migrations.Options) error {

	migrations := cmd.RollbackPlan(options)

	envName := cmd.environment
	if options.Environment != "" {
		envName = options.Environment
	}

	cmd.PrintStatus(len(migrations), envName, options.Offset, options.Limit)
	doc.NewLine()
	return cmd.ExecuteEnvironments(envName, migrations)
}

// RollbackPlan ...
func (cmd *SQLMigration) RollbackPlan(options *migrations.Options) migrations.MigrationCollection {
	return cmd.project.Query().
		Search(cmd.project.Config().Hooks.Rollback...).
		Limit(options.Offset, options.Limit).
		Reverse()
}

// PrintStatus ...
func (cmd *SQLMigration) PrintStatus(count int, environment string, start string, end string) {
	if count == 0 {
		doc.Line(doc.Info("No migrations"))
	} else {
		doc.Line(doc.Infof("%d migrations", count))
	}
	doc.NewLine()
	cmd.PrintBadges(environment, start, end)
}

// PrintBadges ...
func (cmd *SQLMigration) PrintBadges(environment string, start string, end string) {

	if len(environment) > 0 {
		doc.Line(doc.Env(environment), doc.Tag(start, end))
	} else {
		doc.Line(doc.Tag(start, end))
	}

	doc.NewLine()
}

// NewSQLMigration ...
func NewSQLMigration(project *project.Project) *SQLMigration {

	environment := project.Environments.ActiveName()
	return &SQLMigration{
		environment: environment,
		project:     project,
	}
}
