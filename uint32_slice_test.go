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

func setUpUI32SFlagSet(isp *[]uint32) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint32SliceVar(isp, "is", []uint32{}, "Command separated list!")
	return f
}

func setUpUI32SFlagSetWithDefault(isp *[]uint32) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint32SliceVar(isp, "is", []uint32{0, 1}, "Command separated list!")
	return f
}

func TestEmptyUI32S(t *testing.T) {
	var is []uint32
	f := setUpUI32SFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getUI32S, err := f.GetUint32Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint32Slice():", err)
	}
	if len(getUI32S) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getUI32S, len(getUI32S))
	}
}

func TestUI32S(t *testing.T) {
	var is []uint32
	f := setUpUI32SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 32)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint32(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getUI32S, err := f.GetUint32Slice("is")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getUI32S {
		d64, err := strconv.ParseUint(vals[i], 0, 32)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint32(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d from GetUint32Slice", i, vals[i], v)
		}
	}
}

func TestUI32SDefault(t *testing.T) {
	var is []uint32
	f := setUpUI32SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 32)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint32(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getUI32S, err := f.GetUint32Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint32Slice():", err)
	}
	for i, v := range getUI32S {
		d64, err := strconv.ParseUint(vals[i], 0, 32)
		if err != nil {
			t.Fatal("got an error from GetUint32Slice():", err)
		}
		d := uint32(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint32Slice but got: %d", i, d, v)
		}
	}
}

func TestUI32SWithDefault(t *testing.T) {
	var is []uint32
	f := setUpUI32SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 32)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint32(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getUI32S, err := f.GetUint32Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint32Slice():", err)
	}
	for i, v := range getUI32S {
		d64, err := strconv.ParseUint(vals[i], 0, 32)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint32(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint32Slice but got: %d", i, d, v)
		}
	}
}

func TestUI32SAsSliceValue(t *testing.T) {
	var i32s []uint32
	f := setUpUI32SFlagSet(&i32s)

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
	if len(i32s) != 1 || i32s[0] != 3 {
		t.Fatalf("Expected ss to be overwritten with '3.1', but got: %v", i32s)
	}
}

func TestUI32SCalledTwice(t *testing.T) {
	var is []uint32
	f := setUpUI32SFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []uint32{1, 2, 3}
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
