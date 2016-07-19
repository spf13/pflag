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

func setUpBoolFlagSet(isp *[]bool) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(isp, "is", []bool{}, "Command separated list!")
	return f
}

func setUpBoolFlagSetWithDefault(isp *[]bool) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(isp, "is", []bool{true, false}, "Command separated list!")
	return f
}

func TestEmptyBool(t *testing.T) {
	var is []bool
	f := setUpBoolFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getIS, err := f.GetBoolSlice("is")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	if len(getIS) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getIS, len(getIS))
	}
}

func TestBool(t *testing.T) {
	var is []bool
	f := setUpBoolFlagSet(&is)

	vals := []string{"true", "false", "false", "true"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected bool[%d] to be %s but got: %t", i, vals[i], v)
		}
	}
	getIS, err := f.GetBoolSlice("is")
	for i, v := range getIS {
		d, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %t from GetBoolSlice", i, vals[i], v)
		}
	}
}

func TestBoolDefault(t *testing.T) {
	var is []bool
	f := setUpBoolFlagSetWithDefault(&is)

	vals := []string{"true", "false"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %t but got: %t", i, d, v)
		}
	}

	getIS, err := f.GetBoolSlice("is")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	for i, v := range getIS {
		d, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatal("got an error from GetBoolSlice():", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %t from GetBoolSlice but got: %t", i, d, v)
		}
	}
}

func TestBoolWithDefault(t *testing.T) {
	var is []bool
	f := setUpBoolFlagSetWithDefault(&is)

	vals := []string{"true", "false"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %t but got: %t", i, d, v)
		}
	}

	getIS, err := f.GetBoolSlice("is")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	for i, v := range getIS {
		d, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if d != v {
			t.Fatalf("expected is[%d] to be %t from GetBoolSlice but got: %t", i, d, v)
		}
	}
}

func TestBoolCalledTwice(t *testing.T) {
	var is []bool
	f := setUpBoolFlagSet(&is)

	in := []string{"true,false", "true"}
	expected := []bool{true, false, true}
	argfmt := "--is=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		if expected[i] != v {
			t.Fatalf("expected is[%d] to be %t but got: %t", i, expected[i], v)
		}
	}
}
