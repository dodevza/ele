package main

import (
	"ele/commands"
	providers "ele/plugins/sql-providers"
	"ele/plugins/sql-providers/mssql"
	"ele/plugins/sql-providers/postgres"
	"ele/utils/doc"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
	"github.com/mattn/go-colorable"
)

func registerProviders() {
	providers.Register("postgres", postgres.NewProvider())
	providers.Register("mssql", mssql.NewProvider())
}

func main() {

	log.SetOutput(colorable.NewColorableStdout())
	args := os.Args

	command := commands.GetCommand(args...)
	if command != nil {
		registerProviders()
		command.Execute(args)
	} else {
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	doc.Heading(doc.Title("Usage:"), doc.Infof("%s <command>", commands.PROGRAM))

	doc.Line(doc.Info("where <command> is one of:"))
	padding := 27
	doc.Line(doc.GenerateString(" ", padding), commands.INIT)
	doc.Line(doc.GenerateString(" ", padding), commands.LIST)
	doc.Line(doc.GenerateString(" ", padding), commands.MIGRATE)
	doc.Line(doc.GenerateString(" ", padding), commands.ROLLBACK)
	doc.Line(doc.GenerateString(" ", padding), commands.TAG)
	doc.Line(doc.GenerateString(" ", padding), commands.ENV)
	doc.Line(doc.GenerateString(" ", padding), commands.PROMOTE)
	doc.Line(doc.GenerateString(" ", padding), commands.CREATEDB)
	doc.Line(doc.GenerateString(" ", padding), commands.WAIT)

	doc.NewLine()
	doc.Line(doc.Infof("%s <command> -help", commands.PROGRAM), doc.GenerateString(" ", 7), doc.Info("quick help on <command>"))
	doc.NewLine()

}
