package main

import "fmt"

type person1 struct {
	fname string
	lname string
	age   uint8
}

func main2() {
	mark := &person1{fname: "Mark", lname: "Smith", age: 35}

	fmt.Printf("Person's name: %v\tType: %T\n", mark.name(), mark)
	fmt.Printf("Person's age: %v\n", mark.getAge())

	mark.setAge(36)
	fmt.Printf("New person's age: %v", mark.age)
}

func (p person1) name() string {
	return fmt.Sprintf("%v, %v", p.fname, p.lname)
}

func (p person1) getAge() uint8 {
	return p.age
}

func (p *person1) setAge(a uint8) {
	if a <= 150 && a > p.age {
		fmt.Printf("Changing age of %v from %v to %v\n", p.name(), p.age, a)
		p.age = a
	}
}
