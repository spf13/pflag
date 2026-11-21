package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func setUpBSFlagSet(bsp *[]bool) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool{}, "Command separated list!")
	return f
}

func setUpBSFlagSetWithDefault(bsp *[]bool) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(bsp, "bs", []bool{false, true}, "Command separated list!")
	return f
}

func TestEmptyBS(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	if len(getBS) != 0 {
		t.Fatalf("got bs %v with len=%d but expected length=0", getBS, len(getBS))
	}
}

func TestBS(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	vals := []string{"1", "F", "TRUE", "0"}
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected is[%d] to be %s but got: %t", i, vals[i], v)
		}
	}
	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %s but got: %t from GetBoolSlice", i, vals[i], v)
		}
	}
}

func TestBSDefault(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string{"false", "T"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		}
	}

	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatal("got an error from GetBoolSlice():", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		}
	}
}

func TestBSWithDefault(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSetWithDefault(&bs)

	vals := []string{"FALSE", "1"}
	arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range bs {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t but got: %t", i, b, v)
		}
	}

	getBS, err := f.GetBoolSlice("bs")
	if err != nil {
		t.Fatal("got an error from GetBoolSlice():", err)
	}
	for i, v := range getBS {
		b, err := strconv.ParseBool(vals[i])
		if err != nil {
			t.Fatalf("got error: %v", err)
		}
		if b != v {
			t.Fatalf("expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v)
		}
	}
}

func TestBSCalledTwice(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	in := []string{"T,F", "T"}
	expected := []bool{true, false, true}
	argfmt := "--bs=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range bs {
		if expected[i] != v {
			t.Fatalf("expected bs[%d] to be %t but got %t", i, expected[i], v)
		}
	}
}

func TestBSAsSliceValue(t *testing.T) {
	var bs []bool
	f := setUpBSFlagSet(&bs)

	in := []string{"true", "false"}
	argfmt := "--bs=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	f.VisitAll(func(f *Flag) {
		if val, ok := f.Value.(SliceValue); ok {
			_ = val.Replace([]string{"false"})
		}
	})
	if len(bs) != 1 || bs[0] != false {
		t.Fatalf("Expected ss to be overwritten with 'false', but got: %v", bs)
	}
}

func TestBSBadQuoting(t *testing.T) {
	tests := []struct {
		Want    []bool
		FlagArg []string
	}{
		{
			Want:    []bool{true, false, true},
			FlagArg: []string{"1", "0", "true"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"True", "F"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"T", "0"},
		},
		{
			Want:    []bool{true, false},
			FlagArg: []string{"1", "0"},
		},
		{
			Want:    []bool{true, false, false},
			FlagArg: []string{"true,false", "false"},
		},
		{
			Want:    []bool{true, false, false, true, false, true, false},
			FlagArg: []string{`"true,false,false,1,0,     T"`, " false "},
		},
		{
			Want:    []bool{false, false, true, false, true, false, true},
			FlagArg: []string{`"0, False,  T,false  , true,F"`, "true"},
		},
	}

	for i, test := range tests {

		var bs []bool
		f := setUpBSFlagSet(&bs)

		if err := f.Parse([]string{fmt.Sprintf("--bs=%s", strings.Join(test.FlagArg, ","))}); err != nil {
			t.Fatalf("flag parsing failed with error: %s\nparsing:\t%#v\nwant:\t\t%#v",
				err, test.FlagArg, test.Want[i])
		}

		for j, b := range bs {
			if b != test.Want[j] {
				t.Fatalf("bad value parsed for test %d on bool %d:\nwant:\t%t\ngot:\t%t", i, j, test.Want[j], b)
			}
		}
	}
}

// TestBSTwoTokenValue verifies that a bool slice flag consumes a separate
// value token by default: "--active false" must set the flag to [false]
// and leave no positional arguments. Giving boolSliceValue an IsBoolFlag
// method made the parser take NoOptDefVal "true" instead and leave "false"
// as a positional argument, so commands like
// "kueuectl list clusterqueue --active false" selected the wrong queues.
func TestBSTwoTokenValue(t *testing.T) {
	var bs []bool

	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(&bs, "active", nil, "bool slice")
	if err := f.Parse([]string{"--active", "false"}); err != nil {
		t.Fatal("expected no error for --active false; got", err)
	}
	if len(bs) != 1 || bs[0] {
		t.Fatalf("expected [false] for --active false; got %v", bs)
	}
	if len(f.Args()) != 0 {
		t.Fatalf("expected no positional args for --active false; got %v", f.Args())
	}

	// The short form "-a false" must behave the same.
	bs = nil
	f = NewFlagSet("test", ContinueOnError)
	f.BoolSliceVarP(&bs, "active", "a", nil, "bool slice")
	if err := f.Parse([]string{"-a", "false"}); err != nil {
		t.Fatal("expected no error for -a false; got", err)
	}
	if len(bs) != 1 || bs[0] {
		t.Fatalf("expected [false] for -a false; got %v", bs)
	}
	if len(f.Args()) != 0 {
		t.Fatalf("expected no positional args for -a false; got %v", f.Args())
	}

	// The "--active=false" form is unaffected and keeps working.
	bs = nil
	f = NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(&bs, "active", nil, "bool slice")
	if err := f.Parse([]string{"--active=false"}); err != nil {
		t.Fatal("expected no error for --active=false; got", err)
	}
	if len(bs) != 1 || bs[0] {
		t.Fatalf("expected [false] for --active=false; got %v", bs)
	}
	if len(f.Args()) != 0 {
		t.Fatalf("expected no positional args for --active=false; got %v", f.Args())
	}
}

// TestBSNoOptValue verifies that a bool slice flag can be used without
// an explicit value, like a bool flag, once the caller opts in by setting
// NoOptDefVal (see #265). Bare usage is opt-in only: without it the flag
// keeps consuming a separate value token, so existing commands such as
// "--active false" are unaffected.
func TestBSNoOptValue(t *testing.T) {
	var bs []bool
	var b bool

	f := NewFlagSet("test", ContinueOnError)
	f.BoolSliceVar(&bs, "bs", nil, "bool slice")
	f.Lookup("bs").NoOptDefVal = "true"
	if err := f.Parse([]string{"--bs"}); err != nil {
		t.Fatal("expected no error for --bs; got", err)
	}
	if len(bs) != 1 || !bs[0] {
		t.Fatalf("expected [true] for --bs; got %v", bs)
	}

	bs = nil
	f = NewFlagSet("test", ContinueOnError)
	f.BoolSliceVarP(&bs, "bs", "s", nil, "bool slice")
	f.Lookup("bs").NoOptDefVal = "true"
	f.BoolVarP(&b, "bb", "b", false, "bool")
	if err := f.Parse([]string{"-sb"}); err != nil {
		t.Fatal("expected no error for -sb; got", err)
	}
	if len(bs) != 1 || !bs[0] {
		t.Fatalf("expected [true] for -sb; got %v", bs)
	}
	if !b {
		t.Fatal("expected true for -sb; got false")
	}

	bs = nil
	f = NewFlagSet("test", ContinueOnError)
	f.BoolSliceVarP(&bs, "bs", "s", nil, "bool slice")
	f.Lookup("bs").NoOptDefVal = "true"
	if err := f.Parse([]string{"-s"}); err != nil {
		t.Fatal("expected no error for -s; got", err)
	}
	if len(bs) != 1 || !bs[0] {
		t.Fatalf("expected [true] for -s; got %v", bs)
	}
}
