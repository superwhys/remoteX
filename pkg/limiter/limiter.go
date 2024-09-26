// File:		limiter.go
// Created by:	Hoven
// Created on:	2024-07-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package limiter

import (
	"io"
	"sync"
	
	"github.com/superwhys/remoteX/pkg/common"
)

const (
	// 512K
	limiterBurstSize = 4 * 128 << 10
)

type Limiter struct {
	myID             common.NodeID
	mu               sync.Mutex
	writerWaiter     waiter
	readerWaiter     waiter
	maxRecv, maxSend int
}

func NewLimiter(nodeId common.NodeID, maxRecv, maxSend int) *Limiter {
	l := &Limiter{
		myID:         nodeId,
		mu:           sync.Mutex{},
		readerWaiter: NewBaseWaiter(maxRecv),
		writerWaiter: NewBaseWaiter(maxSend),
		maxRecv:      maxRecv,
		maxSend:      maxSend,
	}
	
	return l
}

func (l *Limiter) newLimitReader(r io.Reader) *LimitReader {
	return &LimitReader{
		reader: r,
		waiter: l.readerWaiter,
	}
}

func (l *Limiter) newLimitWriter(w io.Writer) *LimitWriter {
	return &LimitWriter{
		writer: w,
		waiter: l.writerWaiter,
	}
}

func (l *Limiter) GetNodeRateLimiter(rw io.ReadWriter) (*LimitReader, *LimitWriter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	r := l.newLimitReader(rw)
	w := l.newLimitWriter(rw)
	
	return r, w
}

func (l *Limiter) UpdateLimit(maxRecv, maxSend int32) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	l.readerWaiter.SetLimit(maxRecv)
	l.readerWaiter.SetLimit(maxSend)
}
