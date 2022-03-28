// package main

// import "fmt"

// type (
// 	Currency float32
// 	Stringer interface {
// 		String() string
// 	}
// )

// func (c *Currency) String() string {
// 	return fmt.Sprintf("$%.2f", float32(*c))
// }

// func main() {

// 	// var c Currency = 4.5
// 	var c = new(Currency)
// 	*c = 45.67

// 	fmt.Println(c.String())

// 	fmt.Println(c)

// 	var c1 Currency = 45.67

// 	fmt.Printf("value of c1: %v\n", c1.String())

// 	fmt.Printf("value of c1: %v\n", c1)

// 	var mainStringer Stringer = c
// 	// mainStringer = c
// 	fmt.Printf("mainStringer's value: %v, type: %T\n", mainStringer, mainStringer)

// 	var fmtStringer fmt.Stringer = &c1
// 	// fmtStringer = &c1
// 	fmt.Printf("fmtStringer's value: %v, type: %T\n", fmtStringer, fmtStringer)

// }
