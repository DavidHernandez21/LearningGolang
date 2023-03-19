package main

// Why the variable b is allocated on stack but the variable a escapes?
// Aren’t they both used on two
// goroutines? The reason is that the escape analysis
// module is so smart that it detects the variable b
// is never modified and thinks it is a good idea to
// use a (hidden implicit) copy of the variable b in the
// new goroutine
func main() {
	var (
		a = 1 // moved to heap: a
		b = false
		c = make(chan struct{})
	)
	go func() {
		if b {
			a++
		}
		close(c)
	}()
	<-c
	b = !b        //move b to heap
	println(a, b) // 1 false
}

//  scape analysis go run -gcflags=-m .\escape_analysis.go
