package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpUI16SFlagSet(isp *[]uint16) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint16SliceVar(isp, "is", []uint16{}, "Command separated list!")
	return f
}

func setUpUI16SFlagSetWithDefault(isp *[]uint16) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.Uint16SliceVar(isp, "is", []uint16{0, 1}, "Command separated list!")
	return f
}

func TestEmptyUI16S(t *testing.T) {
	var is []uint16
	f := setUpUI16SFlagSet(&is)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getUI16S, err := f.GetUint16Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint16Slice():", err)
	}
	if len(getUI16S) != 0 {
		t.Fatalf("got is %v with len=%d but expected length=0", getUI16S, len(getUI16S))
	}
}

func TestUI16S(t *testing.T) {
	var is []uint16
	f := setUpUI16SFlagSet(&is)

	vals := []string{"1", "2", "4", "3"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d", i, vals[i], v)
		}
	}
	getUI16S, err := f.GetUint16Slice("is")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getUI16S {
		d64, err := strconv.ParseUint(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %s but got: %d from GetUint16Slice", i, vals[i], v)
		}
	}
}

func TestUI16SDefault(t *testing.T) {
	var is []uint16
	f := setUpUI16SFlagSetWithDefault(&is)

	vals := []string{"0", "1"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getUI16S, err := f.GetUint16Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint16Slice():", err)
	}
	for i, v := range getUI16S {
		d64, err := strconv.ParseUint(vals[i], 0, 16)
		if err != nil {
			t.Fatal("got an error from GetUint16Slice():", err)
		}
		d := uint16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint16Slice but got: %d", i, d, v)
		}
	}
}

func TestUI16SWithDefault(t *testing.T) {
	var is []uint16
	f := setUpUI16SFlagSetWithDefault(&is)

	vals := []string{"1", "2"}
	arg := fmt.Sprintf("--is=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range is {
		d64, err := strconv.ParseUint(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d but got: %d", i, d, v)
		}
	}

	getUI16S, err := f.GetUint16Slice("is")
	if err != nil {
		t.Fatal("got an error from GetUint16Slice():", err)
	}
	for i, v := range getUI16S {
		d64, err := strconv.ParseUint(vals[i], 0, 16)
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		d := uint16(d64)
		if d != v {
			t.Fatalf("expected is[%d] to be %d from GetUint16Slice but got: %d", i, d, v)
		}
	}
}

func TestUI16SAsSliceValue(t *testing.T) {
	var i16s []uint16
	f := setUpUI16SFlagSet(&i16s)

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
	if len(i16s) != 1 || i16s[0] != 3 {
		t.Fatalf("Expected ss to be overwritten with '3.1', but got: %v", i16s)
	}
}

func TestUI16SCalledTwice(t *testing.T) {
	var is []uint16
	f := setUpUI16SFlagSet(&is)

	in := []string{"1,2", "3"}
	expected := []uint16{1, 2, 3}
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
