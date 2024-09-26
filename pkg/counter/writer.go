// File:		writer.go
// Created by:	Hoven
// Created on:	2024-07-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package counter

import (
	"io"
	"sync/atomic"
	"time"
)

type CountingWriter struct {
	io.Writer
	idString string
	tot      atomic.Int64 // bytes
	last     atomic.Int64 // unix nanos
}

func (c *CountingWriter) Write(bs []byte) (int, error) {
	n, err := c.Writer.Write(bs)
	c.tot.Add(int64(n))
	totalOutgoing.Add(int64(n))
	c.last.Store(time.Now().UnixNano())
	return n, err
}

func (c *CountingWriter) Tot() int64 { return c.tot.Load() }

func (c *CountingWriter) Last() time.Time {
	return time.Unix(0, c.last.Load())
}

func TotalInOut() (int64, int64) {
	return totalIncoming.Load(), totalOutgoing.Load()
}
