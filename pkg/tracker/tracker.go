package tracker

import "io"

type TrackerReader struct {
	Reader         io.ReadCloser
	trackerManager *Manager
}

func NewTrackerReader(reader io.ReadCloser, trackerManager *Manager) *TrackerReader {
	return &TrackerReader{reader, trackerManager}
}

func (r *TrackerReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err != nil {
		return 0, err
	}

	r.trackerManager.PushDownloaded(int64(n))

	return n, err
}

func (r *TrackerReader) Close() error {
	return r.Reader.Close()
}

type TrackerWriter struct {
	Writer         io.WriteCloser
	trackerManager *Manager
}

func NewTrackerWriter(writer io.WriteCloser, trackerManager *Manager) *TrackerWriter {
	return &TrackerWriter{writer, trackerManager}
}

func (w *TrackerWriter) Write(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}

	w.trackerManager.PushUploaded(int64(n))

	return n, err
}

func (w *TrackerWriter) Close() error {
	return w.Writer.Close()
}
