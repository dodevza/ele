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

type promoteCommand struct {
	help       bool
	add        bool
	remove     bool
	undo       bool
	redo       bool
	reset      bool
	promote    bool
	createTag  bool
	path       string
	search     string
	subCommand string
	flagSet    *flag.FlagSet
	project    *project.Project
}

func (cmd *promoteCommand) Execute(args []string) {
	if cmd.help {
		cmd.Help()
		return
	}

	if !cmd.project.InProject {
		messages.NotInProject()
	}

	cmd.addIfSelected()
	cmd.removeIfSelected()
	cmd.undoIfSelected()
	cmd.redoIfSelected()
	cmd.resetIfSelected()
	cmd.promoteIfSelected()
}

func (cmd *promoteCommand) addIfSelected() {
	if !cmd.add {
		return
	}
	doc.NewLine()
	doc.Progress(doc.Title(doc.FitLeft("Staging Files", sizes.PROGRESS)), doc.Tag(cmd.project.ActiveStart, cmd.project.ActiveEnd))

	count := cmd.project.Stage(cmd.search)

	doc.Complete(doc.Infof("%d files", count))
	doc.NewLine()
}

func (cmd *promoteCommand) removeIfSelected() {
	if !cmd.remove {
		return
	}
	doc.NewLine()
	doc.Progress(doc.Title(doc.FitLeft("Unstaging Files", sizes.PROGRESS)), doc.Tag(cmd.project.ActiveStart, cmd.project.ActiveEnd))

	count := cmd.project.Unstage(cmd.search)

	doc.Complete(doc.Infof("%d files", count))
	doc.NewLine()
}

func (cmd *promoteCommand) undoIfSelected() {
	if !cmd.undo {
		return
	}
	doc.NewLine()
	doc.Progress(doc.Title(doc.FitLeft("Undo last stage command", sizes.PROGRESS)), doc.Tag(cmd.project.ActiveStart, cmd.project.ActiveEnd))

	cmd.project.Undo()

	doc.Complete(doc.Info("Done"))
	doc.NewLine()
}

func (cmd *promoteCommand) redoIfSelected() {
	if !cmd.redo {
		return
	}
	doc.NewLine()
	doc.Progress(doc.Title(doc.FitLeft("Redo last stage command", sizes.PROGRESS)), doc.Tag(cmd.project.ActiveStart, cmd.project.ActiveEnd))

	cmd.project.Redo()

	doc.Complete(doc.Info("Done"))
	doc.NewLine()
}

func (cmd *promoteCommand) resetIfSelected() {
	if !cmd.reset {
		return
	}
	doc.NewLine()
	doc.Progress(doc.Title(doc.FitLeft("Reset promote", sizes.PROGRESS)), doc.Tag(cmd.project.ActiveStart, cmd.project.ActiveEnd))

	cmd.project.Clear()

	doc.Complete(doc.Info("Done"))
	doc.NewLine()
}

func (cmd *promoteCommand) promoteIfSelected() {
	if !cmd.promote {
		return
	}
	doc.NewLine()
	doc.Progress(doc.Title(doc.FitLeft("Promoting", sizes.PROGRESS)), doc.Tag(cmd.project.ActiveStart, cmd.project.ActiveEnd))

	options := project.PromoteOptions{Target: cmd.subCommand, CreateTag: cmd.createTag}
	err := cmd.project.Promote(&options)

	if err != nil {
		doc.Complete(doc.Error("Error"))
		doc.NewLine()
		doc.Paragraph(doc.Infof("%s", err))
		return
	}
	doc.Complete(doc.Info("Done"))
	doc.NewLine()
	doc.Line(doc.Hintf("Activate %s", cmd.subCommand))
	doc.Line(doc.Infof("%s %s %s %s", PROGRAM, TAG, ACTIVATE, cmd.subCommand))
	doc.NewLine()

}

func (cmd *promoteCommand) Help() {

	if cmd.remove {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, PROMOTE, REMOVE))
		doc.Line(doc.Info("Remove files from staged files"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment or group name>", PROGRAM, PROMOTE, REMOVE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.add {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, PROMOTE, ADD))
		doc.Line(doc.Info("Add files to staged files"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment name>", PROGRAM, PROMOTE, ADD))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.undo {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, PROMOTE, UNDO))
		doc.Line(doc.Info("Undo last stage add or remove"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment name>", PROGRAM, PROMOTE, UNDO))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.redo {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, PROMOTE, REDO))
		doc.Line(doc.Info("Redo last stage undo action"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment name>", PROGRAM, PROMOTE, REDO))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.reset {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, PROMOTE, UNSTAGE))
		doc.Line(doc.Info("Unstage all staged files"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <environment name>", PROGRAM, PROMOTE, UNSTAGE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	doc.Heading(doc.Titlef("%s %s", PROGRAM, PROMOTE))
	doc.Line(doc.Hint("Promote is scoped by the active tags"))
	doc.Line(doc.Infof("See: %s %s -help", PROGRAM, TAG))
	doc.NewLine()
	doc.Line(doc.Hint("Usage:"))
	doc.Line(doc.Infof("%s %s <sub command>", PROGRAM, PROMOTE))

	doc.Line(doc.Hint("Pick files to promote by staging them first"))
	doc.Line(doc.Hint("Staging files:"))
	doc.Line(doc.Infof("%s %s %s <file or path>", PROGRAM, PROMOTE, ADD))
	doc.NewLine()

	doc.Line(doc.Info("where <sub command> is one of:"))
	padding := 30
	doc.Line(doc.GenerateString(" ", padding), ADD)
	doc.Line(doc.GenerateString(" ", padding), REMOVE, REMOVELONG)
	doc.Line(doc.GenerateString(" ", padding), UNDO)
	doc.Line(doc.GenerateString(" ", padding), REDO)
	doc.Line(doc.GenerateString(" ", padding), UNSTAGE)

	doc.NewLine()
	doc.Line(doc.Infof("%s %s <sub command> -help", PROGRAM, PROMOTE), doc.GenerateString(" ", 7), doc.Info("quick help on <sub command>"))
	doc.NewLine()
}

func newPromoteCommand(args []string) *promoteCommand {
	flagSet := flag.NewFlagSet(ADD, flag.ExitOnError)
	cmd := promoteCommand{}
	flagSet.BoolVar(&cmd.help, "help", false, "Show help for command")
	flagSet.BoolVar(&cmd.createTag, "new", false, "New Tag")
	flagSet.StringVar(&cmd.path, "cwd", "", "Change working directory")

	flagArgs, lastArg, _ := OnlyFlags(args)
	flagSet.Parse(flagArgs)

	search := "*"
	if flagSet.NArg() > 1 {
		searchValue := flagSet.Arg(flagSet.NArg() - 1)
		search = searchValue
		if searchValue == "." {
			search = "*"
		}

	}

	subCommand := lastArg
	if flagSet.NArg() > 1 {
		subCommand = flagSet.Arg(flagSet.NArg() - 2)
	} else if flagSet.NArg() > 0 {
		subCommand = flagSet.Arg(flagSet.NArg() - 1)
	}
	cmd.subCommand = subCommand
	switch subCommand {
	case ADD:
		cmd.add = true
		break
	case REMOVELONG:
		cmd.remove = true
		break
	case REMOVE:
		cmd.remove = true
		break
	case UNDO:
		cmd.undo = true
		break
	case REDO:
		cmd.redo = true
		break
	case UNSTAGE:
		cmd.reset = true
		break
	default:
		cmd.promote = true
	}

	cmd.search = search

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
