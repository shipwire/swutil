package swio

import (
	"encoding/binary"
	"testing"

	"github.com/bradfitz/iter"
	"github.com/grd/stat"
)

func TestRand(t *testing.T) {
	var i uint8
	ints := make([]uint32, 1<<8)
	for _ = range iter.N(1 << 10) {
		binary.Read(DummyReader, binary.LittleEndian, &i)
		ints[i] += 1
	}

	if skew := stat.Skew(uint32slice(ints)); skew > 1 {
		t.Fatalf("Skew was greater than 1: %d", skew)
	}
}

type uint32slice []uint32

func (u uint32slice) Get(i int) float64 {
	return float64(u[i])
}

func (u uint32slice) Len() int {
	return len(u)
}
