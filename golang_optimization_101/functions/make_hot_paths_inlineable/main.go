package main

func concat(bss ...[]byte) []byte {
	n := len(bss)
	if n == 0 {
		return nil
	} else if n == 1 {
		return bss[0]
	} else if n == 2 {
		return append(bss[0], bss[1]...)
	}
	var m = 0
	for i := 0; i < len(bss); i++ {
		m += len(bss[i])
	}
	var r = make([]byte, 0, m)
	for i := 0; i < len(bss); i++ {
		r = append(r, bss[i]...)
	}
	return r
}

func concat_range(bss ...[]byte) []byte {
	n := len(bss)
	if n == 0 {
		return nil
	} else if n == 1 {
		return bss[0]
	} else if n == 2 {
		return append(bss[0], bss[1]...)
	}
	var m = 0
	for i := range bss {
		m += len(bss[i])
	}
	var r = make([]byte, 0, m)
	for i := range bss {
		r = append(r, bss[i]...)
	}
	return r
}

func main() {
	var bss = [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("!!!"),
	}
	concat(bss...)
	concat_range(bss...)
}
