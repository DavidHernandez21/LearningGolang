package main

import "fmt"

type (
	currency float32
	stringer interface {
		String() string
	}
)

func (c *currency) String() string {
	return fmt.Sprintf("$%.2f", float32(*c))
}

func main1() {
	// var c currency = 4.5
	var c = new(currency)
	*c = 45.67

	fmt.Println(c.String())

	fmt.Println(c)

	var c1 currency = 45.67

	fmt.Printf("value of c1: %v\n", c1.String())

	fmt.Printf("value of c1: %v\n", c1)

	var mainStringer stringer = c
	// mainStringer = c
	fmt.Printf("mainStringer's value: %v, type: %T\n", mainStringer, mainStringer)

	var fmtStringer fmt.Stringer = &c1
	// fmtStringer = &c1
	fmt.Printf("fmtStringer's value: %v, type: %T\n", fmtStringer, fmtStringer)
}
