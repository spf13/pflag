package pflag

import (
	"reflect"
	"testing"
)

var stringGetterTests = []struct {
	name     string
	register func(*FlagSet, string, []string, string) *[]string
	get      func(*FlagSet, string) ([]string, error)
}{
	{"stringArray", (*FlagSet).StringArray, (*FlagSet).GetStringArray},
	{"stringSlice", (*FlagSet).StringSlice, (*FlagSet).GetStringSlice},
}

func TestStringGettersPreserveValues(t *testing.T) {
	for _, getter := range stringGetterTests {
		t.Run(getter.name, func(t *testing.T) {
			for _, tc := range []struct {
				name string
				want []string
			}{
				{"embedded CRLF", []string{"first\r\nsecond"}},
				{"trailing CRLF", []string{"trailing\r\n"}},
				{"mixed line endings", []string{"one\r\ntwo\nthree\rfour"}},
				{"CSV characters", []string{"a,b", "quoted \"value\"", "[brackets]", ""}},
				{"multiple empty values", []string{"", ""}},
				{"plain values", []string{"first", "second"}},
			} {
				t.Run(tc.name, func(t *testing.T) {
					f := NewFlagSet("test", ContinueOnError)
					values := getter.register(f, "value", tc.want, "")
					got, err := getter.get(f, "value")
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tc.want) {
						t.Fatalf("getter returned %q, want %q", got, tc.want)
					}
					first := (*values)[0]
					got[0] = "modified"
					if (*values)[0] != first {
						t.Fatalf("modifying the returned slice changed the flag: %q", *values)
					}
				})
			}
		})
	}
}

func TestStringArrayGetterPreservesParsedCRLF(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	values := f.StringArray("value", nil, "")
	want := []string{"first\r\nsecond", "trailing\r\n"}
	if err := f.Parse([]string{"--value", want[0], "--value", want[1]}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*values, want) {
		t.Fatalf("parsed %q, want %q", *values, want)
	}
	got, err := f.GetStringArray("value")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetStringArray() = %q, want %q", got, want)
	}
}

func TestStringGettersEmptyValues(t *testing.T) {
	for _, getter := range stringGetterTests {
		t.Run(getter.name, func(t *testing.T) {
			for _, values := range [][]string{nil, {}, {""}} {
				f := NewFlagSet("test", ContinueOnError)
				getter.register(f, "value", values, "")
				got, err := getter.get(f, "value")
				if err != nil {
					t.Fatal(err)
				}
				if got == nil || len(got) != 0 {
					t.Errorf("getter for %q returned %#v, want a non-nil empty slice", values, got)
				}
			}
		})
	}
}

type customStringGetterValue struct {
	flagType string
	text     string
}

func (v *customStringGetterValue) Set(s string) error { v.text = s; return nil }
func (v *customStringGetterValue) String() string     { return v.text }
func (v *customStringGetterValue) Type() string       { return v.flagType }

func TestStringGettersCustomValue(t *testing.T) {
	for _, getter := range stringGetterTests {
		t.Run(getter.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.Var(&customStringGetterValue{getter.name, `["first,second",third]`}, "value", "")
			got, err := getter.get(f, "value")
			if err != nil {
				t.Fatal(err)
			}
			if want := []string{"first,second", "third"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("getter returned %q, want %q", got, want)
			}
		})
	}
}

func TestStringGettersErrors(t *testing.T) {
	for _, getter := range stringGetterTests {
		t.Run(getter.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.String("wrong-type", "value", "")
			f.Var(&customStringGetterValue{getter.name, `["unclosed]`}, "invalid", "")
			for _, name := range []string{"missing", "wrong-type", "invalid"} {
				if _, err := getter.get(f, name); err == nil {
					t.Errorf("getter for %q returned no error", name)
				}
			}
		})
	}
}
