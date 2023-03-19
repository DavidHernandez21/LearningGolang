package main

import "log"

func MergeN_MakeCopy(ss ...[]byte) []byte {
	var n = 0
	for _, s := range ss {
		n += len(s)
	}
	i, r := 0, make([]byte, n)
	for _, s := range ss {
		copy(r[i:], s)
		i += len(s)
	}
	return r
}
func MergeN_MakeAppend(ss ...[]byte) []byte {
	var n = 0
	for _, s := range ss {
		n += len(s)
	}
	r := make([]byte, 0, n)
	for _, s := range ss {
		r = append(r, s...)
	}
	return r
}

func main() {
	a := []byte{1, 2, 3}
	b := []byte{4, 5, 6}
	c := []byte{7, 8, 9}
	d := MergeN_MakeAppend(a, b, c)
	e := MergeN_MakeCopy(a, b, c)
	log.Println(d)
	log.Println(e)

}
