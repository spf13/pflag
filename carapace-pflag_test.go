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

func TestNonPosix(t *testing.T) {
	f := NewFlagSet("nonPosix", ContinueOnError)
	f.StringN("stringa", "sa", "0", "string value")
	f.StringN("stringx", "sx", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"
	args := []string{
		"-sa", "somearg",
		"-stringx=something",
	}
	want := []string{
		"stringa", "somearg",
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

func TestOptargDelimiter(t *testing.T) {
	f := NewFlagSet("optargdelimiter", ContinueOnError)
	f.StringN("stringa", "a", "0", "string value")
	f.StringN("stringx", "x", "0", "string value")
	f.Lookup("stringa").NoOptDefVal = "1"
	f.Lookup("stringa").OptargDelimiter = '/'
	f.Lookup("stringx").NoOptDefVal = "2"
	f.Lookup("stringx").OptargDelimiter = ':'

	args := []string{
		"-stringa/somearg",
		"-stringx:something",
	}
	want := []string{
		"stringa", "somearg",
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

func TestNargs(t *testing.T) {
	f := NewFlagSet("nargs", ContinueOnError)
	f.StringSlice("stringa", []string{}, "string value")
	f.StringSlice("stringx", []string{}, "string value")
	f.Lookup("stringa").Nargs = 2
	f.Lookup("stringx").Nargs = -1

	args := []string{
		"--stringa", "one", "two", "three",
		"--stringx", "four", "five",
	}
	want := []string{
		"stringa", "one,two",
		"stringx", "four,five",
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

	// ensure slice is correctly set
	f.Parse(args)
	got, err := f.GetStringSlice("stringa")
	if err != nil {
		t.Error(err.Error())
	}
	want = []string{"one", "two"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Got:  %v", got)
		t.Errorf("Want: %v", want)
	}

}
