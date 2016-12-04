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
