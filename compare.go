// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"

	"golang.org/x/tools/benchmark/parse"
	"regexp"
	"errors"
	"strings"
)

// BenchCmp is a pair of benchmarks.
type BenchCmp struct {
	Before *parse.Benchmark
	After  *parse.Benchmark
}

// Correlate correlates benchmarks from two BenchSets.
func Correlate(benchmarks []*parse.Benchmark, before, after string) ([]BenchCmp, error) {
	cmps := []BenchCmp{}
	beforeRe := regexp.MustCompile(before)
	beforeList := []*parse.Benchmark{}
	beforeHasMain := false
	afterRe := regexp.MustCompile(after)
	afterList := []*parse.Benchmark{}
	afterHasMain := false
	for _, b := range benchmarks {
		if beforeRe.FindStringIndex(b.Name) != nil {
			beforeList = append(beforeList, b)
			if strings.Index(b.Name, "/") < 0 {
				beforeHasMain = true
			}
		}
		if afterRe.FindStringIndex(b.Name) != nil {
			afterList = append(afterList, b)
			if strings.Index(b.Name, "/") < 0 {
				afterHasMain = true
			}
		}
	}
	if len(beforeList) == 0 {
		return nil, errors.New("before benchmark not found.")
	}
	if len(afterList) == 0 {
		return nil, errors.New("after benchmark not found.")
	}
	if len(beforeList) != 1 && beforeHasMain {
		return nil, fmt.Errorf("before benchmark is ambiguous: %v", beforeList)
	}
	if len(afterList) != 1 && afterHasMain {
		return nil, fmt.Errorf("after benchmark is ambiguous: %v", afterList)
	}
	if len(beforeList) != len(afterList) {
		return nil, fmt.Errorf("before and after benchmark length are not aligned\nbefore: %v\nafter: %v", beforeList, afterList)
	}

	if len(beforeList) == 1 {
		cmps = append(cmps, BenchCmp{beforeList[0], afterList[0]})
		return cmps, nil
	}

	for _, b := range beforeList {
		s := strings.SplitN(b.Name, "/", 2)
		for _, a := range afterList {
			if strings.Index(a.Name, s[1]) != -1 {
				cmps = append(cmps, BenchCmp{b, a})
			}
		}
	}

	if len(cmps) != len(beforeList) {
		return nil, fmt.Errorf("before and after benchmark length are not aligned\nbefore: %v\nafter: %v", beforeList, afterList)
	}

	return cmps, nil
}

func (c BenchCmp) Name() string           { return c.Before.Name }
func (c BenchCmp) String() string         { return fmt.Sprintf("<%s, %s>", c.Before, c.After) }
func (c BenchCmp) Measured(flag int) bool { return (c.Before.Measured & c.After.Measured & flag) != 0 }
func (c BenchCmp) DeltaNsPerOp() Delta    { return Delta{c.Before.NsPerOp, c.After.NsPerOp} }
func (c BenchCmp) DeltaMBPerS() Delta     { return Delta{c.Before.MBPerS, c.After.MBPerS} }
func (c BenchCmp) DeltaAllocedBytesPerOp() Delta {
	return Delta{float64(c.Before.AllocedBytesPerOp), float64(c.After.AllocedBytesPerOp)}
}
func (c BenchCmp) DeltaAllocsPerOp() Delta {
	return Delta{float64(c.Before.AllocsPerOp), float64(c.After.AllocsPerOp)}
}

// Delta is the before and after value for a benchmark measurement.
// Both must be non-negative.
type Delta struct {
	Before float64
	After  float64
}

// mag calculates the magnitude of a change, regardless of the direction of
// the change. mag is intended for sorting and has no independent meaning.
func (d Delta) mag() float64 {
	switch {
	case d.Before != 0 && d.After != 0 && d.Before >= d.After:
		return d.After / d.Before
	case d.Before != 0 && d.After != 0 && d.Before < d.After:
		return d.Before / d.After
	case d.Before == 0 && d.After == 0:
		return 1
	default:
		// 0 -> 1 or 1 -> 0
		// These are significant changes and worth surfacing.
		return math.Inf(1)
	}
}

// Changed reports whether the benchmark quantities are different.
func (d Delta) Changed() bool { return d.Before != d.After }

// Float64 returns After / Before. If Before is 0, Float64 returns
// 1 if After is also 0, and +Inf otherwise.
func (d Delta) Float64() float64 {
	switch {
	case d.Before != 0:
		return d.After / d.Before
	case d.After == 0:
		return 1
	default:
		return math.Inf(1)
	}
}

// Percent formats a Delta as a percent change, ranging from -100% up.
func (d Delta) Percent() string {
	return fmt.Sprintf("%+.2f%%", 100*d.Float64()-100)
}

// Multiple formats a Delta as a multiplier, ranging from 0.00x up.
func (d Delta) Multiple() string {
	return fmt.Sprintf("%.2fx", d.Float64())
}

func (d Delta) String() string {
	return fmt.Sprintf("Δ(%f, %f)", d.Before, d.After)
}

// ByParseOrder sorts BenchCmps to match the order in
// which the Before benchmarks were presented to Parse.
type ByParseOrder []BenchCmp

func (x ByParseOrder) Len() int           { return len(x) }
func (x ByParseOrder) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x ByParseOrder) Less(i, j int) bool { return x[i].Before.Ord < x[j].Before.Ord }

// lessByDelta provides lexicographic ordering:
//   * largest delta by magnitude
//   * alphabetic by name
func lessByDelta(i, j BenchCmp, calcDelta func(BenchCmp) Delta) bool {
	iDelta, jDelta := calcDelta(i).mag(), calcDelta(j).mag()
	if iDelta != jDelta {
		return iDelta < jDelta
	}
	return i.Name() < j.Name()
}

// ByDeltaNsPerOp sorts BenchCmps lexicographically by change
// in ns/op, descending, then by benchmark name.
type ByDeltaNsPerOp []BenchCmp

func (x ByDeltaNsPerOp) Len() int           { return len(x) }
func (x ByDeltaNsPerOp) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x ByDeltaNsPerOp) Less(i, j int) bool { return lessByDelta(x[i], x[j], BenchCmp.DeltaNsPerOp) }

// ByDeltaMBPerS sorts BenchCmps lexicographically by change
// in MB/s, descending, then by benchmark name.
type ByDeltaMBPerS []BenchCmp

func (x ByDeltaMBPerS) Len() int           { return len(x) }
func (x ByDeltaMBPerS) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x ByDeltaMBPerS) Less(i, j int) bool { return lessByDelta(x[i], x[j], BenchCmp.DeltaMBPerS) }

// ByDeltaAllocedBytesPerOp sorts BenchCmps lexicographically by change
// in B/op, descending, then by benchmark name.
type ByDeltaAllocedBytesPerOp []BenchCmp

func (x ByDeltaAllocedBytesPerOp) Len() int      { return len(x) }
func (x ByDeltaAllocedBytesPerOp) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x ByDeltaAllocedBytesPerOp) Less(i, j int) bool {
	return lessByDelta(x[i], x[j], BenchCmp.DeltaAllocedBytesPerOp)
}

// ByDeltaAllocsPerOp sorts BenchCmps lexicographically by change
// in allocs/op, descending, then by benchmark name.
type ByDeltaAllocsPerOp []BenchCmp

func (x ByDeltaAllocsPerOp) Len() int      { return len(x) }
func (x ByDeltaAllocsPerOp) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x ByDeltaAllocsPerOp) Less(i, j int) bool {
	return lessByDelta(x[i], x[j], BenchCmp.DeltaAllocsPerOp)
}