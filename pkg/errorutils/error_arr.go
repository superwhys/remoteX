// File:		error_arr.go
// Created by:	Hoven
// Created on:	2024-11-20
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package errorutils

type ErrorArr []error

func (a ErrorArr) Error() string {
	if len(a) == 0 {
		return "no error"
	}

	s := ""
	for _, e := range a {
		s += e.Error() + "\n"
	}
	return s
}
