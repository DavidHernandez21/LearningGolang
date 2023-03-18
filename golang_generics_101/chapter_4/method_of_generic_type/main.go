package main

import "sync"

type Lockable[T any] struct {
	mu   sync.Mutex
	data T
}

func (l *Lockable[T]) Do(f func(*T)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	f(&l.data)
}
func main() {
	var n Lockable[uint32]
	n.Do(func(v *uint32) {
		*v++
	})
	println(n.data)
	var f Lockable[float64]
	f.Do(func(v *float64) {
		*v += 1.23
	})
	println(f.data)
	var b Lockable[bool]
	b.Do(func(v *bool) {
		*v = !*v
	})
	println(b.data)
	var bs Lockable[[]byte]
	bs.Do(func(v *[]byte) {
		*v = append(*v, "Go"...)
	})
	println(string(bs.data))
}
