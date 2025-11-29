// Copyright 2009 The Go Authors. All rights reserved.
// Use of ths2s source code s2s governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestStringToString(t *testing.T) {
	tt := []struct {
		args     []string
		def      map[string]string
		expected map[string]string
	}{
		{
			// should permit no args and defaults
			args:     []string{},
			def:      map[string]string{},
			expected: map[string]string{},
		},
		{
			// should use defaults when no args given
			args:     []string{},
			def:      map[string]string{"a": "1", "b": "2"},
			expected: map[string]string{"a": "1", "b": "2"},
		},
		{
			// should parse single key-value pair
			args:     []string{"--arg", "a=1"},
			def:      map[string]string{},
			expected: map[string]string{"a": "1"},
		},
		{
			// should allow comma-separated key-value pairs
			args:     []string{"--arg", "a=1,b=2"},
			def:      map[string]string{},
			expected: map[string]string{"a": "1", "b": "2"},
		},
		{
			// should correctly parse values with commas
			args:     []string{"--arg", "a=1,2"},
			def:      map[string]string{},
			expected: map[string]string{"a": "1,2"},
		},
		{
			// should correctly parse values with equal symbols
			args:     []string{"--arg", "a=1="},
			def:      map[string]string{},
			expected: map[string]string{"a": "1="},
		},
		{
			// should allow multiple map args, merging into a single result
			args:     []string{"--arg", "a=1,b=2", "--arg", "c=3", "--arg", "a=2"},
			def:      map[string]string{},
			expected: map[string]string{"a": "2", "b": "2", "c": "3"},
		},
		{
			// should ensure command-line args take precedence over defaults
			args:     []string{"--arg", "a=4"},
			def:      map[string]string{"a": "1", "b": "2"},
			expected: map[string]string{"a": "4"},
		},
		{
			// should allow quoting of values to handle values with '=' and ','
			args:     []string{"--arg", `"foo=bar,bar=qix",qix=foo`},
			def:      map[string]string{},
			expected: map[string]string{"foo": "bar,bar=qix", "qix": "foo"},
		},
		{
			// should allow quoting of values to handle values with '=' and ','
			args:     []string{"--arg", `"foo=bar,bar=qix"`, "--arg", "qix=foo"},
			def:      map[string]string{},
			expected: map[string]string{"foo": "bar,bar=qix", "qix": "foo"},
		},
		{
			// should allow stuck values
			args:     []string{`--arg="e=5,6",a=1,b=2,d=4,c=3`},
			def:      map[string]string{},
			expected: map[string]string{"a": "1", "b": "2", "d": "4", "c": "3", "e": "5,6"},
		},
		{
			// should allow stuck values with defaults
			args:     []string{`--arg=a=1,b=2,"e=5,6"`},
			def:      map[string]string{"da": "1", "db": "2", "de": "5,6"},
			expected: map[string]string{"a": "1", "b": "2", "e": "5,6"},
		},
		{
			// should allow multiple stuck value args
			args:     []string{"--arg=a=1,b=2", "--arg=b=3", `--arg="e=5,6"`, `--arg=f=7,8`},
			def:      map[string]string{},
			expected: map[string]string{"a": "1", "b": "3", "e": "5,6", "f": "7,8"},
		},
		{
			// should parse arg with empty key and value
			args:     []string{"--arg", "="},
			def:      map[string]string{},
			expected: map[string]string{"": ""},
		},
		{
			// should parse comma delimited empty mappings
			args:     []string{"--arg", "=,=,="},
			def:      map[string]string{},
			expected: map[string]string{"": ""},
		},
		{
			// should peremit overlapping mappings
			args:     []string{"--arg", "a=1,a=2"},
			def:      map[string]string{},
			expected: map[string]string{"a": "2"},
		},
		{
			// should correctly parse short args
			args:     []string{"-a", "a=1,b=2", "-a=c=3"},
			def:      map[string]string{},
			expected: map[string]string{"a": "1", "b": "2", "c": "3"},
		},
	}

	for num, test := range tt {
		t.Logf("=== TEST %d ===", num)
		t.Logf("    Args:          %v", test.args)
		t.Logf("    Default Value: %v", test.def)
		t.Logf("    Expected:      %v", test.expected)

		f := NewFlagSet("test", ContinueOnError)
		f.StringToStringP("arg", "a", test.def, "test string-to-string arg")

		if err := f.Parse(test.args); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		result, err := f.GetStringToString("arg")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		t.Logf("    Actual:        %v", result)

		for k, v := range test.expected {
			actual, ok := result[k]
			if !ok {
				t.Fatalf("missing key in result: %s", k)
			}
			if actual != v {
				t.Fatalf("unexpected value in result for key '%s': %s", k, actual)
			}
		}

		if len(test.expected) != len(result) {
			t.Fatalf("unexpected extra key-value pairs in result: %v", result)
		}
	}
}

