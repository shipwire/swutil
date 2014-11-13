package swio

import (
	"bytes"
	"testing"
)

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
