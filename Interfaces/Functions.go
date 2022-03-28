package main

import "fmt"

type (
	ID     uint64
	SSN    string
	Person struct {
		Name string
		Age  uint8
		Ssn  SSN
	}

	Empty interface{}

	Printer interface {
		Print() string
	}
)

func main() {

	var e Empty
	PrintInfo(e)
	// e = 13.04
	PrintInfo(13.04)
	// e = ID(11345135)
	PrintInfo(ID(11345135))
	// e = SSN("019-72-1104")
	PrintInfo(SSN("019-72-1104"))
	// e = &Person{Name: "Jane", Age: 35, Ssn: SSN("019-72-1104")}
	PrintInfo(&Person{Name: "Jane", Age: 35, Ssn: SSN("019-72-1104")})

	var f Printer = &Person{Name: "Jane", Age: 35, Ssn: SSN("019-72-1104")}

	var g Printer = SSN("daje")

	var h Printer

	b := &Person{}

	Println(f, g, h, b)

	fmt.Printf("%v, %T", b, b)
}

func PrintInfo(e Empty) {
	fmt.Printf("e's value: %v, type: %T\n", e, e)
}

func Println(e ...Printer) {
	fmt.Println("[main.Println]")
	for i, v := range e {
		if v == nil {
			fmt.Printf("Parameter[%v]'s value is <nil>\n", i)
			continue
		}
		fmt.Printf("Parameter[%v]'s .Print() value: %v\n", i, v.Print())
	}
}

func (p *Person) Print() string {
	if p == nil {
		return "<nil>"
	}
	return p.Name
}

func (ssn SSN) Print() string {
	return string(ssn)
}
