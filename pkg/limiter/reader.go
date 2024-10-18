package limiter

import (
	"io"
)

var _ io.Reader = (*LimitReader)(nil)

type LimitReader struct {
	waiter
	reader io.ReadCloser
}

func (r *LimitReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if !r.UnLimit() {
		r.Take(n)
	}

	return n, err
}

func (r *LimitReader) Close() error {
	return r.reader.Close()
}
