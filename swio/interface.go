package swio

import "io"

// ReadSeekerCloser wraps the io.Reader, io.Seeker, and io.Closer types.
type ReadSeekerCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

// CloseFunc documents a function that satisfies io.Closer.
type CloseFunc func() error

// NewReadCloser wraps r with CloseFunc c.
func NewReadCloser(r io.Reader, c CloseFunc) io.ReadCloser {
	return readCloser{r, c}
}

type readCloser struct {
	io.Reader
	c CloseFunc
}

func (r readCloser) Close() error {
	return r.c()
}

type readSeekerCloser struct {
	io.ReadSeeker
	c CloseFunc
}

func (r readSeekerCloser) Close() error {
	return r.c()
}

// NewReadSeekerCloser wraps r with the close function c.
func NewReadSeekerCloser(r io.ReadSeeker, c CloseFunc) ReadSeekerCloser {
	return readSeekerCloser{r, c}
}

// NopReadSeekerCloser wraps r with a no-op close function.
func NopReadSeekerCloser(r io.ReadSeeker) ReadSeekerCloser {
	return readSeekerCloser{r, func() error { return nil }}
}

// TeeReadCloser returns a ReadCloser that writes everything read from it to w. Its
// content is read from r. All writes must return before anything can be read. Any
// write error will be returned from Read. The Close method on the returned ReadCloser
// is called derived from r.
func TeeReadCloser(r io.ReadCloser, w io.Writer) io.ReadCloser {
	return readCloser{
		io.TeeReader(r, w),
		r.Close,
	}
}
