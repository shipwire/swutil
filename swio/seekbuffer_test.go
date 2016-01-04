package swio

import (
	"bytes"
	"io"
	"os"
	"testing"
)

type shortReader struct {
	r io.Reader
}

func (sr *shortReader) Read(p []byte) (n int, err error) {
	return sr.r.Read(p[:1])
}

func TestShortRead(t *testing.T) {
	r := NewSeekBuffer(&shortReader{bytes.NewBufferString("0123456789")}, 5)
	ch := make([]byte, 5)
	io.ReadFull(r, ch)
	if string(ch) != "01234" {
		t.Fatalf("Expected 01234, got %s", ch)
	}

	// seek back and read "4" again and the rest
	r.Seek(4, os.SEEK_SET)
	ch2 := make([]byte, 6)
	io.ReadFull(r, ch2)
	if string(ch2) != "456789" {
		t.Fatalf("Expected 456789, got %s", ch2)
	}
}

func TestSeekBuffer(t *testing.T) {
	src := bytes.NewBufferString("abcdefghijklmnopqrstuvwxyz")
	r := NewSeekBuffer(src, 25)

	ch := make([]byte, 1)
	r.Read(ch)

	if string(ch) != "a" {
		t.Fatalf("Expected a, got %s", ch)
	}

	r.Seek(5, 0)
	ch = make([]byte, 2)
	r.Read(ch)
	if string(ch) != "fg" {
		t.Fatalf("Expected fg, got %s", ch)
	}

	ch = make([]byte, 2)
	r.ReadAt(ch, 5)
	if string(ch) != "fg" {
		t.Fatalf("Expected fg, got %s", ch)
	}

	r.Seek(5, 1)
	r.Read(ch)
	if string(ch) != "mn" {
		t.Fatalf("Expected mn, got %s", ch)
	}

	r.Seek(5, 0)
	ch = make([]byte, 2)
	r.Read(ch)
	if string(ch) != "fg" {
		t.Fatalf("Expected fg, got %s", ch)
	}

	ch = make([]byte, 2)
	r.ReadAt(ch, 16)
	if string(ch) != "qr" {
		t.Fatalf("Expected qr, got %s", ch)
	}

	r.Seek(-5, 2)
	ch = make([]byte, 2)
	r.Read(ch)
	if string(ch) != "vw" {
		t.Fatalf("Expected vw, got %s", ch)
	}
}
