package pflag

import (
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestBoolFunc(t *testing.T) {
	var count int
	fn := func(_ string) error {
		count++
		return nil
	}

	fset := NewFlagSet("test", ContinueOnError)
	fset.BoolFunc("func", "Callback function", fn)

	err := fset.Parse([]string{"--func", "--func=1", "--func=false"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if count != 3 {
		t.Fatalf("expected 3 calls to the callback, got %d calls", count)
	}
}

func TestBoolFuncP(t *testing.T) {
	var count int
	fn := func(_ string) error {
		count++
		return nil
	}

	fset := NewFlagSet("test", ContinueOnError)
	fset.BoolFuncP("bfunc", "b", "Callback function", fn)

	err := fset.Parse([]string{"--bfunc", "--bfunc=0", "--bfunc=false", "-b", "-b=0"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	if count != 5 {
		t.Fatalf("expected 5 calls to the callback, got %d calls", count)
	}
}

func TestBoolFuncCompat(t *testing.T) {
	// compare behavior with the stdlib 'flag' package
	type BoolFuncFlagSet interface {
		BoolFunc(name string, usage string, fn func(string) error)
		Parse([]string) error
	}

	unitTestErr := errors.New("unit test error")
	runCase := func(f BoolFuncFlagSet, name string, args []string) (values []string, err error) {
		fn := func(s string) error {
			values = append(values, s)
			if s == "err" {
				return unitTestErr
			}
			return nil
		}
		f.BoolFunc(name, "Callback function", fn)

		err = f.Parse(args)
		return values, err
	}

	t.Run("regular parsing", func(t *testing.T) {
		flagName := "bflag"
		args := []string{"--bflag", "--bflag=false", "--bflag=1", "--bflag=bar", "--bflag="}

		// It turns out that, even though the function is called "BoolFunc",
		// the standard flag package does not try to parse the value assigned to
		// that cli flag as a boolean. The string provided on the command line is
		// passed as is to the callback.
		//   e.g: with "--bflag=not_a_bool" on the command line, the FlagSet does not
		// generate an error stating "invalid boolean value", and `fn` will be called
		// with "not_a_bool" as an argument.

		stdFSet := flag.NewFlagSet("std test", flag.ContinueOnError)
		stdValues, err := runCase(stdFSet, flagName, args)
		if err != nil {
			t.Fatalf("std flag: expected no error, got %v", err)
		}
		expected := []string{"true", "false", "1", "bar", ""}
		if !cmpLists(expected, stdValues) {
			t.Fatalf("std flag: expected %v, got %v", expected, stdValues)
		}

		fset := NewFlagSet("pflag test", ContinueOnError)
		pflagValues, err := runCase(fset, flagName, args)
		if err != nil {
			t.Fatalf("pflag: expected no error, got %v", err)
		}
		if !cmpLists(stdValues, pflagValues) {
			t.Fatalf("pflag: expected %v, got %v", stdValues, pflagValues)
		}
	})

	t.Run("error triggered by callback", func(t *testing.T) {
		flagName := "bflag"
		args := []string{"--bflag", "--bflag=err", "--bflag=after"}

		// test behavior of standard flag.Fset with an error triggered by the callback:
		// (note: as can be seen in 'runCase()', if the callback sees "err" as a value
		//  for the bool flag, it will return an error)
		stdFSet := flag.NewFlagSet("std test", flag.ContinueOnError)
		stdFSet.SetOutput(io.Discard) // suppress output

		// run test case with standard flag.Fset
		stdValues, err := runCase(stdFSet, flagName, args)

		// double check the standard behavior:
		// - .Parse() should return an error, which contains the error message
		if err == nil {
			t.Fatalf("std flag: expected an error triggered by callback, got no error instead")
		}
		if !strings.HasSuffix(err.Error(), unitTestErr.Error()) {
			t.Fatalf("std flag: expected unittest error, got unexpected error value: %T %v", err, err)
		}
		// - the function should have been called twice, with the first two values,
		//   the final "=after" should not be recorded
		expected := []string{"true", "err"}
		if !cmpLists(expected, stdValues) {
			t.Fatalf("std flag: expected %v, got %v", expected, stdValues)
		}

		// now run the test case on a pflag FlagSet:
		fset := NewFlagSet("pflag test", ContinueOnError)
		pflagValues, err := runCase(fset, flagName, args)

		// check that there is a similar error (note: pflag will _wrap_ the error, while the stdlib
		// currently keeps the original message but creates a flat errors.Error)
		if !errors.Is(err, unitTestErr) {
			t.Fatalf("pflag: got unexpected error value: %T %v", err, err)
		}
		// the callback should be called the same number of times, with the same values:
		if !cmpLists(stdValues, pflagValues) {
			t.Fatalf("pflag: expected %v, got %v", stdValues, pflagValues)
		}
	})
}

func TestBoolFuncUsage(t *testing.T) {
	t.Run("regular func flag", func(t *testing.T) {
		// regular boolfunc flag:
		// expect to see '--flag1' followed by the usageMessage, and no mention of a default value
		fset := NewFlagSet("unittest", ContinueOnError)
		fset.BoolFunc("flag1", "usage message", func(s string) error { return nil })
		usage := fset.FlagUsagesWrapped(80)

		usage = strings.TrimSpace(usage)
		expected := "--flag1   usage message"
		if usage != expected {
			t.Fatalf("unexpected generated usage message\n  expected: %s\n       got: %s", expected, usage)
		}
	})

	t.Run("func flag with placeholder name", func(t *testing.T) {
		// func flag, with a placeholder name:
		// if usageMesage contains a placeholder, expect '--flag2 {placeholder}'; still expect no mention of a default value
		fset := NewFlagSet("unittest", ContinueOnError)
		fset.BoolFunc("flag2", "usage message with `name` placeholder", func(s string) error { return nil })
		usage := fset.FlagUsagesWrapped(80)

		usage = strings.TrimSpace(usage)
		expected := "--flag2 name   usage message with name placeholder"
		if usage != expected {
			t.Fatalf("unexpected generated usage message\n  expected: %s\n       got: %s", expected, usage)
		}
	})
}
