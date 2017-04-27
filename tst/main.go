package main

import (
	"cuckoofilter"
	"fmt"
)

func main() {
	filter := cuckoofilter.NewCuckooFilterEightBit(int64(5))
	status := filter.Add(123)
	status = filter.Add(12)
	status = filter.Add(1)

	fmt.Println(byte(status))
	fmt.Println(filter.Info())
	for i := 1; i < 200; i++ {
		if cuckoofilter.Ok == filter.Contain(uint64(i)) {
			fmt.Printf("positive: %d\n", i)
		}
	}
}
