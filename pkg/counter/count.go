// File:		count.go
// Created by:	Hoven
// Created on:	2024-07-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package counter

import "sync/atomic"

var (
	totalIncoming atomic.Int64
	totalOutgoing atomic.Int64
)
