package limiter

import "io"

var _ io.Writer = (*LimitWriter)(nil)

type LimitWriter struct {
	writer io.WriteCloser
	waiter
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
		w.Take(toWrite)
		n, err := w.writer.Write(buf[written : written+toWrite])
		written += n
		if err != nil {
			return written, err
		}
	}

	return written, nil
}

func (w *LimitWriter) Close() error {
	return w.writer.Close()
}
