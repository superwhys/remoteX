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
	
	"github.com/superwhys/remoteX/domain/node"
	"golang.org/x/time/rate"
)

const (
	// 512K
	limiterBurstSize = 4 * 128 << 10
)

type Limiter struct {
	myID                node.NodeID
	mu                  sync.Mutex
	write               *rate.Limiter
	read                *rate.Limiter
	deviceReadLimiters  map[node.NodeID]*rate.Limiter
	deviceWriteLimiters map[node.NodeID]*rate.Limiter
}

func NewLimiter(myID node.NodeID) *Limiter {
	return &Limiter{
		myID:                myID,
		mu:                  sync.Mutex{},
		write:               rate.NewLimiter(rate.Inf, limiterBurstSize),
		read:                rate.NewLimiter(rate.Inf, limiterBurstSize),
		deviceReadLimiters:  make(map[node.NodeID]*rate.Limiter),
		deviceWriteLimiters: make(map[node.NodeID]*rate.Limiter),
	}
}

func (l *Limiter) SetLimiter(device node.DeviceConfiguration) {
	readLimiter := l.getReadLimiterLocked(device.DeviceId)
	writeLimiter := l.getWriteLimiterLocked(device.DeviceId)
	
	readLimit := rate.Limit(device.MaxRecvKbps) * 1024
	writeLimit := rate.Limit(device.MaxSendKbps) * 1024
	if readLimit <= 0 {
		readLimit = rate.Inf
	}
	if writeLimit <= 0 {
		writeLimit = rate.Inf
	}
	
	readLimiter.SetLimit(readLimit)
	writeLimiter.SetLimit(writeLimit)
}

func (l *Limiter) getReadLimiterLocked(deviceId node.NodeID) *rate.Limiter {
	return l.getRateLimiter(l.deviceReadLimiters, deviceId)
}

func (l *Limiter) getWriteLimiterLocked(deviceId node.NodeID) *rate.Limiter {
	return l.getRateLimiter(l.deviceWriteLimiters, deviceId)
}

func (l *Limiter) getRateLimiter(m map[node.NodeID]*rate.Limiter, deviceId node.NodeID) *rate.Limiter {
	limiter, ok := m[deviceId]
	if !ok {
		limiter = rate.NewLimiter(rate.Inf, limiterBurstSize)
		m[deviceId] = limiter
	}
	return limiter
}

func (l *Limiter) newLimitReader(deviceId node.NodeID, r io.Reader) io.Reader {
	return &LimitReader{
		reader:     r,
		baseWaiter: (*baseWaiter)(l.getReadLimiterLocked(deviceId)),
	}
}

func (l *Limiter) getDeviceRateLimiter(deviceId node.NodeID, rw io.ReadWriter) (io.Reader, io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	r := l.newLimitReader(deviceId, rw)
	
	return r, nil
}
