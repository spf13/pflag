package pflag

import (
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func cmpLists(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFunc(t *testing.T) {
	var values []string
	fn := func(s string) error {
		values = append(values, s)
		return nil
	}

	fset := NewFlagSet("test", ContinueOnError)
	fset.Func("fnflag", "Callback function", fn)

	err := fset.Parse([]string{"--fnflag=aa", "--fnflag", "bb"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	expected := []string{"aa", "bb"}
	if !cmpLists(expected, values) {
		t.Fatalf("expected %v, got %v", expected, values)
	}
}

func TestFuncP(t *testing.T) {
	var values []string
	fn := func(s string) error {
		values = append(values, s)
		return nil
	}

	fset := NewFlagSet("test", ContinueOnError)
	fset.FuncP("fnflag", "f", "Callback function", fn)

	err := fset.Parse([]string{"--fnflag=a", "--fnflag", "b", "-fc", "-f=d", "-f", "e"})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	expected := []string{"a", "b", "c", "d", "e"}
	if !cmpLists(expected, values) {
		t.Fatalf("expected %v, got %v", expected, values)
	}
}

func TestFuncCompat(t *testing.T) {
	// compare behavior with the stdlib 'flag' package
	type FuncFlagSet interface {
		Func(name string, usage string, fn func(string) error)
		Parse([]string) error
	}

	unitTestErr := errors.New("unit test error")
	runCase := func(f FuncFlagSet, name string, args []string) (values []string, err error) {
		fn := func(s string) error {
			values = append(values, s)
			if s == "err" {
				return unitTestErr
			}
			return nil
		}
		f.Func(name, "Callback function", fn)

		err = f.Parse(args)
		return values, err
	}

	t.Run("regular parsing", func(t *testing.T) {
		flagName := "fnflag"
		args := []string{"--fnflag=xx", "--fnflag", "yy", "--fnflag=zz"}

		stdFSet := flag.NewFlagSet("std test", flag.ContinueOnError)
		stdValues, err := runCase(stdFSet, flagName, args)
		if err != nil {
			t.Fatalf("std flag: expected no error, got %v", err)
		}
		expected := []string{"xx", "yy", "zz"}
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
		flagName := "fnflag"
		args := []string{"--fnflag", "before", "--fnflag", "err", "--fnflag", "after"}

		// test behavior of standard flag.Fset with an error triggered by the callback:
		// (note: as can be seen in 'runCase()', if the callback sees "err" as a value
		//  for the flag, it will return an error)
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
		expected := []string{"before", "err"}
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

func TestFuncUsage(t *testing.T) {
	t.Run("regular func flag", func(t *testing.T) {
		// regular func flag:
		// expect to see '--flag1 value' followed by the usageMessage, and no mention of a default value
		fset := NewFlagSet("unittest", ContinueOnError)
		fset.Func("flag1", "usage message", func(s string) error { return nil })
		usage := fset.FlagUsagesWrapped(80)

		usage = strings.TrimSpace(usage)
		expected := "--flag1 value   usage message"
		if usage != expected {
			t.Fatalf("unexpected generated usage message\n  expected: %s\n       got: %s", expected, usage)
		}
	})

	t.Run("func flag with placeholder name", func(t *testing.T) {
		// func flag, with a placeholder name:
		// if usageMesage contains a placeholder, expect that name; still expect no mention of a default value
		fset := NewFlagSet("unittest", ContinueOnError)
		fset.Func("flag2", "usage message with `name` placeholder", func(s string) error { return nil })
		usage := fset.FlagUsagesWrapped(80)

		usage = strings.TrimSpace(usage)
		expected := "--flag2 name   usage message with name placeholder"
		if usage != expected {
			t.Fatalf("unexpected generated usage message\n  expected: %s\n       got: %s", expected, usage)
		}
	})
}
