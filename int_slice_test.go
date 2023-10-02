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

func TestIntSlice(t *testing.T) {
	t.Parallel()

	newFlag := func(isp *[]int) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.IntSliceVar(isp, "is", []int{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		is := make([]int, 0)
		f := newFlag(&is)
		require.NoError(t, f.Parse([]string{}))

		getIS, err := f.GetIntSlice("is")
		require.NoErrorf(t, err,
			"got an error from GetIntSlice(): %v", err,
		)
		require.Empty(t, getIS)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"1", "2", "4", "3"}
		is := make([]int, 0, len(vals))
		f := newFlag(&is)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--is=%s", strings.Join(vals, ",")),
		}))

		for i, v := range is {
			d, err := strconv.Atoi(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected is[%d] to be %d but got: %d", i, v, d,
			)
		}

		getIS, eri := f.GetIntSlice("is")
		require.NoError(t, eri)

		for i, v := range getIS {
			d, err := strconv.Atoi(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected is[%d] to be %d but got: %d from GetIntSlice", i, v, d,
			)
		}
	})

	newFlagWithDefault := func(isp *[]int) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.IntSliceVar(isp, "is", []int{0, 1}, "Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"0", "1"}
		is := make([]int, 0, len(vals))
		f := newFlagWithDefault(&is)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range is {
			d, err := strconv.Atoi(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected is[%d] to be %d but got: %d", i, v, d,
			)
		}

		getIS, eri := f.GetIntSlice("is")
		require.NoErrorf(t, eri,
			"got an error from GetIntSlice(): %v", eri,
		)

		for i, v := range getIS {
			d, err := strconv.Atoi(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected is[%d] to be %d from GetIntSlice but got: %d", i, v, d,
			)
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"1", "2"}
		is := make([]int, 0, len(vals))
		f := newFlagWithDefault(&is)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--is=%s", strings.Join(vals, ",")),
		}))

		for i, v := range is {
			d, err := strconv.Atoi(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected is[%d] to be %d but got: %d", i, v, d,
			)
		}

		getIS, eri := f.GetIntSlice("is")
		require.NoErrorf(t, eri,
			"got an error from GetIntSlice(): %v", eri,
		)

		for i, v := range getIS {
			d, err := strconv.Atoi(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, d,
				"expected is[%d] to be %d from GetIntSlice but got: %d", i, v, d,
			)
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--is=%s"
		in := []string{"1,2", "3"}
		is := make([]int, 0, len(in))
		f := newFlag(&is)
		expected := []int{1, 2, 3}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, is)
	})
}
