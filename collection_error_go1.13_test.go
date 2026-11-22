//go:build go1.13
// +build go1.13

package pflag

import (
	"errors"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestCollectionArgumentErrorCause(t *testing.T) {
	for _, tc := range collectionErrorCases {
		if tc.name == "durationSlice" {
			continue // time.ParseDuration does not return a *strconv.NumError.
		}
		t.Run(tc.name, func(t *testing.T) {
			inputs := []struct {
				value string
				cause error
			}{{tc.invalid, strconv.ErrSyntax}}
			if tc.overflow != "" {
				inputs = append(inputs, struct {
					value string
					cause error
				}{tc.overflow, strconv.ErrRange})
			}
			for _, input := range inputs {
				f := NewFlagSet("test", ContinueOnError)
				f.SetOutput(ioutil.Discard)
				tc.register(f)
				raw := tc.prefix + input.value
				err := f.Parse([]string{"--value=" + raw})
				var invalid *InvalidValueError
				if !errors.As(err, &invalid) || invalid.GetFlag() != f.Lookup("value") || invalid.GetValue() != raw {
					t.Errorf("Parse(%q) error = %v, want InvalidValueError for the flag and input", raw, err)
				}
				var numErr *strconv.NumError
				if !errors.As(err, &numErr) || numErr.Num != input.value || !errors.Is(err, input.cause) {
					t.Errorf("Parse(%q) error = %v, want NumError for %q wrapping %v", raw, err, input.value, input.cause)
				}
			}
		})
	}
}
