package counter

import (
	"io"
	"sync/atomic"
	"time"
)

type CountingReader struct {
	io.Reader
	idString string
	tot      atomic.Int64 // bytes
	last     atomic.Int64 // unix nanos
}

func (c *CountingReader) Read(bs []byte) (int, error) {
	n, err := c.Reader.Read(bs)
	c.tot.Add(int64(n))
	totalIncoming.Add(int64(n))
	c.last.Store(time.Now().UnixNano())
	return n, err
}

func (c *CountingReader) Tot() int64 { return c.tot.Load() }

func (c *CountingReader) Last() time.Time {
	return time.Unix(0, c.last.Load())
}
