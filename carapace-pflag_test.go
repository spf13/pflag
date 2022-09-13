package pflag

import (
	"reflect"
	"testing"
)

func TestLongShorthand(t *testing.T) {
	f := NewFlagSet("longShorthand", ContinueOnError)
	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "ab", false, "bool2 value")
	f.BoolP("boolc", "c", false, "bool value")
	f.StringP("stringa", "s", "0", "string value")
	f.StringP("stringx", "sx", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"
	args := []string{
		"-ab",
		"-sx=something",
	}
	want := []string{
		"boolb", "true",
		"stringx", "something",
	}
	got := []string{}
	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}
	if err := f.ParseAll(args, store); err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if !f.Parsed() {
		t.Errorf("f.Parse() = false after Parse")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("f.TestLongShorthand() fail to restore the args")
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}
}
