package main

import (
	"testing"
)

func BenchmarkIncrement1Value(b *testing.B) {
	a := new(IntVars)
	IncrementManyTimes(&a.i1, b.N)
}

func BenchmarkIncrement2ValuesInParallelNotPadded(b *testing.B) {
	a := new(IntVars)
	IncrementParallel(&a.i1, &a.i2, b.N)
}

func BenchmarkIncrement2ValuesInParallelPadded(b *testing.B) {
	a := new(IntVarsPadded)
	IncrementParallel(&a.i1, &a.i2, b.N)
}
