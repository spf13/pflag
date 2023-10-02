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

func TestInt64Slice(t *testing.T) {
	t.Parallel()

	newFlag := func(isp *[]int64) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.Int64SliceVar(isp, "is", []int64{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		is := make([]int64, 0)
		f := newFlag(&is)
		require.NoError(t, f.Parse([]string{}))

		getI64S, err := f.GetInt64Slice("is")
		require.NoErrorf(t, err,
			"got an error from GetInt64Slice(): %v", err,
		)
		require.Empty(t, getI64S)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"1", "2", "4", "3"}
		is := make([]int64, 0, len(vals))
		f := newFlag(&is)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--is=%s", strings.Join(vals, ",")),
		}))

		for i, v := range is {
			d, err := strconv.ParseInt(vals[i], 0, 64)
			require.NoError(t, err)

			require.Equalf(t, v, d,
				"expected is[%d] to be %s but got: %d", i, vals[i], v,
			)
		}

		getI64S, eri := f.GetInt64Slice("is")
		require.NoError(t, eri)

		for i, v := range getI64S {
			d, err := strconv.ParseInt(vals[i], 0, 64)
			require.NoError(t, err)

			require.Equalf(t, v, d,
				"expected is[%d] to be %s but got: %d from GetInt64Slice", i, vals[i], v,
			)
		}
	})

	newFlagWithDefault := func(isp *[]int64) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.Int64SliceVar(isp, "is", []int64{0, 1}, "Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"0", "1"}
		is := make([]int64, 0, len(vals))
		f := newFlagWithDefault(&is)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range is {
			d, err := strconv.ParseInt(vals[i], 0, 64)
			require.NoError(t, err)

			require.Equalf(t, v, d,
				"expected is[%d] to be %d but got: %d", i, v, d,
			)
		}

		getI64S, eri := f.GetInt64Slice("is")
		require.NoErrorf(t, eri,
			"got an error from GetInt64Slice(): %v", eri,
		)

		for i, v := range getI64S {
			d, err := strconv.ParseInt(vals[i], 0, 64)
			require.NoErrorf(t, err,
				"got an error from GetInt64Slice(): %v", err,
			)
			require.Equalf(t, v, d,
				"expected is[%d] to be %d from GetInt64Slice but got: %d", i, v, d,
			)
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"1", "2"}
		is := make([]int64, 0, len(vals))
		f := newFlagWithDefault(&is)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--is=%s", strings.Join(vals, ",")),
		}))

		for i, v := range is {
			d, err := strconv.ParseInt(vals[i], 0, 64)
			require.NoError(t, err)

			require.Equalf(t, v, d,
				"expected is[%d] to be %d but got: %d", i, v, d,
			)
		}

		getI64S, eri := f.GetInt64Slice("is")
		require.NoErrorf(t, eri,
			"got an error from GetInt64Slice(): %v", eri,
		)

		for i, v := range getI64S {
			d, err := strconv.ParseInt(vals[i], 0, 64)
			require.NoError(t, err)

			require.Equalf(t, v, d,
				"expected is[%d] to be %d from GetInt64Slice but got: %d", i, d, v,
			)
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--is=%s"
		in := []string{"1,2", "3"}
		is := make([]int64, 0, len(in))
		f := newFlag(&is)
		expected := []int64{1, 2, 3}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, is)
	})

	t.Run("as slice value", func(t *testing.T) {
		const argfmt = "--is=%s"
		in := []string{"1", "2"}
		i64s := make([]int64, 0, len(in))
		f := newFlag(&i64s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"3"}))
			}
		})

		require.Equalf(t, []int64{3}, i64s,
			"expected ss to be overwritten with '3.1', but got: %v", i64s,
		)
	})
}
