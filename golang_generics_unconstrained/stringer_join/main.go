package main

import (
	"fmt"
	"strings"
	"time"
)

func join[T fmt.Stringer](elems []T, sep string) string {
	var buf strings.Builder
	for i, elem := range elems {
		if i > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(elem.String())
	}
	return buf.String()
}

func main() {
	durations := []time.Duration{time.Second, time.Minute, time.Hour}
	fmt.Println(join(durations, ", "))
}
