package main

func f4a(is []int, bs []byte) {
	if len(is) >= 256 {
		for _, n := range bs {
			_ = is[n] // Found IsInBounds
		}
	}
}

func f4b(is []int, bs []byte) {
	if len(is) >= 256 {
		is = is[:256] // a successful hint
		for _, n := range bs {
			_ = is[n] // BCEed!
		}
	}
}
func f4c(is []int, bs []byte) {
	if len(is) >= 256 {
		_ = is[:256] // a non-workable hint
		for _, n := range bs {
			_ = is[n] // Found IsInBounds
		}
	}
}
func f4d(is []int, bs []byte) {
	if len(is) >= 256 {
		_ = is[255] // a non-workable hint
		for _, n := range bs {
			_ = is[n] // Found IsInBounds
		}
	}
}

func main() {}
