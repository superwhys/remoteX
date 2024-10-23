package errorutils

import "github.com/superwhys/remoteX/domain/command"

type CommandError struct {
	*BaseError
	cmdType int32
	args    command.Args
}

func ErrCommand(cmdType int32, args command.Args, opts ...ErrOption) *CommandError {
	err := &CommandError{
		BaseError: &BaseError{},
		cmdType:   cmdType,
		args:      args,
	}

	for _, opt := range opts {
		opt(err.BaseError)
	}
	return err
}

func ErrCommandMissingArguments(cmdType int32, args command.Args) *CommandError {
	return ErrCommand(cmdType, args, WithMsg("missing arguments"))
}
