package main

import (
	"strconv"
	"strings"
	"time"
)

func join[T any](elems []T, sep string, toString func(T) string) string {
	var buf strings.Builder
	for i, elem := range elems {
		if i > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(toString(elem))
	}
	return buf.String()
}

func main() {
	durations := []time.Duration{time.Second, time.Minute, time.Hour}
	println(join(durations, ", ", time.Duration.String))

	println(join([]int{1, 2, 3}, ", ", strconv.Itoa))

	println(join([]string{"foo", "bar"}, ", ", strconv.Quote))

}
