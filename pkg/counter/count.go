package counter

import "sync/atomic"

var (
	totalIncoming atomic.Int64
	totalOutgoing atomic.Int64
)
