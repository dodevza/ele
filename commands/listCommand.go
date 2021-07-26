package commands

import (
	"ele/constants"
	"ele/project"
	"ele/utils/doc"
	"flag"
	"log"
	"os"
)

type listCommand struct {
	help    bool
	path    string
	search  string
	staged  bool
	repeat  bool
	all     bool
	flagSet *flag.FlagSet
	project *project.Project
}

func (cmd *listCommand) Execute(args []string) {
	if cmd.help {
		cmd.Help()
		return
	}

	searchOptions := project.SearchOptions{VersionStart: cmd.project.ActiveStart, VersionEnd: cmd.project.ActiveEnd}
	searchOptions.Search = cmd.search
	if cmd.staged && !cmd.all {
		searchOptions.OnlyStaged = true
	}

	if !cmd.staged && !cmd.all {
		searchOptions.ExcludeStaged = true
	}
	searchOptions.VersionStart = cmd.project.ActiveStart
	searchOptions.VersionEnd = cmd.project.ActiveEnd

	options := doc.NewSlice()
	options = append(options, doc.Title(LISTLONG))
	if searchOptions.OnlyStaged {
		options = append(options, doc.Option("Staged"))
	} else if searchOptions.ExcludeStaged {
		options = append(options, doc.Option("Un-Staged"))
	} else {
		options = append(options, doc.Option("All"))
	}

	if cmd.repeat {
		searchOptions.VersionStart = constants.REPEATABLE
		searchOptions.VersionEnd = constants.REPEATABLE
	}
	options = append(options, doc.Env(cmd.project.Environments.ActiveName()))
	options = append(options, doc.Tag(searchOptions.VersionStart, searchOptions.VersionEnd))

	doc.Heading(options...)
	doc.Line(doc.Info(cmd.search))
	doc.NewLine()

	collection := cmd.project.Search(&searchOptions)

	rows := collection.PrintList()

	if rows > 0 {
		cmd.project.PrintBadges()
		return
	}
	if searchOptions.OnlyStaged {
		doc.Paragraph(doc.Info("No staged migrations"))
	} else if searchOptions.ExcludeStaged {
		doc.Paragraph(doc.Info("No un-staged migrations"))
	} else {
		doc.Paragraph(doc.Info("No migrations for the active tag range"), doc.Tag(searchOptions.VersionStart, searchOptions.VersionEnd))
	}
}

func (cmd *listCommand) Help() {
	doc.Heading(doc.Titlef("%s %s", PROGRAM, LISTLONG))
	doc.Line(doc.Hint("By default staged migrations is excluded from list"))
	doc.Line(doc.Infof("See: %s %s -help", PROGRAM, PROMOTE))
	doc.NewLine()
	doc.Line("Optional Parameters:")
	doc.NewLine()
	cmd.flagSet.PrintDefaults()
	doc.NewLine()
}

func newListCommand(args []string) *listCommand {
	flagSet := flag.NewFlagSet(LIST, flag.ExitOnError)
	cmd := listCommand{}
	flagSet.BoolVar(&cmd.help, "help", false, "Show help for command")
	flagSet.BoolVar(&cmd.repeat, "repeat", false, "Search only repeatable migrations")
	flagSet.BoolVar(&cmd.staged, "staged", false, "Search only staged migrations")
	flagSet.BoolVar(&cmd.all, "all", false, "Search all migrations")
	flagSet.StringVar(&cmd.path, "cwd", "", "Change working directory")

	flagSet.Parse(args)

	search := "*"
	if flagSet.NArg() > 0 {
		search = flagSet.Arg(flagSet.NArg() - 1)
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
