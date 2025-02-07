package pflag

import (
	"os"
	"reflect"
	"testing"
)

func setUpCount(c *int) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.CountVarP(c, "verbose", "v", "a counter")
	return f
}

func TestCountValueImplementsGetter(t *testing.T) {
	var v Value = new(countValue)

	if _, ok := v.(Getter); !ok {
		t.Fatalf("%T should implement the Getter interface", v)
	}
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

			c2, err := f.Get("verbose")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}

			if !reflect.DeepEqual(c, c2) {
				t.Fatalf("expected %v with type %T but got %v with type %T", c, c, c2, c2)
			}
		}
	}
}
