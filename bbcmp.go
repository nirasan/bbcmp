package main

import (
	"golang.org/x/tools/benchmark/parse"
	"flag"
	"bufio"
	"os"
	"fmt"
	"sort"
	"text/tabwriter"
	"strconv"
)

var (
	changedOnly = flag.Bool("changed", false, "show only benchmarks that have changed")
	magSort     = flag.Bool("mag", false, "sort benchmarks by magnitude of change")
	filename = flag.String("f", "", "bench result filename")
)

func main() {
	flag.Parse()

	var file *os.File

	if *filename == "" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(*filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	}

	benchmarks, err := Parse(file)
	if err != nil {
		panic(err)
	}

	cmps, err := Correlate(benchmarks, flag.Arg(0), flag.Arg(1))
	if err != nil {
		panic(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 5, ' ', 0)
	defer w.Flush()

	var header bool // Has the header has been displayed yet for a given block?

	if *magSort {
		sort.Sort(ByDeltaNsPerOp(cmps))
	} else {
		sort.Sort(ByParseOrder(cmps))
	}
	for _, cmp := range cmps {
		if !cmp.Measured(parse.NsPerOp) {
			continue
		}
		if delta := cmp.DeltaNsPerOp(); !*changedOnly || delta.Changed() {
			if !header {
				fmt.Fprint(w, "benchmark\told ns/op\tnew ns/op\tdelta\n")
				header = true
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", cmp.Name(), formatNs(cmp.Before.NsPerOp), formatNs(cmp.After.NsPerOp), delta.Percent())
		}
	}

	header = false
	if *magSort {
		sort.Sort(ByDeltaMBPerS(cmps))
	}
	for _, cmp := range cmps {
		if !cmp.Measured(parse.MBPerS) {
			continue
		}
		if delta := cmp.DeltaMBPerS(); !*changedOnly || delta.Changed() {
			if !header {
				fmt.Fprint(w, "\nbenchmark\told MB/s\tnew MB/s\tspeedup\n")
				header = true
			}
			fmt.Fprintf(w, "%s\t%.2f\t%.2f\t%s\n", cmp.Name(), cmp.Before.MBPerS, cmp.After.MBPerS, delta.Multiple())
		}
	}

	header = false
	if *magSort {
		sort.Sort(ByDeltaAllocsPerOp(cmps))
	}
	for _, cmp := range cmps {
		if !cmp.Measured(parse.AllocsPerOp) {
			continue
		}
		if delta := cmp.DeltaAllocsPerOp(); !*changedOnly || delta.Changed() {
			if !header {
				fmt.Fprint(w, "\nbenchmark\told allocs\tnew allocs\tdelta\n")
				header = true
			}
			fmt.Fprintf(w, "%s\t%d\t%d\t%s\n", cmp.Name(), cmp.Before.AllocsPerOp, cmp.After.AllocsPerOp, delta.Percent())
		}
	}

	header = false
	if *magSort {
		sort.Sort(ByDeltaAllocedBytesPerOp(cmps))
	}
	for _, cmp := range cmps {
		if !cmp.Measured(parse.AllocedBytesPerOp) {
			continue
		}
		if delta := cmp.DeltaAllocedBytesPerOp(); !*changedOnly || delta.Changed() {
			if !header {
				fmt.Fprint(w, "\nbenchmark\told bytes\tnew bytes\tdelta\n")
				header = true
			}
			fmt.Fprintf(w, "%s\t%d\t%d\t%s\n", cmp.Name(), cmp.Before.AllocedBytesPerOp, cmp.After.AllocedBytesPerOp, cmp.DeltaAllocedBytesPerOp().Percent())
		}
	}
}

func Parse(file *os.File) ([]*parse.Benchmark, error) {
	scanner := bufio.NewScanner(file)

	bb := []*parse.Benchmark{}

	for scanner.Scan() {
		line := scanner.Text()

		b, err := parse.ParseLine(line)
		if err == nil {
			bb = append(bb, b)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Errorf("Invalid input: %s", err)
		return nil, err
	}

	return bb, nil
}

// formatNs formats ns measurements to expose a useful amount of
// precision. It mirrors the ns precision logic of testing.B.
func formatNs(ns float64) string {
	prec := 0
	switch {
	case ns < 10:
		prec = 2
	case ns < 100:
		prec = 1
	}
	return strconv.FormatFloat(ns, 'f', prec, 64)
}