// File:		reader.go
// Created by:	Hoven
// Created on:	2024-07-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package limiter

import (
	"io"
)

type LimitReader struct {
	reader io.Reader
	*baseWaiter
}

func (r *LimitReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if !r.UnLimit() {
		r.take(n)
	}
	
	return n, err
}
