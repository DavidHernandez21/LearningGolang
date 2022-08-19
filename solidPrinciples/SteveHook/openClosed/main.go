package main

import (
	"log"
	"math"
)

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

type square struct {
	length float64
}

func (s square) area() float64 {
	return s.length * s.length
}

type triangle struct {
	height float64
	base   float64
}

func (t triangle) area() float64 {
	return t.base * t.height / 2
}

type shape interface {
	area() float64
}

type calculator struct {
}

func (a calculator) areaSum(shapes ...shape) float64 {
	var sum float64
	for i := range shapes {
		sum += shapes[i].area()
	}
	return sum
}

func main() {
	c := circle{radius: 5}
	s := square{length: 7}
	t := triangle{height: 3, base: 7}
	calc := calculator{}
	const text = "area sum:"
	log.Println(text, calc.areaSum(c, s, t))
}
