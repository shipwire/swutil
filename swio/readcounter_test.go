package swio

import (
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/bradfitz/iter"
)

func TestReadCounter(t *testing.T) {
	var sum int64
	counter := NewReadCounter(DummyReader)
	for _ = range iter.N(1 << 8) {
		expect := rand.Int63n(1 << 16)
		actual, err := io.CopyN(ioutil.Discard, counter, expect)
		if err != nil {
			t.Fatal(err)
		}
		if expect != actual {
			t.Fatalf("Read %d, expected %d", actual, expect)
		}
		sum += actual
		if counted := counter.BytesRead(); counted != sum {
			t.Fatalf("Read total %d, counted %d", sum, counted)
		}
	}
}
