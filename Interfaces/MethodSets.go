// package main

// import "fmt"

// type Person struct {
// 	fname string
// 	lname string
// 	age   uint8
// }

// func main() {

// 	mark := &Person{fname: "Mark", lname: "Smith", age: 35}

// 	fmt.Printf("Person's name: %v\tType: %T\n", mark.Name(), mark)
// 	fmt.Printf("Person's age: %v\n", mark.Age())

// 	mark.SetAge(36)
// 	fmt.Printf("New person's age: %v", mark.age)

// }

// func (p Person) Name() string {
// 	return fmt.Sprintf("%v, %v", p.fname, p.lname)
// }

// func (p Person) Age() uint8 {
// 	return p.age
// }

// func (p *Person) SetAge(a uint8) {
// 	if a <= 150 && a > p.age {
// 		fmt.Printf("Changing age of %v from %v to %v\n", p.Name(), p.age, a)
// 		p.age = a
// 	}
// }
