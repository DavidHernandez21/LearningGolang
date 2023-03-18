package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkSprintf(b *testing.B) {
	docs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		docs[i] = fmt.Sprintf("newsfeed-%.4d.xml", i)
	}
}

func BenchmarkStringsBuilder(b *testing.B) {
	docs := make([]string, b.N)
	sb := strings.Builder{}
	sb.Grow(15)
	for i := 0; i < b.N; i++ {
		doc, err := joinBuilder(&sb, "newsfeed-", strconv.Itoa(i), ".xml")
		if err != nil {
			log.Fatalf("failed to join: %v", err)
		}

		docs[i] = doc

	}
}
