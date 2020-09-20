// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpUI64SFlagSet(isp *[]uint64) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint64SliceVar(isp, "is", []uint64{}, "Command separated list!")
	return f
}

func setUpUI64SFlagSetWithDefault(isp *[]uint64) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint64SliceVar(isp, "is", []uint64{0, 1}, "Command separated list!")
	return f
}

func TestEmptyUI64S(t *testing.T) {
	var is []uint64
	f := setUpUI64SFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getI64S, err := f.GetUint64Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint64Slice():", err)
	}
	if len(getI64S) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getI64S, len(getI64S))
	}
}

func TestUI64S(t *testing.T) {
	var is []uint64
	f := setUpUI64SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.ParseUint(vals[i], 0, 64)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getI64S, err := f.GetUint64Slice("is")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getI64S {
		d, err := strconv.ParseUint(vals[i], 0, 64)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d from GetUint64Slice", i, vals[i], v)
		}
	}
}

func TestUI64SDefault(t *testing.T) {
	var is []uint64
	f := setUpUI64SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.ParseUint(vals[i], 0, 64)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getI64S, err := f.GetUint64Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint64Slice():", err)
	}
	for i, v := range getI64S {
		d, err := strconv.ParseUint(vals[i], 0, 64)
		if err != nil {
			t.Fatal("got an error from GetUint64Slice():", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint64Slice but got: %d", i, d, v)
		}
	}
}

func TestIU64SWithDefault(t *testing.T) {
	var is []uint64
	f := setUpUI64SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.ParseUint(vals[i], 0, 64)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getI64S, err := f.GetUint64Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint64Slice():", err)
	}
	for i, v := range getI64S {
		d, err := strconv.ParseUint(vals[i], 0, 64)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint64Slice but got: %d", i, d, v)
		}
	}
}

func TestUI64SAsSliceValue(t *testing.T) {
	var i64s []uint64
	f := setUpUI64SFlagSet(&i64s)

	in := []string{"1", "2"}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	f.VisitAll(func(f *Flag) {
		if val, ok := f.Value.(SliceValue); ok {
			_ = val.Replace([]string{"3"})
		}
	})
	if len(i64s) != 1 || i64s[0] != 3 {
		t.Fatalf("Expected ss to be overwritten with '3.1', but got: %v", i64s)
	}
}

func TestUI64SCalledTwice(t *testing.T) {
	var is []uint64
	f := setUpUI64SFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []uint64{1, 2, 3}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		if expected[i] != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, expected[i], v)
		}
	}
}
