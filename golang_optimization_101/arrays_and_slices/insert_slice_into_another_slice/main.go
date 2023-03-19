package main

import "log"

func Insert1(s []byte, k int, vs []byte) []byte {
	s2 := make([]byte, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}

// the elements within s2[:k] is not zeroed within the make call
func Insert2(s []byte, k int, vs []byte) []byte {
	a := s[:k]
	s2 := make([]byte, len(s)+len(vs))
	copy(s2, a)
	copy(s2[len(a):], vs)
	copy(s2[len(a)+len(vs):], s[k:])
	return s2
}

// If the free capacity of the base slice
// is large enough to hold all the inserted elements, and it is allowed
// to let the result slice and the base slice share elements,
// then the following way is the most efficient,
// for this way doesn’t allocate
func Insert3(s []byte, k int, vs []byte) []byte {
	s = s[:len(s)+len(vs)]
	copy(s[k+len(vs):], s[k:])
	copy(s[k:], vs)
	return s
}

func main() {
	s := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// log.Printf("%p", s)
	s1 := []byte{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	// s2 := Insert1(s, 3, s1)
	// s3 := Insert2(s, 3, s1)

	s5 := make([]byte, len(s)+len(s1), len(s)+len(s1)+len(s))
	log.Printf("%v", s5)
	// log.Printf("%p", s5)
	s4 := Insert3(s5, 3, s)
	// log.Printf("%p", s4)
	log.Printf("%v", s4)
	log.Printf("%v", s5)

	// log.Println(s2)
	// log.Println(s3)
}
