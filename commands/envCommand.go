package commands

import (
	"ele/commands/messages"
	"ele/commands/sizes"
	"ele/project"
	"ele/utils/doc"
	"flag"
	"log"
	"os"
)

type envCommand struct {
	help       bool
	inspect    bool
	create     bool
	remove     bool
	activate   bool
	deactivate bool
	list       bool
	path       string
	envName    string
	subCommand string
	flagSet    *flag.FlagSet
	project    *project.Project
}

func (cmd *envCommand) Execute(args []string) {
	if cmd.help {
		cmd.Help()
		return
	}

	if !cmd.project.InProject {
		messages.NotInProject()
	}

	if cmd.activate {
		messages.CheckEnvironmentName(cmd.envName)
		doc.NewLine()
		doc.Progress(doc.Title(doc.FitLeft("Activating", sizes.PROGRESS)), doc.Env(cmd.envName))
		err := cmd.project.Environments.Activate(cmd.envName)

		if err != nil {
			doc.Complete(doc.Error("Error"))
			doc.NewLine()
			doc.Paragraph(doc.Errorf("%s", err))
			doc.Line(doc.Hintf("%s env init <evironment name>", PROGRAM))
			doc.NewLine()
			return
		}
		doc.Complete(doc.Info("Done"))
		doc.NewLine()
		cmd.project.PrintStatus()

		return
	}

	if cmd.deactivate {
		doc.NewLine()
		doc.Progress(doc.Title(doc.FitLeft("Deactivating", sizes.PROGRESS)), doc.Env(cmd.project.Environments.ActiveName()))
		err := cmd.project.Environments.Deactivate()

		if err != nil {
			doc.Complete(doc.Error("Error"))
			doc.NewLine()
			doc.Paragraph(doc.Errorf("%s", err))
			return
		}
		doc.Complete(doc.Info("Done"))
		doc.NewLine()
		cmd.project.PrintStatus()

		return
	}

	if cmd.inspect {
		messages.CheckEnvironmentName(cmd.envName)
		doc.Heading(doc.Title("Inspect"))
		doc.Line(doc.Hintf("%s env activate <environment / group name>\n", PROGRAM))
		cmd.project.Environments.Inspect(cmd.envName)
		return
	}

	if cmd.create {
		messages.CheckEnvironmentName(cmd.envName)
		doc.Heading(doc.Titlef("Create %s", cmd.envName))
		cmd.project.Environments.Init(cmd.envName)
		return
	}

	if cmd.remove {
		doc.Heading(doc.Error("Roadmap Item - Marked for future release"))
		messages.CheckEnvironmentName(cmd.envName)
		return
	}

	if cmd.list {
		environments := cmd.project.Environments.All()
		doc.Heading(doc.Title("Environments"))
		doc.Line(doc.Hintf("%s env inspect <environment / group name>\n", PROGRAM))

		count := cmd.project.Environments.Print(environments)
		if count == 0 {
			doc.Line(doc.Info("No environments found"))
			doc.Line(doc.Hintf("%s env init <environment>\n", PROGRAM))
		}
		doc.NewLine()
		return
	}

	cmd.Help()
}

func (cmd *envCommand) Help() {
	if cmd.inspect {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, ENV, INSPECT))
		doc.Line(doc.Info("Output information for an environment"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment name>", PROGRAM, ENV, INSPECT))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.remove {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, ENV, REMOVE))
		doc.Line(doc.Info("Remove an environment in the project - Future Release"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment or group name>", PROGRAM, ENV, REMOVE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.create {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, ENV, CREATE))
		doc.Line(doc.Info("Create a new environment"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment name>", PROGRAM, ENV, CREATE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.activate {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, ENV, ACTIVATE))
		doc.Line(doc.Info("Set the active environment for the project"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment name or group name>", PROGRAM, ENV, ACTIVATE))
		doc.NewLine()
		doc.Line(doc.Hint("Examples:"))
		doc.Line(doc.Infof("%s %s %s development", PROGRAM, ENV, ACTIVATE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.deactivate {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, ENV, DEACTIVATE))
		doc.Line(doc.Info("Deactivate active environment"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s", PROGRAM, ENV, DEACTIVATE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.list {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, ENV, LIST))
		doc.Line(doc.Info("View all environments"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s", PROGRAM, ENV, LIST))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	doc.Heading(doc.Titlef("%s %s", PROGRAM, ENV))
	doc.Line(doc.Hint("Usage:"))
	doc.Line(doc.Infof("%s %s <sub command>", PROGRAM, ENV))

	doc.Line(doc.Info("where <sub command> is one of:"))
	padding := 30
	doc.Line(doc.GenerateString(" ", padding), LIST, LISTLONG)
	doc.Line(doc.GenerateString(" ", padding), CREATE, INIT)
	// doc.Line(doc.GenerateString(" ", padding), REMOVE, REMOVELONG)
	doc.Line(doc.GenerateString(" ", padding), INSPECT)
	doc.Line(doc.GenerateString(" ", padding), ACTIVATE)
	doc.Line(doc.GenerateString(" ", padding), DEACTIVATE)

	doc.NewLine()
	doc.Line(doc.Infof("%s %s <sub command> -help", PROGRAM, ENV), doc.GenerateString(" ", 7), doc.Info("quick help on <sub command>"))
	doc.NewLine()
}

func newEnvCommand(args []string) *envCommand {
	flagSet := flag.NewFlagSet(TAG, flag.ExitOnError)
	cmd := envCommand{}
	flagSet.BoolVar(&cmd.help, "help", false, "Show help for cleanup local")
	flagSet.StringVar(&cmd.path, "cwd", "", "Change working directory")

	cf := ParseFlagSet(flagSet, args)
	envName := ""
	subCommand := cf.subCommand

	if cf.NArg() > 1 {
		envName = flagSet.Arg(cf.NArg() - 1)
		subCommand = flagSet.Arg(cf.NArg() - 2)
	}

	cmd.envName = envName
	cmd.subCommand = subCommand

	switch cmd.subCommand {
	case INSPECT:
		cmd.inspect = true
		break
	case CREATE, INIT:
		cmd.create = true
		break
	case REMOVE, REMOVELONG:
		cmd.remove = true
		break
	case ACTIVATE:
		cmd.activate = true
		break
	case DEACTIVATE:
		cmd.deactivate = true
		break
	case LIST, LISTLONG:
		cmd.list = true
		break
	default:
		break
	}

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
