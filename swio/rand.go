package swio

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/rand"
	"time"
)

// DummyReader is a reader of pseudo-random data. It is meant to be more
// efficient than cryptographically random data, but is useful only in
// limited cases such as testing other readers.
var DummyReader = newDummy(time.Now().Unix())

// newDummy returns a new reader of pseudo-random data. Its returned reader
// is meant to be faster than cryptographically random data, but is useful
// in limited cases such as testing other readers. Use the same seed value
// to produce deterministic data.
func newDummy(seed int64) io.Reader {
	return &dummy{src: rand.New(rand.NewSource(seed)), buf: &bytes.Buffer{}}
}

type dummy struct {
	src *rand.Rand
	buf *bytes.Buffer
}

const bits32 = 1 << 32

func (d *dummy) Read(b []byte) (n int, err error) {
	n, err = d.buf.Read(b)

	var nn int
	for n < len(b) && err == io.EOF {
		bits := d.src.Int63n(bits32)
		binary.Write(d.buf, binary.LittleEndian, uint32(bits))
		nn, err = d.buf.Read(b[n:])
		n += nn
	}
	return n, err
}
