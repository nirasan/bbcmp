package sample

import (
	"testing"
)

func BenchmarkSample1Double1(b *testing.B) {
	f := func(n int) int { return n * 2}
	for i := 0; i < b.N; i++  {
		f(i)
	}
}

func BenchmarkSample1Double2(b *testing.B) {
	f := func(n int) int { return n << 1}
	for i := 0; i < b.N; i++  {
		f(i)
	}
}
