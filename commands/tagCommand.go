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

type tagCommand struct {
	help       bool
	remove     bool
	create     bool
	activate   bool
	deactivate bool
	inspect    bool
	list       bool
	path       string
	tagName    string
	subCommand string
	start      string
	end        string
	flagSet    *flag.FlagSet
	project    *project.Project
}

func (cmd *tagCommand) Execute(args []string) {
	if cmd.help {
		cmd.Help()
		return
	}

	if !cmd.project.InProject {
		messages.NotInProject()
	}

	if cmd.inspect {
		messages.CheckTagName(cmd.tagName)

		doc.Heading(doc.Title("Inspect"))
		count := cmd.project.InspectTag(cmd.tagName)

		if count == 0 {
			doc.Paragraph(doc.Errorf("No migrations found for tag: %s", cmd.tagName))
		}
		doc.NewLine()

		return
	}

	if cmd.remove {
		messages.CheckTagName(cmd.tagName)
		doc.NewLine()
		doc.Progress(doc.Title(doc.FitLeft("Removing tag", sizes.PROGRESS)), doc.Tag(cmd.tagName, cmd.tagName))

		cmd.project.RemoveTag(cmd.tagName)

		doc.Complete("Done")
		doc.NewLine()
		return
	}

	if cmd.create {
		messages.CheckTagName(cmd.tagName)

		doc.NewLine()
		doc.Progress(doc.Title(doc.FitLeft("Creating tag", sizes.PROGRESS)), doc.Tag(cmd.tagName, cmd.tagName))
		cmd.project.AddTag(cmd.tagName)

		doc.Complete("Done")
		doc.NewLine()
		return
	}

	if cmd.activate {
		messages.CheckTagName(cmd.tagName)
		doc.NewLine()
		doc.Progress(doc.Title(doc.FitLeft("Activating", sizes.PROGRESS)), doc.Tag(cmd.start, cmd.end))
		err := cmd.project.ActivateTags(cmd.start, cmd.end)

		if err != nil {
			doc.Complete(doc.Error("Error"))
			doc.NewLine()
			doc.Paragraph(doc.Errorf("%s", err))
			doc.Line(doc.Hintf("%s tag create <tag name>", PROGRAM))
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
		doc.Progress(doc.Title(doc.FitLeft("Deactivating", sizes.PROGRESS)), doc.Tag(cmd.project.ActiveStart, cmd.project.ActiveEnd))
		err := cmd.project.DeactivateTags()

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

	if cmd.list {
		tags := cmd.project.GetTags()
		doc.Heading(doc.Title("Tags"))
		doc.Line(doc.Hintf("%s tag inspect <tag name>\n", PROGRAM))

		count := tags.Print()
		if count == 0 {
			doc.Line(doc.Info("No tags found"))
			doc.Line(doc.Hintf("%s %s %s <tag name>\n", PROGRAM, TAG, CREATE))
			doc.Line(doc.Hint("Or create folders using with version patterns \n"))
		}
		doc.NewLine()
		return
	}

	cmd.Help()

}

func (cmd *tagCommand) Help() {
	if cmd.inspect {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, TAG, INSPECT))
		doc.Line(doc.Info("Output information for a tag"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <tag name>", PROGRAM, TAG, INSPECT))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.remove {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, TAG, REMOVE))
		doc.Line(doc.Info("Remove a custom tag in the project"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <tag name>", PROGRAM, TAG, REMOVE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.create {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, TAG, CREATE))
		doc.Line(doc.Info("Create a custom tag in the project"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <tag name>", PROGRAM, TAG, CREATE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.activate {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, TAG, ACTIVATE))
		doc.Line(doc.Info("Set the active tag range for the project"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s <tag name>", PROGRAM, TAG, ACTIVATE))
		doc.Line(doc.Infof("%s %s %s <tag start> <tag end>", PROGRAM, TAG, ACTIVATE))
		doc.Line(doc.Infof("%s %s %s -start=<tag start> -end=<tag end>", PROGRAM, TAG, ACTIVATE))
		doc.NewLine()
		doc.Line(doc.Hint("Examples:"))
		doc.Line(doc.Infof("%s %s %s V1.0.0", PROGRAM, TAG, ACTIVATE))
		doc.Line(doc.Infof("%s %s %s V1.0.0 V2.0.0", PROGRAM, TAG, ACTIVATE))
		doc.Line(doc.Infof("%s %s %s -start=V1.0.0", PROGRAM, TAG, ACTIVATE))
		doc.Line(doc.Infof("%s %s %s -end=V2.0.0", PROGRAM, TAG, ACTIVATE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.deactivate {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, TAG, DEACTIVATE))
		doc.Line(doc.Info("Reset tag range, set all version tags active"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s", PROGRAM, TAG, DEACTIVATE))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	if cmd.list {
		doc.Heading(doc.Titlef("%s %s %s", PROGRAM, TAG, LIST))
		doc.Line(doc.Info("View all tags"))
		doc.Line(doc.Hint("Usage:"))
		doc.Line(doc.Infof("%s %s %s", PROGRAM, TAG, LIST))
		doc.NewLine()
		doc.Line("Optional Parameters:")
		doc.NewLine()
		cmd.flagSet.PrintDefaults()
		doc.NewLine()
		return
	}

	doc.Heading(doc.Titlef("%s %s", PROGRAM, TAG))
	doc.Line(doc.Hint("Usage:"))
	doc.Line(doc.Infof("%s %s <sub command>", PROGRAM, TAG))

	doc.Line(doc.Info("where <sub command> is one of:"))
	padding := 30
	doc.Line(doc.GenerateString(" ", padding), LIST, LISTLONG)
	doc.Line(doc.GenerateString(" ", padding), CREATE)
	doc.Line(doc.GenerateString(" ", padding), REMOVE, REMOVELONG)
	doc.Line(doc.GenerateString(" ", padding), INSPECT)
	doc.Line(doc.GenerateString(" ", padding), ACTIVATE)
	doc.Line(doc.GenerateString(" ", padding), DEACTIVATE)

	doc.NewLine()
	doc.Line(doc.Infof("%s %s <sub command> -help", PROGRAM, TAG), doc.GenerateString(" ", 7), doc.Info("quick help on <sub command>"))
	doc.NewLine()
}

func newTagCommand(args []string) *tagCommand {
	flagSet := flag.NewFlagSet(TAG, flag.ExitOnError)
	cmd := tagCommand{}
	flagSet.BoolVar(&cmd.help, "help", false, "Show help for command")
	flagSet.StringVar(&cmd.start, "start", "", "Activate start of range")
	flagSet.StringVar(&cmd.end, "end", "", "Activate end of range")
	flagSet.StringVar(&cmd.path, "cwd", "", "Change working directory")

	cf := ParseFlagSet(flagSet, args)

	tagName := ""
	subCommand := cf.subCommand
	if cf.NArg() > 0 {
		tagName = cf.Arg(cf.NArg() - 1)
	}

	if cf.NArg() > 1 {
		subCommand = cf.Arg(cf.NArg() - 2)
	}

	if cf.NArg() > 2 {
		subCommand = cf.Arg(cf.NArg() - 3)
		tagName = cf.Arg(cf.NArg() - 2)
		cmd.start = tagName
		cmd.end = cf.Arg(cf.NArg() - 1)
	} else {
		cmd.start = tagName
		cmd.end = tagName
	}

	cmd.tagName = tagName
	cmd.subCommand = subCommand

	switch cmd.subCommand {
	case CREATE:
		cmd.create = true
		break
	case REMOVE, REMOVELONG:
		cmd.remove = true
		break
	case INSPECT:
		cmd.inspect = true
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
