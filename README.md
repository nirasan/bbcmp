# bbcmp

bbcmp is better benchcmp

# What is better ?

* bbcmp can compare benchmarks of the same file.
* read from stdin and file.

# Install

```
go get golang.org/x/tools/benchmark/parse
go get github.com/nirasan/bbcmp
go install github.com/nirasan/bbcmp
```

# How to use

## Read from stdin

```
go test -run=NONE -bench . | bbcmp <BEFORE_BENCH_REGEX> <AFTER_BENCH_REGEX>
```

## Read from file

```
go test -run=NONE -bench . > result.txt
bbcmp -f result.txt <BEFORE_BENCH_REGEX> <AFTER_BENCH_REGEX>
```

# Example

## Compare normal benchmarks

### sample1_test.go

```go
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
```

### sample1.txt

```
BenchmarkSample1Double1-4   	1000000000	         2.29 ns/op
BenchmarkSample1Double2-4   	1000000000	         2.26 ns/op
PASS
ok  	github.com/nirasan/bbcmp/sample	5.016s
```

### Run

```
> bbcmp -f sample1.txt Double1 Double2
benchmark                     old ns/op     new ns/op     delta
BenchmarkSample1Double1-4     2.29          2.26          -1.31%
```

## Compare sub benchmarks

### sample2_test.go

```go
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
```

### sample2.txt

```
BenchmarkSample2DoubleBefore/Double1-4         	1000000000	         2.28 ns/op
BenchmarkSample2DoubleBefore/Double2-4         	1000000000	         2.26 ns/op
BenchmarkSample2DoubleAfter/Double1-4          	1000000000	         2.29 ns/op
BenchmarkSample2DoubleAfter/Double2-4          	1000000000	         2.26 ns/op
PASS
ok  	github.com/nirasan/bbcmp/sample	10.040s
```

### Run

```
> bbcmp -f sample2.txt Before After
benchmark                                  old ns/op     new ns/op     delta
BenchmarkSample2DoubleBefore/Double1-4     2.28          2.29          +0.44%
BenchmarkSample2DoubleBefore/Double2-4     2.26          2.26          +0.00%
```
