// Copyright 2009 The Go Authors. All rights reserved.
// Use of ths2i source code s2i governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func setUpS2IFlagSet(s2ip *map[string]int) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.StringToIntVar(s2ip, "s2i", map[string]int{}, "Command separated ls2it!")
	return f
}

func setUpS2IFlagSetWithDefault(s2ip *map[string]int) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.StringToIntVar(s2ip, "s2i", map[string]int{"a": 1, "b": 2}, "Command separated ls2it!")
	return f
}

func createS2IFlag(vals map[string]int) string {
	var buf bytes.Buffer
	i := 0
	for k, v := range vals {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.Itoa(v))
		i++
	}
	return buf.String()
}

func TestEmptyS2I(t *testing.T) {
	var s2i map[string]int
	f := setUpS2IFlagSet(&s2i)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getS2I, err := f.GetStringToInt("s2i")
	if err != nil {
		t.Fatal("got an error from GetStringToInt():", err)
	}
	if len(getS2I) != 0 {
		t.Fatalf("got s2i %v with len=%d but expected length=0", getS2I, len(getS2I))
	}
}

func TestS2I(t *testing.T) {
	var s2i map[string]int
	f := setUpS2IFlagSet(&s2i)

	vals := map[string]int{"a": 1, "b": 2, "d": 4, "c": 3}
	arg := fmt.Sprintf("--s2i=%s", createS2IFlag(vals))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2i {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", k, vals[k], v)
		}
	}
	getS2I, err := f.GetStringToInt("s2i")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for k, v := range getS2I {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d from GetStringToInt", k, vals[k], v)
		}
	}
}

func TestS2IDefault(t *testing.T) {
	var s2i map[string]int
	f := setUpS2IFlagSetWithDefault(&s2i)

	vals := map[string]int{"a": 1, "b": 2}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2i {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", k, vals[k], v)
		}
	}

	getS2I, err := f.GetStringToInt("s2i")
	if err != nil {
		t.Fatal("got an error from GetStringToInt():", err)
	}
	for k, v := range getS2I {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d from GetStringToInt but got: %d", k, vals[k], v)
		}
	}
}

func TestS2IWithDefault(t *testing.T) {
	var s2i map[string]int
	f := setUpS2IFlagSetWithDefault(&s2i)

	vals := map[string]int{"a": 1, "b": 2}
	arg := fmt.Sprintf("--s2i=%s", createS2IFlag(vals))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2i {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", k, vals[k], v)
		}
	}

	getS2I, err := f.GetStringToInt("s2i")
	if err != nil {
		t.Fatal("got an error from GetStringToInt():", err)
	}
	for k, v := range getS2I {
		if vals[k] != v {
			t.Fatalf("expected s2i[%s] to be %d from GetStringToInt but got: %d", k, vals[k], v)
		}
	}
}

func TestS2ICalledTwice(t *testing.T) {
	var s2i map[string]int
	f := setUpS2IFlagSet(&s2i)

	in := []string{"a=1,b=2", "b=3"}
	expected := map[string]int{"a": 1, "b": 3}
	argfmt := "--s2i=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range s2i {
		if expected[i] != v {
			t.Fatalf("expected s2i[%s] to be %d but got: %d", i, expected[i], v)
		}
	}
}

func TestS2IPreservesKeys(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	tests := []struct {
		name     string
		defaults map[string]int
		sets     []string
		want     map[string]int
	}{
		{name: "opening bracket", want: map[string]int{"[name": -1}},
		{name: "closing bracket", want: map[string]int{"]name": 0}},
		{name: "bracketed", want: map[string]int{"[name]": 1}},
		{name: "repeated brackets", want: map[string]int{"[[name]]": maxInt}},
		{name: "only brackets", want: map[string]int{"[]][": -maxInt - 1}},
		{name: "empty key", want: map[string]int{"": 0}},
		{name: "normal key", want: map[string]int{"name": -1}},
		{name: "multiple entries", want: map[string]int{"[min]": -maxInt - 1, "]max[": maxInt, "[[zero]]": 0}},
		{
			name:     "repeated set",
			defaults: map[string]int{"default": 42},
			sets:     []string{"[name]=-1,]=2", "[name]=3,[[new]]=0"},
			want:     map[string]int{"[name]": 3, "]": 2, "[[new]]": 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []string{"default", "set"} {
				t.Run(mode, func(t *testing.T) {
					defaults := tt.want
					sets := tt.sets
					if sets != nil {
						defaults = tt.defaults
					} else if mode == "set" {
						defaults = map[string]int{"default": 42}
						sets = []string{createS2IFlag(tt.want)}
					}
					f := NewFlagSet("test", ContinueOnError)
					f.StringToInt("s2i", defaults, "")
					for _, value := range sets {
						if err := f.Set("s2i", value); err != nil {
							t.Fatal(err)
						}
					}
					got, err := f.GetStringToInt("s2i")
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tt.want) {
						t.Fatalf("expected %v, got %v", tt.want, got)
					}
				})
			}
		})
	}
}
