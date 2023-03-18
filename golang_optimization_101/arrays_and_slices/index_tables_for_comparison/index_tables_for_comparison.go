package main

import "log"

func withoutIndexTable(n int) {
	switch n % 10 {
	case 1, 2, 6, 7, 9:
		// do something 1
		log.Println("match case n % 10 in 1,2,6,7,9")
	default:
		// do something 2
		log.Println("match case n % 10 not in 1,2,6,7,9")
	}
}

var indexTable = [10]bool{
	1: true, 2: true, 6: true, 7: true, 9: true,
}

func withIndexTable(n int) {
	switch {
	case indexTable[n%10]:
		// do something 1
		log.Println("indextable true value")
	default:
		// do something 2
		log.Println("indextable false value")

	}
}

func main() {
	n := 105
	log.Println(n % 10)
	withoutIndexTable(n)
	withIndexTable(n)
}
