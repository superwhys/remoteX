// File:		writer.go
// Created by:	Hoven
// Created on:	2024-08-01
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package limiter

import "io"

type LimitWriter struct {
	writer io.Writer
	*baseWaiter
}

func (w *LimitWriter) Write(buf []byte) (int, error) {
	if w.UnLimit() {
		return w.writer.Write(buf)
	}
	
	// Calculate the data size that can be written within 10ms and
	// round it up to the next KB to avoid too many small writes
	singleWriteSize := int(w.Limit() / 100)                 // 10ms worth of data
	singleWriteSize = ((singleWriteSize / 1024) + 1) * 1024 // round up to the next KB
	if singleWriteSize > limiterBurstSize {
		singleWriteSize = limiterBurstSize
	}
	
	written := 0
	for written < len(buf) {
		toWrite := singleWriteSize
		if toWrite > len(buf)-written {
			toWrite = len(buf) - written
		}
		w.take(toWrite)
		n, err := w.writer.Write(buf[written : written+toWrite])
		written += n
		if err != nil {
			return written, err
		}
	}
	
	return written, nil
}
