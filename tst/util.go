package main

import (
	"cuckoofilter"
	"math/rand"
	"time"
)

type testTable interface {
	Add(uint64) cuckoofilter.CuckoofilterStatus
	Contain(uint64) cuckoofilter.CuckoofilterStatus
	Info() string
	Size() int64
	SizeInBytes() int64
}

func generateRandom64s(n int) []uint64 {
	result := make([]uint64, 0, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < n; i++ {
		result = append(result, r.Uint64())
	}
	return result
}
