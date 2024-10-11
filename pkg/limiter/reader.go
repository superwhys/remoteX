package limiter

import (
	"io"
)

var _ io.Reader = (*LimitReader)(nil)

type LimitReader struct {
	waiter
	reader io.Reader
}

func (r *LimitReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if !r.UnLimit() {
		r.Take(n)
	}
	
	return n, err
}
