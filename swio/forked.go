package swio

import (
	"bytes"
	"io"
	"sync"
)

// ForkReader accepts a reader and forks it at a given point. It returns two readers,
// one that continues from the beginning of the original reader, and another that
// reads from the nth byte. Readers are a FIFO stream, so reads from the second reader
// can't actually jump ahead. To solve this problem, reads from the tail cause
// the entire remaining contents of the head to be transparently read into memory.
func ForkReader(r io.Reader, n int) (head, tail io.Reader) {
	h := &headReader{
		r:    io.LimitReader(r, int64(n)),
		lock: &sync.Mutex{},
	}
	t := &tailReader{
		r: r,
		h: h,
	}
	return h, t
}

type headReader struct {
	b    *bytes.Buffer
	r    io.Reader
	lock *sync.Mutex
}

func (h *headReader) Cache() {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.b == nil {
		h.b = &bytes.Buffer{}
		h.b.ReadFrom(h.r)
	}
}

func (h *headReader) Read(b []byte) (n int, err error) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.b == nil {
		n, err = h.r.Read(b)
	} else {
		n, err = h.b.Read(b)
	}
	return
}

type tailReader struct {
	r io.Reader
	h *headReader
}

func (t *tailReader) Read(b []byte) (int, error) {
	t.h.Cache()
	return t.r.Read(b)
}
