// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringArray(t *testing.T) {
	t.Parallel()

	newFlag := func(sap *[]string) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringArrayVar(sap, "sa", []string{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		sa := make([]string, 0)
		f := newFlag(&sa)
		require.NoError(t, f.Parse([]string{}))

		getSA, err := f.GetStringArray("sa")
		require.NoErrorf(t, err,
			"got an error from GetStringArray(): %v", err,
		)
		require.Empty(t, getSA)
	})

	t.Run("with empty value", func(t *testing.T) {
		sa := make([]string, 0)
		f := newFlag(&sa)
		require.NoError(t, f.Parse([]string{"--sa="}))

		getSA, err := f.GetStringArray("sa")
		require.NoErrorf(t, err,
			"got an error from GetStringArray(): %v", err,
		)
		require.Empty(t, getSA)
	})

	newFlagWithDefault := func(sap *[]string) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringArrayVar(sap, "sa", []string{"default", "values"}, "Command separated list!")
		return f
	}

	t.Run("with default (1)", func(t *testing.T) {
		vals := []string{"default", "values"}
		sa := make([]string, 0, len(vals))
		f := newFlagWithDefault(&sa)

		require.NoError(t, f.Parse([]string{}))
		require.Equal(t, vals, sa)

		getSA, err := f.GetStringArray("sa")
		require.NoError(t, err)
		require.Equal(t, vals, getSA)
	})

	t.Run("with default (2)", func(t *testing.T) {
		val := "one"
		sa := make([]string, 0, len(val))
		f := newFlagWithDefault(&sa)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--sa=%s", val),
		}))
		require.Equal(t, []string{val}, sa)

		getSA, err := f.GetStringArray("sa")
		require.NoErrorf(t, err,
			"got an error from GetStringArray(): %v", err,
		)
		require.Equal(t, []string{val}, getSA)
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--sa=%s"
		in := []string{"one", "two"}
		sa := make([]string, 0, len(in))
		f := newFlag(&sa)
		expected := []string{"one", "two"}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))
		require.Equal(t, expected, sa)

		values, err := f.GetStringArray("sa")
		require.NoError(t, err)
		require.Equal(t, expected, values)
	})

	t.Run("with special char", func(t *testing.T) {
		const argfmt = "--sa=%s"
		in := []string{"one,two", `"three"`, `"four,five",six`, "seven eight"}
		sa := make([]string, 0, len(in))
		f := newFlag(&sa)
		expected := []string{"one,two", `"three"`, `"four,five",six`, "seven eight"}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
			fmt.Sprintf(argfmt, in[2]),
			fmt.Sprintf(argfmt, in[3]),
		}))
		require.Equal(t, expected, sa)

		values, err := f.GetStringArray("sa")
		require.NoError(t, err)
		require.Equal(t, expected, values)
	})

	t.Run("with square bracket", func(t *testing.T) {
		const argfmt = "--sa=%s"
		in := []string{"][]-[", "[a-z]", "[a-z]+"}
		sa := make([]string, 0, len(in))
		f := newFlag(&sa)
		expected := []string{"][]-[", "[a-z]", "[a-z]+"}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
			fmt.Sprintf(argfmt, in[2]),
		}))
		require.Equal(t, expected, sa)

		values, err := f.GetStringArray("sa")
		require.NoError(t, err)
		require.Equal(t, expected, values)
	})

	t.Run("with slice as value", func(t *testing.T) {
		const argfmt = "--sa=%s"
		in := []string{"1ns", "2ns"}
		sa := make([]string, 0, len(in))
		f := newFlag(&sa)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"3ns"}))
			}
		})
		require.Equalf(t, []string{"3ns"}, sa,
			"expected ss to be overwritten with '3ns', but got: %v", sa,
		)
	})
}

func TestStringArrayConv(t *testing.T) {
	t.Run("with empty string", func(t *testing.T) {
		_, err := stringArrayConv("")
		require.NoError(t, err)
	})
}
