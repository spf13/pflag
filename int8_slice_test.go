// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpI8SFlagSet(isp *[]int8) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Int8SliceVar(isp, "is", []int8{}, "Command separated list!")
	return f
}

func setUpI8SFlagSetWithDefault(isp *[]int8) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Int8SliceVar(isp, "is", []int8{0, 1}, "Command separated list!")
	return f
}

func TestEmptyI8S(t *testing.T) {
	var is []int8
	f := setUpI8SFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getI8S, err := f.GetInt8Slice("is")
	if err != nil {
		t.Fatal("got an error from GetInt8Slice():", err)
	}
	if len(getI8S) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getI8S, len(getI8S))
	}
}

func TestI8S(t *testing.T) {
	var is []int8
	f := setUpI8SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getI8S, err := f.GetInt8Slice("is")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getI8S {
		d64, err := strconv.ParseInt(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d from GetInt8Slice", i, vals[i], v)
		}
	}
}

func TestI8SDefault(t *testing.T) {
	var is []int8
	f := setUpI8SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getI8S, err := f.GetInt8Slice("is")
	if err != nil {
		t.Fatal("got an error from GetInt8Slice():", err)
	}
	for i, v := range getI8S {
		d64, err := strconv.ParseInt(vals[i], 0, 8)
		if err != nil {
			t.Fatal("got an error from GetInt8Slice():", err)
		}
		d := int8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetInt8Slice but got: %d", i, d, v)
		}
	}
}

func TestI8SWithDefault(t *testing.T) {
	var is []int8
	f := setUpI8SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseInt(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getI8S, err := f.GetInt8Slice("is")
	if err != nil {
		t.Fatal("got an error from GetInt8Slice():", err)
	}
	for i, v := range getI8S {
		d64, err := strconv.ParseInt(vals[i], 0, 8)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := int8(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetInt8Slice but got: %d", i, d, v)
		}
	}
}

func TestI8SAsSliceValue(t *testing.T) {
	var i8s []int8
	f := setUpI8SFlagSet(&i8s)

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

func TestI8SCalledTwice(t *testing.T) {
	var is []int8
	f := setUpI8SFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []int8{1, 2, 3}
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
