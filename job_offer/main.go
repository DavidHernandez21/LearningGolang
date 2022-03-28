package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
)

func main() {

	var a float64
	s := "527918932189"

	for i := 1; i < len(s)-1; i++ {

		number_1, err := strconv.ParseFloat(fmt.Sprintf("%c", s[i-1]), 32)
		if err != nil {
			log.Fatalf("could not parse to float %c: %v", s[i-1], err)
		}

		number_2, err := strconv.ParseFloat(fmt.Sprintf("%c", s[i+1]), 32)
		if err != nil {
			log.Fatalf("could not parse to float %c: %v", s[i-1], err)
		}

		number_3, err := strconv.ParseFloat(fmt.Sprintf("%c", s[i]), 32)
		if err != nil {
			log.Fatalf("could not parse to float %c: %v", s[i-1], err)
		}

		if math.Abs((number_1)-(number_2)) == number_3 {

			fmt.Printf("%c", s[i])
			a += number_3

		}
		// fmt.Printf("%c", s[i-1])
		// fmt.Println(s[i+1])
	}

	fmt.Println()
	fmt.Println(a)

}
