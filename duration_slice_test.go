// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code ds governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationSlice(t *testing.T) {
	t.Parallel()

	newFlag := func(dsp *[]time.Duration) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.DurationSliceVar(dsp, "ds", []time.Duration{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		ds := make([]time.Duration, 0)
		f := newFlag(&ds)
		require.NoError(t, f.Parse([]string{}))

		getDS, err := f.GetDurationSlice("ds")
		require.NoErrorf(t, err,
			"got an error from GetDurationSlice(): %v", err,
		)
		require.Empty(t, getDS)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"1ns", "2ms", "3m", "4h"}
		ds := make([]time.Duration, 0, len(vals))
		f := newFlag(&ds)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--ds=%s", strings.Join(vals, ",")),
		}))

		for i, v := range ds {
			d, err := time.ParseDuration(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected ds[%d] to be %s but got: %d", i, vals[i], v,
			)
		}

		getDS, erd := f.GetDurationSlice("ds")
		require.NoError(t, erd)

		for i, v := range getDS {
			d, err := time.ParseDuration(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected ds[%d] to be %s but got: %d from GetDurationSlice", i, vals[i], v,
			)
		}
	})

	newFlagWithDefault := func(dsp *[]time.Duration) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.DurationSliceVar(dsp, "ds", []time.Duration{0, 1}, "Command separated list!")
		return f
	}

	t.Run("with default (1)", func(t *testing.T) {
		vals := []string{"0s", "1ns"}
		ds := make([]time.Duration, 0, len(vals))
		f := newFlagWithDefault(&ds)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range ds {
			d, err := time.ParseDuration(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected ds[%d] to be %d but got: %d", i, d, v,
			)
		}

		getDS, erd := f.GetDurationSlice("ds")
		require.NoErrorf(t, erd,
			"got an error from GetDurationSlice(): %v", erd,
		)

		for i, v := range getDS {
			d, err := time.ParseDuration(vals[i])
			require.NoErrorf(t, err,
				"got an error from GetDurationSlice(): %v", err,
			)
			require.Equalf(t, v, d,
				"expected ds[%d] to be %d from GetDurationSlice but got: %d", i, d, v,
			)
		}
	})

	t.Run("with default (2)", func(t *testing.T) {
		vals := []string{"1ns", "2ns"}
		ds := make([]time.Duration, 0, len(vals))
		f := newFlagWithDefault(&ds)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--ds=%s", strings.Join(vals, ",")),
		}))

		for i, v := range ds {
			d, err := time.ParseDuration(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected ds[%d] to be %d but got: %d", i, d, v,
			)
		}

		getDS, erd := f.GetDurationSlice("ds")
		require.NoErrorf(t, erd,
			"got an error from GetDurationSlice(): %v", erd,
		)

		for i, v := range getDS {
			d, err := time.ParseDuration(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected ds[%d] to be %d from GetDurationSlice but got: %d", i, d, v,
			)
		}
	})

	t.Run("as SliceValue", func(t *testing.T) {
		in := []string{"1ns", "2ns"}
		ds := make([]time.Duration, 0, len(in))
		f := newFlag(&ds)

		argfmt := "--ds=%s"
		arg1 := fmt.Sprintf(argfmt, in[0])
		arg2 := fmt.Sprintf(argfmt, in[1])
		require.NoError(t, f.Parse([]string{arg1, arg2}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"3ns"}))
			}
		})

		require.Equalf(t, []time.Duration{time.Duration(3)}, ds,
			"expected ss to be overwritten with '3ns', but got: %v", ds,
		)
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--ds=%s"
		in := []string{"1ns,2ns", "3ns"}
		ds := make([]time.Duration, 0, len(in))
		f := newFlag(&ds)
		expected := []time.Duration{1, 2, 3}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, ds)
	})
}
