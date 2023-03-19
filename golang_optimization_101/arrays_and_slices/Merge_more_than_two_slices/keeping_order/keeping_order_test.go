package main

import "testing"

func BenchmarkMergeNSlicesCopy(b *testing.B) {
	slide1 := make([]byte, 1024)
	slide2 := make([]byte, 1024)
	slide3 := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MergeN_MakeCopy(slide1, slide2, slide3)
	}

}

func BenchmarkMergeNSlicesAppend(b *testing.B) {
	slide1 := make([]byte, 1024)
	slide2 := make([]byte, 1024)
	slide3 := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MergeN_MakeAppend(slide1, slide2, slide3)
	}

}
