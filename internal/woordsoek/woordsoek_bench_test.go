package woordsoek

import (
	"fmt"
	"testing"
)

func makeWords(n int) []string {
	words := make([]string, 0, n+3)
	for i := 0; i < n; i++ {
		words = append(words, fmt.Sprintf("word%06d", i))
	}
	// add some matching words
	words = append(words, "alpha", "alphabet", "alpine")
	return words
}

func BenchmarkSearch_Serial_20k(b *testing.B) {
	words := makeWords(20000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := searchFromWords(words, "a", "lph", 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_Parallel_20k(b *testing.B) {
	words := makeWords(20000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := searchFromWordsParallel(words, "a", "lph", 0); err != nil {
			b.Fatal(err)
		}
	}
}
