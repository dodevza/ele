package commands

import (
	"ele/commands/messages"
	"ele/plugins"
	"ele/project"
	"ele/utils/doc"
	"flag"
	"log"
	"os"
)

type createCommand struct {
	help        bool
	environment string
	path        string
	flagSet     *flag.FlagSet
	project     *project.Project
}

func (cmd *createCommand) Execute(args []string) {
	if cmd.help {
		cmd.Help()
		return
	}

	if !cmd.project.InProject {
		messages.NotInProject()
	}

	createDbPlugin := plugins.NewCreateDB(cmd.project)
	err := createDbPlugin.CreateDB(cmd.environment)
	if err != nil {
		doc.Paragraph(doc.Errorf("%s", err))

		if cmd.environment == "" {
			doc.Line(doc.Hint("No environment specified:"))
			doc.Line(doc.Infof("%s %s -env=<environment name or group name>", PROGRAM, CREATEDB))
			doc.Line(doc.Infof("%s %s <environment name or group name>", PROGRAM, CREATEDB))
			doc.NewLine()
		}
		if cmd.project.Environments.ActiveName() == "" {
			doc.Line(doc.Hint("No environment active:"))
			doc.Line(doc.Infof("%s %s %s <environment name or group name>", PROGRAM, ENV, ACTIVATE))
			doc.NewLine()
		}
		os.Exit(1)
	}

}

func (cmd *createCommand) Help() {
	doc.Heading(doc.Titlef("%s %s", PROGRAM, CREATEDB))
	doc.Line("Optional Parameters:")
	doc.NewLine()
	cmd.flagSet.PrintDefaults()
	doc.NewLine()
}

func newCreateDBCommand(args []string) *createCommand {
	flagSet := flag.NewFlagSet(CREATEDB, flag.ExitOnError)
	cmd := createCommand{}
	flagSet.StringVar(&cmd.environment, "env", "", "Override environment to run (defaults to active environment)")
	flagSet.BoolVar(&cmd.help, "help", false, "Show help for command")
	flagSet.StringVar(&cmd.path, "cwd", "", "Change working directory")

	cf := ParseFlagSet(flagSet, args)

	if len(cmd.path) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not determine working directory, error: %s", err)
		}
		cmd.path = cwd
	}

	if cmd.environment == "" {
		cmd.environment = cf.subCommand
	}

	cmd.project = project.New(cmd.path)
	cmd.flagSet = flagSet

	return &cmd
}
