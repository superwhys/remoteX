// File:		command.go
// Created by:	Hoven
// Created on:	2024-09-30
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package command

func EmptyCommand() *Command {
	return &Command{
		Type: Empty,
		Args: []string{"empty"},
	}
}
