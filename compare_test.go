package main

import (
	"testing"
	"golang.org/x/tools/benchmark/parse"
	"strings"
)

func TestCorrelate(t *testing.T) {
	benchmarks := []*parse.Benchmark{
		&parse.Benchmark{ Name: "BenchmarkSample1" },
		&parse.Benchmark{ Name: "BenchmarkSample2" },
		&parse.Benchmark{ Name: "BenchmarkSample3/Sub1" },
		&parse.Benchmark{ Name: "BenchmarkSample3/Sub2" },
		&parse.Benchmark{ Name: "BenchmarkSample4/Sub1" },
		&parse.Benchmark{ Name: "BenchmarkSample4/Sub2" },
		&parse.Benchmark{ Name: "BenchmarkSample5/Sub1" },
	}

	cmps, err := Correlate(benchmarks, "Sample1", "Sample2")
	if err != nil {
		t.Error(err)
	}
	if len(cmps) != 1 {
		t.Errorf("correlate failed: %v", cmps)
	}

	_, err = Correlate(benchmarks, "Notfound1", "Sample2")
	if err == nil || err.Error() != "before benchmark not found." {
		t.Errorf("correlate failed: %v", err)
	}

	_, err = Correlate(benchmarks, "Sample1", "Notfound2")
	if err == nil || err.Error() != "after benchmark not found." {
		t.Errorf("correlate failed: %v", err)
	}

	_, err = Correlate(benchmarks, "Sample", "Sample2")
	if err == nil || strings.Index(err.Error(), "before benchmark is ambiguous") < 0 {
		t.Errorf("correlate failed: %v", err)
	}

	_, err = Correlate(benchmarks, "Sample1", "Sample")
	if err == nil || strings.Index(err.Error(), "after benchmark is ambiguous") < 0 {
		t.Errorf("correlate failed: %v", err)
	}

	cmps, err = Correlate(benchmarks, "Sample3", "Sample4")
	if err != nil {
		t.Error(err)
	}
	if len(cmps) != 2 {
		t.Errorf("correlate failed: %v", cmps)
	}

}
