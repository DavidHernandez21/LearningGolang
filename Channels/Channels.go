package main

import (
	"fmt"
	"sync"
)

// "math"
// "time"
// "math"

// func SendDataToChannel(ch chan [3]float32, arr [3]float32) {
// 	ch <- arr
// 	close(ch)
// }

// type Person struct {
// 	Name string
// 	Age  int
// }

// func SendPerson(ch chan Person, p Person) {
// 	fmt.Println("now on go routine function")
// 	// time.Sleep(6 * time.Second)
// 	fmt.Println("go routine end")
// 	ch <- p
// 	close(ch)

// }

// func main() {

// 	// arr := [3]float32{12.3, 44.67, 3.4}

// 	// ch := make(chan [3]float32)

// 	// chp := make(chan Person)

// 	// p := Person{
// 	// 	Name: "David",
// 	// 	Age:  32,
// 	// }

// 	// go SendPerson(chp, p)
// 	// fmt.Println("In main thread")

// 	// fmt.Printf("time unix now and type: %v\t%T", time.Now(), time.Now())

// 	// d := (<-chp).Name

// 	// go SendDataToChannel(ch, arr)

// 	// f, ok := <-ch

// 	// for _, v := range f {

// 	// 	fmt.Println(v)
// 	// }

// 	// println(d)
// 	// println(ok)

// 	Ncont := 10072.0
// 	Xcont := 974.0

// 	Nexp := 9886.0
// 	Xexp := 1242.0

// 	z := 1.96

// 	POcont := Xcont / Ncont

// 	POexp := Xexp / Nexp

// 	dhat := math.Abs(POcont - POexp)

// 	Ppool := (Xcont + Xexp) / (Ncont + Nexp)

// 	SEpool := math.Sqrt(Ppool * (1 - Ppool) * (1/Ncont + 1/Nexp))

// 	m := z * SEpool

// 	// fmt.Println(SEpool)
// 	fmt.Println("true probability difference: ", dhat)
// 	fmt.Println("margin: ", m)
// 	fmt.Printf("lower margin: %f\nupper margin: %f\n", dhat-m, dhat+m)

// 	if (Xexp / Nexp) > (Ppool + m) {

// 		fmt.Println("launch it")
// 	}

// 	if (Xexp / Nexp) < (Ppool + m) {

// 		fmt.Println("Do not launch it")
// 	}

// 	// SE := math.Sqrt(po * (1 - po) / N)

// 	// if (N*po) > 5 && ((1-po)*N) > 5 {
// 	// 	fmt.Println("We can safely treat it as a normal distribution")

// 	// 	m := z * SE

// 	// 	fmt.Printf("confidence intervals:\nlower:%f\t upper:%f", (po - m), (po + m))

// 	// }

// }

func sum(s []int16, c chan int16) {
	var sum int16
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func main1() {
	// s := []int16{7, 2, 8, -9, 4, 0}
	// // d := s
	// // d[3] = 8
	// c := make(chan int16)
	// go sum(s[:len(s)/2], c)
	// go sum(s[len(s)/2:], c)
	// x, y := <-c, <-c // receive from c
	// fmt.Println(x, y, x+y)
	// ch := make(chan int, 3)
	// ch <- 1
	// ch <- 2
	// ch <- 3
	// fmt.Println(<-ch)
	// fmt.Println(<-ch)
	var wg sync.WaitGroup

	wg.Add(2)

	ch := make(chan int16)

	go func(ch <-chan int16) {
		defer wg.Done()

		for i := range ch {
			fmt.Println(i)
		}
	}(ch)

	go func(ch chan<- int16) {
		defer wg.Done()
		ch <- 42
		ch <- 27
		close(ch)
	}(ch)

	wg.Wait()
}
