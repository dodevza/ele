package commands

import (
	"ele/migrations"
	"ele/plugins"
	"ele/project"
	"ele/utils/doc"
	"flag"
	"log"
	"os"
)

type rollbackCommand struct {
	help    bool
	plan    bool
	path    string
	options *migrations.Options
	flagSet *flag.FlagSet
	project *project.Project
}

func (cmd *rollbackCommand) Execute(args []string) {
	if cmd.help {
		cmd.Help()
		return
	}

	runner := plugins.NewSQLMigration(cmd.project)

	if cmd.plan {
		doc.Heading(doc.Title("Rollback Plan"))
		collection := runner.RollbackPlan(cmd.options)
		collection.PrintExcecutionPlan()
		doc.NewLine()
		runner.PrintStatus(len(collection), cmd.project.Environments.ActiveName(), cmd.options.Offset, cmd.options.Limit)
		doc.NewLine()
		return
	}

	doc.Heading(doc.Title("Rollback"))

	err := runner.Rollback(cmd.options)
	if err != nil {
		doc.Paragraph(doc.Errorf("%s", err))
		return
	}

	doc.Line(doc.Success("Completed All"))
	doc.NewLine()
}

func (cmd *rollbackCommand) Help() {
	doc.Heading(doc.Titlef("%s %s", PROGRAM, ROLLBACK))
	doc.Line(doc.Hint("Migrations defaults to settings provided in by the active environment and tags"))
	doc.Line(doc.Infof("See: %s %s -help", PROGRAM, ENV))
	doc.Line(doc.Infof("See: %s %s -help", PROGRAM, TAG))
	doc.NewLine()
	doc.Line(doc.Hint("Examples:"))
	doc.Line(doc.Infof("%s %s -start=V1.0.0", PROGRAM, ROLLBACK))
	doc.Line(doc.Infof("%s %s -start=UAT", PROGRAM, ROLLBACK))
	doc.Line(doc.Infof("%s %s -start=V1.0.0 -end=V.1.1.0 -env=develop", PROGRAM, ROLLBACK))
	doc.NewLine()
	doc.Line("Optional Parameters:")
	doc.NewLine()
	cmd.flagSet.PrintDefaults()
	doc.NewLine()
}

func newRollbackCommand(args []string) *rollbackCommand {
	flagSet := flag.NewFlagSet(ROLLBACK, flag.ExitOnError)
	cmd := rollbackCommand{}
	opts := migrations.Options{}

	flagSet.BoolVar(&cmd.help, "help", false, "Show help for command")
	flagSet.StringVar(&opts.Offset, "start", "", "Execute from specific version (defaults to active tag range)")
	flagSet.StringVar(&opts.Limit, "end", "", "Execute up to specific version (defaults to active tag range)")
	flagSet.StringVar(&opts.Environment, "env", "", "Override environment to run (defaults to active environment)")
	flagSet.BoolVar(&cmd.plan, "plan", false, "Output files that will be run")
	flagSet.StringVar(&cmd.path, "cwd", "", "Change working directory")

	flagSet.Parse(args)

	if len(cmd.path) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not determine working directory, error: %s", err)
		}
		cmd.path = cwd
	}

	cmd.flagSet = flagSet

	cmd.project = project.New(cmd.path)

	if len(opts.Offset) == 0 && len(opts.Limit) == 0 {
		opts.Offset = cmd.project.ActiveStart
		opts.Limit = cmd.project.ActiveEnd
	}

	cmd.options = &opts

	return &cmd
}
