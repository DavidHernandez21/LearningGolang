package main

import "log"

func main() {
	a := make([]int, 10)
	for k, _ := range a {
		a[k] = k
	}
	log.Println(a[2:5])
}
