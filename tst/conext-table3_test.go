package main

import (
	"cuckoofilter"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var fprSampleSize = 1000 * 1000

type cuckooFilterMetrics struct {
	AddCount int64
	Space    float64
	Fpr      float64
	Speed    float64
}

func cuckooBenchmarkConextTable3(table testTable, input []uint64) cuckooFilterMetrics {
	i := 0
	start := time.Now()
	for ; i < len(input); i++ {
		if status := table.Add(input[i]); status != cuckoofilter.Ok {
			break
		}
	}
	duration := time.Now().Sub(start)
	spew.Dump(duration)
	fpCount := 0
	absent := 0
	for ; i+absent < len(input) && absent < fprSampleSize; absent++ {
		if status := table.Contain(input[i+absent]); status == cuckoofilter.Ok {
			fpCount++
		}
	}
	return cuckooFilterMetrics{
		AddCount: table.Size(),
		Space:    float64(int64(8)*table.SizeInBytes()) / float64(i),
		Fpr:      float64(100*fpCount) / float64(absent),
		Speed:    (float64(i) / duration.Seconds()) / 1000000.0,
	}
}

// TestConextBenchmark reproduces the CoNEXT 2014 results found in "Table 3: Space efficiency
// and construction speed."
//
// Results   (CF - 12 bit)  (CF - 8 bit)
// AddCount:  127887330		 128150215
// Space:     12.593997    	 8.378774
// Fpr:       0.1859		 2.9713
// Speed:     2.711674   	 2.820561
func TestConextTable3Benchmark(t *testing.T) {
	var addCount int = 127.78 * 1000 * 1000
	maxAddCount := 2 * addCount
	ipSize := int(maxAddCount + fprSampleSize)
	input := generateRandom64s(ipSize)
	if len(input) != int(ipSize) {
		t.Fatalf("Expected %d random uint64s but got %d", ipSize, len(input))
	}
	table12 := cuckoofilter.NewCuckooFilterTwelveBit(int64(addCount))
	benchMark12 := cuckooBenchmarkConextTable3(table12, input)
	spew.Dump(benchMark12)

	table8 := cuckoofilter.NewCuckooFilterEightBit(int64(addCount))
	benchMark8 := cuckooBenchmarkConextTable3(table8, input)
	spew.Dump(benchMark8)
}
