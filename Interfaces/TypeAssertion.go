// package main

// import (
// 	"fmt"
// )

// type (
// 	ID     string
// 	Person struct {
// 		Name string
// 	}
// )

// func main() {

// 	// var foo interface{}
// 	// foo = float32(34.566)

// 	// var f float32

// 	// f = foo.(float32)

// 	// fmt.Printf("%v, %T", f, f)
// 	done := make(chan struct{})

// 	for i := 0; i < 10; i++ {
// 		myPrintln(getValue(done))
// 		v, ok := <-done
// 		fmt.Printf("is the done channel closed: %t, value sent: %v\n", ok, v)

// 	}

// }

// func myPrintln(i interface{}) {
// 	switch i.(type) {
// 	case bool:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i.(bool))
// 	case int:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i.(int))
// 	case float64:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i.(float64))
// 	case float32:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i.(float32))
// 	case complex128:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i.(complex128))
// 	case string:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i.(string))
// 	case *Person:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i.(*Person))
// 	default:
// 		fmt.Printf("type is: %T, value is: %v\n", i, i)

// 	}
// }

// func getValue(done chan struct{}) interface{} {
// 	ch := make(chan interface{})

// 	go func() {
// 		select {
// 		case ch <- true:
// 		case ch <- 2010:
// 		case ch <- 9.15:
// 		case ch <- float32(45.6):
// 		case ch <- 7 + 6i:
// 		case ch <- "Hello World!":
// 		case ch <- &Person{Name: "Jane"}:
// 		case ch <- ID("123124213"):
// 		}

// 		fmt.Println("sending acknowledge mex to done channel")
// 		done <- struct{}{}
// 	}()

// 	return <-ch
// }
