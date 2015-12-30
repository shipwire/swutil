// Package swio provides additional utilities on top of the standard io package.
package swio

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// ReadSeekerAt is the interface that groups io.ReadSeeker and io.ReaderAt.
type ReadSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
}

// SeekerPosition returns the current offset of an io.Seeker
func SeekerPosition(r io.Seeker) (int64, error) {
	return r.Seek(0, 1)
}

// SeekBuffer contains a reader where all data read from it is buffered, and
// thus both readable and seekable. It essentially augments bytes.Reader so
// that all data does not need to be read at once.
type SeekBuffer struct {
	r ReadSeekerAt
}

// Read reads from the internal buffer or source reader until b is filled or an error
// occurs.
func (s *SeekBuffer) Read(b []byte) (int, error) {
	return s.r.Read(b)
}

// Seek implements io.Seeker for SeekBuffer. If the new offset goes beyond what is
// buffered, SeekBuffer will read from the source reader until offset can be reached
// or an error is returned. If whence is 2, SeekBuffer will read until EOF or an error
// is reached.
func (s *SeekBuffer) Seek(offset int64, whence int) (int64, error) {
	return s.r.Seek(offset, whence)
}

// ReadAt implements io.ReaderAt for SeekBuffer.
func (s *SeekBuffer) ReadAt(b []byte, off int64) (int, error) {
	return s.r.ReadAt(b, off)
}

type seekBuffer struct {
	src      ReadCounter
	buf      []byte
	position int64
}

// NewSeekBuffer creates a new SeekBuffer using reader r. l is the initial capacity
// of the internal buffer. If r happens to be a ReadSeeker already, r is used directly
// without any additional copies of the data.
func NewSeekBuffer(r io.Reader, l int64) *SeekBuffer {
	if rs, ok := r.(ReadSeekerAt); ok {
		return &SeekBuffer{rs}
	}
	return &SeekBuffer{&seekBuffer{
		buf: make([]byte, 0, l),
		src: NewReadCounter(r),
	}}
}

// NewSeekBufferCloser creates a new SeekBuffer using ReadCloser r. l is the initial
// capacity of the internal buffer. The close method from r is forwarded to the returned
// ReadSeekerCloser.
func NewSeekBufferCloser(r io.ReadCloser, l int64) ReadSeekerCloser {
	return NewReadSeekerCloser(NewSeekBuffer(r, l), r.Close)
}

func (r *seekBuffer) ReadAt(b []byte, off int64) (n int, err error) {
	// cannot modify state - see io.ReaderAt
	if off < 0 {
		return 0, errors.New("swio.ReadAt: negative offset")
	}

	if off < int64(len(r.buf)) {
		n = copy(b, r.buf[off:])
		off += int64(n)
		b = b[n:]
	}

	// is the offset past what we have buffered?
	if off > int64(len(r.buf)) {
		buf := make([]byte, off-int64(len(r.buf)))
		// get src to the offset
		var read int
		read, err = io.ReadFull(r.src, buf)
		r.buf = append(r.buf, buf[:read]...)
		if err != nil {
			return
		}
	}

	secondRead, err := io.ReadFull(r.src, b)
	r.buf = append(r.buf, b[:secondRead]...)

	n += secondRead

	return
}

func (r *seekBuffer) Read(b []byte) (n int, err error) {
	unread := r.buf[r.position:]

	// Get as much as we can from unread
	n = copy(b, unread)

	// If we got all we need from unread, we're done
	if len(b) == n {
		r.position += int64(n)
		return
	}

	// Fill the rest of b with r.src
	srcn, err := r.src.Read(b[n:])

	// Add new read from src to r.buf
	r.buf = append(r.buf, b[n:]...)

	// Seek to new position in r.buf, so we don't read the latest from src again
	r.position += int64(srcn + n)

	return srcn + n, err
}

func (r *seekBuffer) Seek(offset int64, whence int) (int64, error) {
	var err error
	pos := r.position

	// Get new offset (relative to beginning)
	switch whence {
	case os.SEEK_SET:
		// offset stays as it is
	case os.SEEK_CUR:
		offset = offset + pos
	case os.SEEK_END:
		var remainder []byte
		remainder, err = ioutil.ReadAll(r.src)
		r.buf = append(r.buf, remainder...)
		offset = int64(len(r.buf)) + offset
	default:
		return pos, errors.New("seek: invalid whence")
	}

	originalOffset := offset

	// Have we read far enough into src?
	if diff := offset - r.src.BytesRead(); diff > 0 {
		var n int

		buf := make([]byte, diff)

		// No. Read diff bytes into src.
		n, err = r.src.Read(buf)

		// Copy onto r.buf
		r.buf = append(r.buf, buf[:n]...)
		r.position += int64(n)

		// Make offset correction if we didn't read buf full
		offset -= diff - int64(n)
	}

	// If we reached the EOF, but we got to where we're going, ignore EOF.
	if err == io.EOF && offset == originalOffset {
		err = nil
	}

	r.position = offset
	return offset, err
}
