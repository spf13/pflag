package pflag

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCount(t *testing.T) {
	newFlag := func(c *int) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.CountVarP(c, "verbose", "v", "a counter")
		return f
	}

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
		f := newFlag(&count)

		tc := &testCases[i]

		err := f.Parse(tc.input)
		if !tc.success {
			require.Errorf(t, err,
				"expected failure with %q, got success", tc.input,
			)

			continue
		}

		require.NoError(t, err)

		c, err := f.GetCount("verbose")
		require.NoError(t, err)
		require.Equal(t, tc.expected, c)
	}
}
