package commands

import (
	"flag"
	"strings"
)

// CommandFlags ...
type CommandFlags struct {
	set        *flag.FlagSet
	foundFlags bool
	args       []string
	subCommand string
}

// ParseFlagSet ...
func ParseFlagSet(set *flag.FlagSet, args []string) *CommandFlags {
	onlyFlags, subCommand, foundFlags := OnlyFlags(args)
	set.Parse(onlyFlags)
	return &CommandFlags{set: set, foundFlags: foundFlags, subCommand: subCommand, args: args}
}

// NArg ...
func (cf *CommandFlags) NArg() int {
	if cf.foundFlags {
		return cf.set.NArg()
	}
	return len(cf.args)
}

// Arg ...
func (cf *CommandFlags) Arg(index int) string {
	if cf.foundFlags {
		return cf.set.Arg(index)
	}
	return cf.args[index]
}

// OnlyFlags ...
func OnlyFlags(args []string) ([]string, string, bool) {
	list := make([]string, 0)
	foundFlag := false
	lastArg := ""
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			foundFlag = true
		}

		if foundFlag {
			list = append(list, arg)
		} else {
			lastArg = arg
		}
	}

	if foundFlag == false {
		if len(args) == 0 {
			return args, "", foundFlag
		}
		return args, args[0], foundFlag
	}

	return list, lastArg, foundFlag
}
