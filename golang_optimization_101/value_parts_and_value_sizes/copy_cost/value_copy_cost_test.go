package copycost

import "testing"

type T4 struct{ a, b, c, d float32 }
type T5 struct{ a, b, c, d, e float32 }

var t4 T4
var t5 T5

//go:noinline
func Add4(x, y T4) (z T4) {
	z.a = x.a + y.a
	z.b = x.b + y.b
	z.c = x.c + y.c
	z.d = x.d + y.d
	return
}

//go:noinline
func Add5(x, y T5) (z T5) {
	z.a = x.a + y.a
	z.b = x.b + y.b
	z.c = x.c + y.c
	z.d = x.d + y.d
	z.e = x.e + y.e
	return
}

// The //go:noinline compiler directives used here are to prevent the calls to the two function from
// being inlined. If the directives are removed, the Add4 function will become even more performant.

func Benchmark_Add4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x, y T4
		t4 = Add4(x, y)
	}
}
func Benchmark_Add5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var x, y T5
		t5 = Add5(x, y)
	}
}
