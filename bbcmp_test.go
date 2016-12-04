package main

import (
	"testing"
	"os"
)

func TestParse(t *testing.T) {
	file, err := os.Open("sample/sample1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	benchmarks, err := Parse(file)

	if err != nil {
		t.Fatal(err)
	}

	if len(benchmarks) != 2 {
		t.Errorf("parse failed: %v", benchmarks)
	}
}