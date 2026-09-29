package pflag

import (
	"os"
	"strings"
	"testing"
)

func setUpCount(c *int) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.CountVarP(c, "verbose", "v", "a counter")
	return f
}

func TestCount(t *testing.T) {
	testCases := []struct {
		input    []string
		success  bool
		expected int
	}{
		{[]string{}, true, 0},
		{[]string{"-v"}, true, 1},
		{[]string{"-vvv"}, true, 3},
		{[]string{"-v", "-v", "-v"}, true, 3},
		{[]string{"-v", "--verbose", "-v"}, true, 3},
		{[]string{"-v=3", "-v"}, true, 4},
		{[]string{"--verbose=0"}, true, 0},
		{[]string{"-v=0"}, true, 0},
		{[]string{"-v=a"}, false, 0},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		var count int
		f := setUpCount(&count)

		tc := &testCases[i]

		err := f.Parse(tc.input)
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure, got success")
			continue
		} else if tc.success {
			c, err := f.GetCount("verbose")
			if err != nil {
				t.Errorf("Got error trying to fetch the counter flag")
			}
			if c != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, c)
			}
		}
	}
}

// Count flags take an optional --flag=N value (NoOptDefVal), so --foo 2
// treats "2" as a positional argument. Help used to print a required-looking
// "count" argument, which implied the space-separated form would work (#267).
func TestCountUsageShowsOptionalValue(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	var n int
	f.CountVar(&n, "foo", "increment")
	usage := f.FlagUsages()
	if strings.Contains(usage, "--foo count") {
		t.Fatalf("usage presents count as a required argument:\n%s", usage)
	}
	if !strings.Contains(usage, "--foo[=count]") {
		t.Fatalf("usage should show optional --foo[=count], got:\n%s", usage)
	}

	f2 := NewFlagSet("test", ContinueOnError)
	var n2 int
	f2.CountVar(&n2, "foo", "increment")
	if err := f2.Parse([]string{"--foo", "2"}); err != nil {
		t.Fatal(err)
	}
	if n2 != 1 {
		t.Fatalf("space-separated --foo 2 should increment once, got %d", n2)
	}
	if got := f2.Args(); len(got) != 1 || got[0] != "2" {
		t.Fatalf("space-separated 2 should remain a positional arg, got %v", got)
	}
}
