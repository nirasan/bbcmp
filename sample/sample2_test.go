package sample

import "testing"

func BenchmarkSample2DoubleBefore(b *testing.B) {
	b.Run("Double1", func(b *testing.B){
		f := func(n int) int { return n * 2 }
		for i := 0; i < b.N; i++ {
			f(i)
		}
	})
	b.Run("Double2", func(b *testing.B){
		f := func(n int) int { return n << 1 }
		for i := 0; i < b.N; i++ {
			f(i)
		}
	})
}

func BenchmarkSample2DoubleAfter(b *testing.B) {
	b.Run("Double1", func(b *testing.B){
		f := func(n int) int { return 2 * n }
		for i := 0; i < b.N; i++ {
			f(i)
		}
	})
	b.Run("Double2", func(b *testing.B){
		f := func(n int) int { return n + n }
		for i := 0; i < b.N; i++ {
			f(i)
		}
	})
}