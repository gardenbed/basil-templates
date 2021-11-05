package command

const (
	// Success is the exit code when a command execution is successful.
	Success int = iota
	// FlagError is the exit code when an undefined or invalid flag is provided to a command.
	FlagError
	// ArgError is the exit code when an invalid argument is provided to a command.
	ArgError
	// GenericError is the generic exit code when something fails.
	GenericError
)
