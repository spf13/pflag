// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat64Slice(t *testing.T) {
	t.Parallel()

	newFlag := func(f64sp *[]float64) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.Float64SliceVar(f64sp, "f64s", []float64{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		f64s := make([]float64, 0)
		f := newFlag(&f64s)

		require.NoError(t, f.Parse([]string{}))

		getF64S, err := f.GetFloat64Slice("f64s")
		require.NoErrorf(t, err,
			"got an error from GetFloat64Slice(): %v", err,
		)
		require.Empty(t, getF64S)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"1.0", "2.0", "4.0", "3.0"}
		f64s := make([]float64, 0, len(vals))
		f := newFlag(&f64s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--f64s=%s", strings.Join(vals, ",")),
		}))

		for i, v := range f64s {
			d, err := strconv.ParseFloat(vals[i], 64)
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected f64s[%d] to be %s but got: %f", i, vals[i], v,
			)
		}

		getF64S, err := f.GetFloat64Slice("f64s")
		require.NoError(t, err)

		for i, v := range getF64S {
			d, err := strconv.ParseFloat(vals[i], 64)
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected f64s[%d] to be %s but got: %f from GetFloat64Slice", i, vals[i], v,
			)
		}
	})

	newFlagWithDefault := func(f64sp *[]float64) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.Float64SliceVar(f64sp, "f64s", []float64{0.0, 1.0}, "Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"0.0", "1.0"}
		f64s := make([]float64, 0, len(vals))
		f := newFlagWithDefault(&f64s)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range f64s {
			d, err := strconv.ParseFloat(vals[i], 64)
			require.NoError(t, err)

			require.Equalf(t, v, d,
				"expected f64s[%d] to be %f but got: %f", i, d, v,
			)
		}

		getF64S, erf := f.GetFloat64Slice("f64s")
		require.NoErrorf(t, erf,
			"got an error from GetFloat64Slice(): %v", erf,
		)

		for i, v := range getF64S {
			d, err := strconv.ParseFloat(vals[i], 64)
			require.NoErrorf(t, err,
				"got an error from GetFloat64Slice(): %v", err,
			)
			require.Equalf(t, v, d,
				"expected f64s[%d] to be %f from GetFloat64Slice but got: %f", i, d, v,
			)
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"1.0", "2.0"}
		f64s := make([]float64, 0, len(vals))
		f := newFlagWithDefault(&f64s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--f64s=%s", strings.Join(vals, ",")),
		}))

		for i, v := range f64s {
			d, err := strconv.ParseFloat(vals[i], 64)
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected f64s[%d] to be %f but got: %f", i, d, v,
			)
		}

		getF64S, erf := f.GetFloat64Slice("f64s")
		require.NoErrorf(t, erf,
			"got an error from GetFloat64Slice(): %v", erf,
		)

		for i, v := range getF64S {
			d, err := strconv.ParseFloat(vals[i], 64)
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected f64s[%d] to be %f from GetFloat64Slice but got: %f", i, d, v,
			)
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--f64s=%s"
		in := []string{"1.0,2.0", "3.0"}
		f64s := make([]float64, 0, len(in))
		f := newFlag(&f64s)
		expected := []float64{1.0, 2.0, 3.0}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, f64s)
	})

	t.Run("as slice value", func(t *testing.T) {
		const argfmt = "--f64s=%s"
		in := []string{"1.0", "2.0"}
		f64s := make([]float64, 0, len(in))
		f := newFlag(&f64s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"3.1"}))
			}
		})

		require.Equalf(t, []float64{3.1}, f64s,
			"expected ss to be overwritten with '3.1', but got: %v", f64s,
		)
	})
}
