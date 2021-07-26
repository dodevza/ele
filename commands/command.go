package commands

// Command ...
type Command interface {
	Execute(args []string)
	Help()
}

func getCommandByName(commandName string) Command {
	return GetCommand([]string{"executable", commandName}...)
}

// GetCommand ...
func GetCommand(args ...string) Command {
	if len(args) < 2 {
		return nil
	}
	commandName := args[1]

	commandArgs := make([]string, 0)
	if len(args) > 2 {
		commandArgs = args[2:]
	}
	switch commandName {
	case INIT:
		return newInitCommand(commandArgs)
	case MIGRATE:
		return newMigrationCommand(commandArgs)
	case ROLLBACK:
		return newRollbackCommand(commandArgs)
	case CREATEDB:
		return newCreateDBCommand(commandArgs)
	case PROMOTE:
		return newPromoteCommand(commandArgs)
	case LIST:
		return newListCommand(commandArgs)
	case LISTLONG:
		return newListCommand(commandArgs)
	case ENV:
		return newEnvCommand(commandArgs)
	case TAG:
		return newTagCommand(commandArgs)
	case WAIT:
		return newWaitCommand(commandArgs)
	default:
		return nil
	}

}
