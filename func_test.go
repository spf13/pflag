package pflag

import (
	"errors"
	"testing"
)

func TestFunc(t *testing.T) {
	testCases := []struct {
		input    []string
		success  bool
		expected map[string]int
	}{
		{[]string{}, true, map[string]int{}},
		{[]string{"-f=foo"}, true, map[string]int{"foo": 1}},
		{[]string{"-f", "quux", "--fooq", "quux"}, true, map[string]int{"quux": 2}},
		{[]string{"--fooq=quux", "-bbaz", "--barz", "baz", "-b=bar", "-f", "foo"}, true,
			map[string]int{"foo": 1, "bar": 1, "baz": 2, "quux": 1}},
		{[]string{"-f=foo", "-f=quux", "-fquux", "--fooq=quux", "--fooq", "foo"}, true,
			map[string]int{"foo": 2, "quux": 3}},
		{[]string{"--barz=quux"}, false, map[string]int{}},
		{[]string{"--barz=", "bar"}, false, map[string]int{}},
	}
	var m map[string]int
	setUpFunc := func() *FlagSet {
		m = make(map[string]int, 4)
		f := NewFlagSet("test", ContinueOnError)
		fooquux := func(s string) error {
			switch s {
			case "foo", "quux":
				m[s]++
				return nil
			}
			return errors.New("unrecognized arg")
		}
		barbaz := func(s string) error {
			switch s {
			case "bar", "baz":
				m[s]++
				return nil
			}
			return errors.New("unrecognized arg")
		}
		f.FuncP("fooq", "f", "a counter for foo and quux args", fooquux)
		f.FuncP("barz", "b", "a counter for bar and baz args", barbaz)
		return f
	}
	for i := range testCases {
		f := setUpFunc()

		tc := &testCases[i]

		err := f.Parse(tc.input)
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure, got success")
			continue
		} else if tc.success {
			if got, expected := len(m), len(tc.expected); got != expected {
				t.Errorf("%d: expected map length %d, got %d", i, expected, got)
			}
			for key := range m {
				got, expected := m[key], tc.expected[key]
				if got != expected {
					t.Errorf("%d: expected [%q]:%d, got [%q]:%d", i, key, expected, key, got)
				}
			}
		}
	}
}
