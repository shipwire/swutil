package swio

import (
	"bytes"
	"testing"
)

func TestForkReader(t *testing.T) {
	buf := bytes.NewBufferString("abcdefghijklmnopqrstuvwxyz")
	head, tail := ForkReader(buf, 5)

	one := make([]byte, 1)

	head.Read(one)
	if string(one) != "a" {
		t.Errorf("Expected first byte of head to be 'a'. Got '%s'.", string(one))
	}

	tail.Read(one)
	if string(one) != "f" {
		t.Errorf("Expected first byte of head to be 'f'. Got '%s'.", string(one))
	}

	head.Read(one)
	if string(one) != "b" {
		t.Errorf("Expected second byte of head to be 'b'. Got '%s'.", string(one))
	}

	tail.Read(one)
	if string(one) != "g" {
		t.Errorf("Expected first byte of head to be 'g'. Got '%s'.", string(one))
	}
}
