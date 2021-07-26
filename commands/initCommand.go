package commands

import (
	"ele/project"
	"ele/utils/doc"
	"flag"
	"log"
	"os"
)

type initCommand struct {
	help    bool
	path    string
	flagSet *flag.FlagSet
	project *project.Project
}

func (cmd *initCommand) Execute(args []string) {
	if cmd.help {
		cmd.Help()
		return
	}

	doc.NewLine()
	doc.Progress(doc.Title("Init Project"))
	cmd.project.Init()

	doc.Complete(doc.Info("Done"))
	doc.NewLine()
}

func (cmd *initCommand) Help() {
	doc.Heading(doc.Titlef("%s %s", PROGRAM, INIT))
	doc.Line("Optional Parameters:")
	doc.NewLine()
	cmd.flagSet.PrintDefaults()
	doc.NewLine()
}

func newInitCommand(args []string) *initCommand {
	flagSet := flag.NewFlagSet(INIT, flag.ExitOnError)
	cmd := initCommand{}
	flagSet.BoolVar(&cmd.help, "help", false, "Show help for command")

	flagSet.StringVar(&cmd.path, "cwd", "", "Change working directory")

	flagSet.Parse(args)

	if len(cmd.path) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not determine working directory, error: %s", err)
		}
		cmd.path = cwd
	}

	cmd.project = project.New(cmd.path)
	cmd.flagSet = flagSet

	return &cmd
}
