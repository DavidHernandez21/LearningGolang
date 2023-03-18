package main

func main() {
	var (
		a = 1    // doesn't escape
		b = true // doesn't escape
		c = make(chan int)
	)
	b1 := b
	go func(a int) {
		if b1 {
			a++
			c <- a
		}
	}(a)
	a = <-c
	b = !b
	println(a, b) // 2 false
}
