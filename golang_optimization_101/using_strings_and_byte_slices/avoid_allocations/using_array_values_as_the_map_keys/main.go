package main

var ma = make(map[[2]string]struct{})
var ms = make(map[string]struct{})
var keyparts = []string{
	"docs", "aaa",
	"pictures", "bbb",
	"downloads", "ccc",
}

func fa(a, b string) {
	ma[[2]string{a, b}] = struct{}{}
}
func fs(a, b string) {
	ms[a+"/"+b] = struct{}{}
	// ms[a] = struct{}{}
}

func fs1(a string) {
	ms[a] = struct{}{}
	// ms[a] = struct{}{}
}

func main() {
	// fmt.Println(len(keyparts))
	for i := 0; i < len(keyparts); i += 2 {
		// fmt.Println(keyparts[i], keyparts[i+1])
		fa(keyparts[i], keyparts[i+1])
	}
	for key := range ma {
		delete(ma, key)
	}
}
