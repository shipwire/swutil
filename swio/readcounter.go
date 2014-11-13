package swio

import "io"

// ReadCounter is a reader that keeps track of the total number of bytes read.
type ReadCounter interface {
	io.Reader
	BytesRead() int64
}

type readCounter struct {
	r    io.Reader
	read int64
}

// NewReadCounter wraps a Reader in a ReadCounter.
func NewReadCounter(r io.Reader) ReadCounter {
	return &readCounter{r, 0}
}

func (r *readCounter) Read(b []byte) (n int, err error) {
	n, err = r.r.Read(b)
	r.read += int64(n)
	return
}

func (r readCounter) BytesRead() int64 {
	return r.read
}

// WriteCounter is a writer that keeps track of the total number of bytes written.
type WriteCounter interface {
	io.Writer
	BytesWritten() int64
}

type writeCounter struct {
	w       io.Writer
	written int64
}

// NewWriteCounter wraps a Writer in a WriteCounter
func NewWriteCounter(w io.Writer) WriteCounter {
	return &writeCounter{w, 0}
}

func (w writeCounter) BytesWritten() int64 {
	return w.written
}

func (w *writeCounter) Write(b []byte) (n int, err error) {
	n, err = w.w.Write(b)
	w.written += int64(n)
	return
}
