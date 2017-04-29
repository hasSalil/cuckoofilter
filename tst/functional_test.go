package main

import (
	"cuckoofilter"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func doTestNoFalseNegative(t *testing.T, table testTable, input []uint64) {
	i := 0
	for ; i < len(input); i++ {
		if table.Add(input[i]) != cuckoofilter.Ok {
			break
		}
	}
	spew.Dump(i)
	for j := 0; j < i; j++ {
		if table.Contain(input[j]) != cuckoofilter.Ok {
			t.Fatalf("Found false negative")
		}
	}
}

func Test_NoFalseNegative(t *testing.T) {
	addCount := 2006100
	maxAddCount := 2 * addCount
	toAdd := generateRandom64s(maxAddCount)
	if len(toAdd) != maxAddCount {
		t.Fatalf("Expected %d random uint64s but got %d", maxAddCount, len(toAdd))
	}
	table12 := cuckoofilter.NewCuckooFilterTwelveBit(int64(addCount))
	doTestNoFalseNegative(t, table12, toAdd)

	table8 := cuckoofilter.NewCuckooFilterEightBit(int64(addCount))
	doTestNoFalseNegative(t, table8, toAdd)
}