// This test ensures that [FlagSet.GetStringToString] always return the pointers which were given during flag
// initialization.
//
// This behaviour is important as it ensures consumers of the library can access the underlying map in a stable,
// consistent manner.
func TestS2SStablePointers(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)

	defval := map[string]string{"a": "1", "b": "2"}

	ptr := f.StringToString("map-flag", defval, "test for s2s arg")

	if reflect.ValueOf(*ptr).Pointer() != reflect.ValueOf(defval).Pointer() {
		t.Fatal("pointer mismatch")
	}

	// initially, arg should have defaults
	result0, err := f.GetStringToString("map-flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := result0["a"]; !ok || v != "1" {
		t.Fatalf("value not present in map or unexpected value: %v", result0)
	}
	if v, ok := result0["b"]; !ok || v != "2" {
		t.Fatalf("value not present in map or unexpected value: %v", result0)
	}

	if reflect.ValueOf(result0).Pointer() != reflect.ValueOf(defval).Pointer() {
		t.Fatal("pointer mismatch")
	}
	if reflect.ValueOf(*ptr).Pointer() != reflect.ValueOf(result0).Pointer() {
		t.Fatal("pointer mismatch")
	}

	// manipulate the map; the map should now have a single mapping and the pointers should be stable
	if err := f.Set("map-flag", "c=3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result1, err := f.GetStringToString("map-flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reflect.ValueOf(*ptr).Pointer() != reflect.ValueOf(result1).Pointer() {
		t.Fatal("pointer mismatch")
	}
	if reflect.ValueOf(result0).Pointer() != reflect.ValueOf(result1).Pointer() {
		t.Fatal("pointer mismatch")
	}

	// manipulate the map once more
	if err := f.Set("map-flag", "d=4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result2, err := f.GetStringToString("map-flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if reflect.ValueOf(*ptr).Pointer() != reflect.ValueOf(result2).Pointer() {
		t.Fatal("pointer mismatch")
	}
	if reflect.ValueOf(result1).Pointer() != reflect.ValueOf(result2).Pointer() {
		t.Fatal("pointer mismatch")
	}

	// check that the newly added flag value was updated
	if v, ok := result1["c"]; !ok || v != "3" {
		t.Fatalf("value not present in map or unexpected value: %v", result1)
	}
	if v, ok := result1["d"]; !ok || v != "4" {
		t.Fatalf("value not present in map or unexpected value: %v", result1)
	}

	// finally, if we clear the map, it should reset flag
	for k := range result1 {
		delete(result1, k)
	}

	result3, err := f.GetStringToString("map-flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result3) != 0 {
		t.Fatalf("unexpected map values: %v", result3)
	}
}

func setUpS2SFlagSet(s2sp *map[string]string) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.StringToStringVar(s2sp, "s2s", map[string]string{}, "Command separated ls2st!")
	return f
}

func setUpS2SFlagSetWithDefault(s2sp *map[string]string) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.StringToStringVar(s2sp, "s2s", map[string]string{"da": "1", "db": "2", "de": "5,6"}, "Command separated ls2st!")
	return f
}

func createS2SFlag(vals map[string]string) string {
	records := make([]string, 0, len(vals)>>1)
	for k, v := range vals {
		records = append(records, k+"="+v)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(records); err != nil {
		panic(err)
	}
	w.Flush()
	return strings.TrimSpace(buf.String())
}

func TestEmptyS2S(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	if len(getS2S) != 0 {
		t.Fatalf("got s2s %v with len=%d but expected length=0", getS2S, len(getS2S))
	}
}

func TestS2S(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)

	vals := map[string]string{"a": "1", "b": "2", "d": "4", "c": "3", "e": "5,6"}
	arg := fmt.Sprintf("--s2s=%s", createS2SFlag(vals))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", k, vals[k], v)
		}
	}
	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s from GetStringToString", k, vals[k], v)
		}
	}
}

func TestS2SDefault(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSetWithDefault(&s2s)

	vals := map[string]string{"da": "1", "db": "2", "de": "5,6"}

	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", k, vals[k], v)
		}
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s from GetStringToString but got: %s", k, vals[k], v)
		}
	}
}

func TestS2SWithDefault(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSetWithDefault(&s2s)

	vals := map[string]string{"a": "1", "b": "2", "e": "5,6"}
	arg := fmt.Sprintf("--s2s=%s", createS2SFlag(vals))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for k, v := range s2s {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", k, vals[k], v)
		}
	}

	getS2S, err := f.GetStringToString("s2s")
	if err != nil {
		t.Fatal("got an error from GetStringToString():", err)
	}
	for k, v := range getS2S {
		if vals[k] != v {
			t.Fatalf("expected s2s[%s] to be %s from GetStringToString but got: %s", k, vals[k], v)
		}
	}
}

func TestS2SCalledTwice(t *testing.T) {
	var s2s map[string]string
	f := setUpS2SFlagSet(&s2s)

	in := []string{"a=1,b=2", "b=3", `"e=5,6"`, `f=7,8`}
	expected := map[string]string{"a": "1", "b": "3", "e": "5,6", "f": "7,8"}
	argfmt := "--s2s=%s"
	arg0 := fmt.Sprintf(argfmt, in[0])
	arg1 := fmt.Sprintf(argfmt, in[1])
	arg2 := fmt.Sprintf(argfmt, in[2])
	arg3 := fmt.Sprintf(argfmt, in[3])
	err := f.Parse([]string{arg0, arg1, arg2, arg3})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	if len(s2s) != len(expected) {
		t.Fatalf("expected %d flags; got %d flags", len(expected), len(s2s))
	}
	for i, v := range s2s {
		if expected[i] != v {
			t.Fatalf("expected s2s[%s] to be %s but got: %s", i, expected[i], v)
		}
	}
}
