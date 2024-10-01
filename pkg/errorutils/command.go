// File:		command.go
// Created by:	Hoven
// Created on:	2024-10-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package errorutils

type CommandError struct {
	*BaseError
	cmdType int32
	args    map[string]string
}

func ErrCommand(cmdType int32, args map[string]string, opts ...ErrOption) *CommandError {
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

func ErrCommandMissingArguments(cmdType int32, args map[string]string) *CommandError {
	return ErrCommand(cmdType, args, WithMsg("missing arguments"))
}
