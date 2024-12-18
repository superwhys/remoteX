// File:		errors_test.go
// Created by:	Hoven
// Created on:	2024-12-18
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package errorutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiWrapError(t *testing.T) {
	first := &RemoteXError{
		cause:  fmt.Errorf("first errro"),
		code:   400,
		errMsg: "first error msg",
	}

	second := &RemoteXError{
		cause:  first,
		code:   500,
		errMsg: "second error msg",
	}

	assert.Equal(t, second.Error(), "RemoteXError(code=500, cause=first errro, errMsg=first error msg: second error msg)")
}
