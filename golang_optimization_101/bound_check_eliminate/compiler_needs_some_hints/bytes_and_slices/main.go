package main

type sliceOfNumbers interface {
	int | int32 | int64 | float32 | float64 | string
}

func NumSameBytes_3[T sliceOfNumbers](x, y []T) int {
	if len(x) > len(y) {
		x, y = y, x
	}
	y = y[:len(x)] // a hint, only works if T is slice
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return i
		}
	}
	return len(x)
}
func NumSameBytes_4[T sliceOfNumbers](x, y []T) int {
	if len(x) > len(y) {
		x, y = y, x
	}
	_ = y[:len(x)] // a hint, only works if T is string
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return i
		}
	}
	return len(x)
}

func main() {}
