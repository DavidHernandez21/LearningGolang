package main

import "sync/atomic"

func IncrementManyTimesAtomic(addr *int64, times int) {
	for i := 0; i < times; i++ {
		atomic.AddInt64(addr, 1)
	}
}
