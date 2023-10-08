// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringSlice(t *testing.T) {
	t.Parallel()

	newFlag := func(ssp *[]string) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringSliceVar(ssp, "ss", []string{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		ss := make([]string, 0)
		f := newFlag(&ss)

		require.NoError(t, f.Parse([]string{}))

		getSS, err := f.GetStringSlice("ss")
		require.NoErrorf(t, err,
			"got an error from GetStringSlice(): %v", err,
		)
		require.Empty(t, getSS)
	})

	t.Run("with empty values", func(t *testing.T) {
		ss := make([]string, 0)
		f := newFlag(&ss)
		require.NoError(t, f.Parse([]string{"--ss="}))

		getSS, err := f.GetStringSlice("ss")
		require.NoErrorf(t, err,
			"got an error from GetStringSlice(): %v", err,
		)
		require.Empty(t, getSS)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"one", "two", "4", "3"}
		ss := make([]string, 0, len(vals))
		f := newFlag(&ss)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--ss=%s", strings.Join(vals, ",")),
		}))

		require.Equal(t, vals, ss)

		getSS, err := f.GetStringSlice("ss")
		require.NoErrorf(t, err,
			"got an error from GetStringSlice(): %v", err,
		)
		require.Equal(t, vals, getSS)
	})

	newFlagWithDefault := func(ssp *[]string) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringSliceVar(ssp, "ss", []string{"default", "values"}, "Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"default", "values"}
		ss := make([]string, 0, len(vals))
		f := newFlagWithDefault(&ss)

		require.NoError(t, f.Parse([]string{}))
		require.Equal(t, vals, ss)

		getSS, err := f.GetStringSlice("ss")
		require.NoErrorf(t, err,
			"got an error from GetStringSlice(): %v", err,
		)
		require.Equal(t, vals, getSS)
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"one", "two", "4", "3"}
		ss := make([]string, 0, len(vals))
		f := newFlagWithDefault(&ss)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--ss=%s", strings.Join(vals, ",")),
		}))
		require.Equal(t, vals, ss)

		getSS, err := f.GetStringSlice("ss")
		require.NoErrorf(t, err,
			"got an error from GetStringSlice(): %v", err,
		)
		require.Equal(t, vals, getSS)
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--ss=%s"
		in := []string{"one,two", "three"}
		ss := make([]string, 0, len(in))
		f := newFlag(&ss)
		expected := []string{"one", "two", "three"}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))
		require.Equal(t, expected, ss)

		values, err := f.GetStringSlice("ss")
		require.NoError(t, err)
		require.Equal(t, expected, values)
	})

	t.Run("with comma", func(t *testing.T) {
		const argfmt = "--ss=%s"
		in := []string{`"one,two"`, `"three"`, `"four,five",six`}
		ss := make([]string, 0, len(in))
		f := newFlag(&ss)
		expected := []string{"one,two", "three", "four,five", "six"}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
			fmt.Sprintf(argfmt, in[2]),
		}))

		require.Equal(t, expected, ss)

		values, err := f.GetStringSlice("ss")
		require.NoError(t, err)
		require.Equal(t, expected, values)
	})

	t.Run("with square bracket", func(t *testing.T) {
		const argfmt = "--ss=%s"
		in := []string{`"[a-z]"`, `"[a-z]+"`}
		ss := make([]string, 0, len(in))
		f := newFlag(&ss)
		expected := []string{"[a-z]", "[a-z]+"}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, ss)

		values, err := f.GetStringSlice("ss")
		require.NoError(t, err)
		require.Equal(t, expected, values)
	})

	t.Run("as slice value", func(t *testing.T) {
		const argfmt = "--ss=%s"
		in := []string{"one", "two"}
		ss := make([]string, 0, len(in))
		f := newFlag(&ss)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"three"}))
			}
		})
		require.Equalf(t, []string{"three"}, ss,
			"expected ss to be overwritten with 'three', but got: %s", ss,
		)
	})
}
