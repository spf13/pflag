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

func setUpUI8SFlagSet(isp *[]uint8) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint8SliceVar(isp, "is", []uint8{}, "Command separated list!")
	return f
}

func setUpUI8SFlagSetWithDefault(isp *[]uint8) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint8SliceVar(isp, "is", []uint8{0, 1}, "Command separated list!")
	return f
}

func TestEmptyUI8S(t *testing.T) {
	var is []uint8
	f := setUpUI8SFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getUI8S, err := f.GetUint8Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint8Slice():", err)
	}
	if len(getUI8S) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getUI8S, len(getUI8S))
	}
}

func TestUI8S(t *testing.T) {
	var is []uint8
	f := setUpUI8SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getUI8S, err := f.GetUint8Slice("is")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getUI8S {
		d64, err := strconv.ParseUint(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d from GetUint8Slice", i, vals[i], v)
		}
	}
}

func TestUI8SDefault(t *testing.T) {
	var is []uint8
	f := setUpUI8SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getUI8S, err := f.GetUint8Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint8Slice():", err)
	}
	for i, v := range getUI8S {
		d64, err := strconv.ParseUint(vals[i], 0, 8)
		if err != nil {
			t.Fatal("got an error from GetUint8Slice():", err)
		}
		d := uint8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint8Slice but got: %d", i, d, v)
		}
	}
}

func TestUI8SWithDefault(t *testing.T) {
	var is []uint8
	f := setUpUI8SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getUI8S, err := f.GetUint8Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint8Slice():", err)
	}
	for i, v := range getUI8S {
		d64, err := strconv.ParseUint(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint8Slice but got: %d", i, d, v)
		}
	}
}

func TestUI8SAsSliceValue(t *testing.T) {
	var i8s []uint8
	f := setUpUI8SFlagSet(&i8s)

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
	if len(i8s) != 1 || i8s[0] != 3 {
		t.Fatalf("Expected ss to be overwritten with '3.1', but got: %v", i8s)
	}
}

func TestUI8SCalledTwice(t *testing.T) {
	var is []uint8
	f := setUpUI8SFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []uint8{1, 2, 3}
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
