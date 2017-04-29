package main

import (
	"math/rand"
	"testing"
	"time"

	"cuckoofilter"

	"github.com/davecgh/go-spew/spew"
)

var lookupSampleSize = 1000000

func mixIn(x, y []uint64, yProbability float64, res []uint64) []uint64 {
	// Clone x
	for i := 0; i < len(x); i++ {
		res[i] = x[i]
	}

	// Replace first yProb*100 % of elements from random elements selected from y
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ySize := len(y)
	for i := 0; i < int(yProbability*float64(len(x))); i++ {
		res[i] = y[r.Intn(ySize)]
	}

	// Shuffle res
	for i := range res {
		j := r.Intn(i + 1)
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func cuckooBenchmarkConextFigure5(t *testing.T, table testTable, toAdd []uint64, toLookup []uint64) []float64 {
	res := []float64{}
	for i := 0; i < len(toAdd); i++ {
		if status := table.Add(toAdd[i]); status != cuckoofilter.Ok {
			break
		}
	}
	foundPercents := []float64{0.0, 0.25, 0.5, 0.75, 1.0}
	foundCount := 0
	mixed := make([]uint64, len(toLookup))
	for _, foundPercent := range foundPercents {
		mixIn(toLookup, toAdd, foundPercent, mixed)
		start := time.Now()
		for _, mixedElement := range mixed {
			if table.Contain(mixedElement) == cuckoofilter.Ok {
				foundCount++
			}
		}
		duration := time.Since(start)
		lookupTps := (float64(lookupSampleSize) / duration.Seconds()) / 1000000.0
		res = append(res, lookupTps)
	}
	if 6*lookupSampleSize == foundCount {
		t.Fatalf("All found")
	}
	return res
}

// TestConextBenchmark reproduces the CoNEXT 2014 results found in "Figure 5: Lookup
// performance when a filter achieves its capacity."
func TestConextFigure5Benchmark(t *testing.T) {
	var addCount int = 127.78 * 1000 * 1000
	maxAddCount := 2 * addCount
	toAdd := generateRandom64s(maxAddCount)
	if len(toAdd) != maxAddCount {
		t.Fatalf("Expected %d random uint64s but got %d", maxAddCount, len(toAdd))
	}
	toLookup := generateRandom64s(lookupSampleSize)
	if len(toLookup) != lookupSampleSize {
		t.Fatalf("Expected %d random uint64s but got %d", lookupSampleSize, len(toLookup))
	}
	table12 := cuckoofilter.NewCuckooFilterTwelveBit(int64(addCount))
	benchMark12 := cuckooBenchmarkConextFigure5(t, table12, toAdd, toLookup)
	spew.Dump(benchMark12)

	table8 := cuckoofilter.NewCuckooFilterEightBit(int64(addCount))
	benchMark8 := cuckooBenchmarkConextFigure5(t, table8, toAdd, toLookup)
	spew.Dump(benchMark8)
}
