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

func TestFloat32Slice(t *testing.T) {
	t.Parallel()

	newFlag := func(f32sp *[]float32) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.Float32SliceVar(f32sp, "f32s", []float32{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		f32s := make([]float32, 0)
		f := newFlag(&f32s)

		require.NoError(t, f.Parse([]string{}))

		getF32S, err := f.GetFloat32Slice("f32s")
		require.NoErrorf(t, err,
			"got an error from GetFloat32Slice(): %v", err,
		)
		require.Empty(t, getF32S)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"1.0", "2.0", "4.0", "3.0"}
		f32s := make([]float32, 0, len(vals))
		f := newFlag(&f32s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--f32s=%s", strings.Join(vals, ",")),
		}))

		for i, v := range f32s {
			d64, err := strconv.ParseFloat(vals[i], 32)
			require.NoError(t, err)

			d := float32(d64)
			require.Equalf(t, v, d,
				"expected f32s[%d] to be %s but got: %f", i, vals[i], v,
			)
		}

		getF32S, erf := f.GetFloat32Slice("f32s")
		require.NoError(t, erf)

		for i, v := range getF32S {
			d64, err := strconv.ParseFloat(vals[i], 32)
			require.NoError(t, err)

			d := float32(d64)
			require.Equalf(t, v, d,
				"expected f32s[%d] to be %s but got: %f from GetFloat32Slice", i, vals[i], v,
			)
		}
	})

	newFlagWithDefault := func(f32sp *[]float32) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.Float32SliceVar(f32sp, "f32s", []float32{0.0, 1.0}, "Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"0.0", "1.0"}
		f32s := make([]float32, 0, len(vals))
		f := newFlagWithDefault(&f32s)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range f32s {
			d64, err := strconv.ParseFloat(vals[i], 32)
			require.NoError(t, err)

			d := float32(d64)
			require.Equalf(t, v, d,
				"expected f32s[%d] to be %f but got: %f", i, d, v,
			)
		}

		getF32S, erf := f.GetFloat32Slice("f32s")
		require.NoErrorf(t, erf,
			"got an error from GetFloat32Slice(): %v", erf,
		)

		for i, v := range getF32S {
			d64, err := strconv.ParseFloat(vals[i], 32)
			require.NoErrorf(t, err,
				"got an error from GetFloat32Slice(): %v", err,
			)

			require.Equalf(t, v, float32(d64),
				"expected f32s[%d] to be %f from GetFloat32Slice but got: %f", i, float32(d64), v,
			)
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"1.0", "2.0"}
		f32s := make([]float32, 0, len(vals))
		f := newFlagWithDefault(&f32s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--f32s=%s", strings.Join(vals, ",")),
		}))

		for i, v := range f32s {
			d64, err := strconv.ParseFloat(vals[i], 32)
			require.NoError(t, err)

			require.Equalf(t, v, float32(d64),
				"expected f32s[%d] to be %f but got: %f", i, float32(d64), v,
			)
		}

		getF32S, erf := f.GetFloat32Slice("f32s")
		require.NoErrorf(t, erf,
			"got an error from GetFloat32Slice(): %v", erf,
		)

		for i, v := range getF32S {
			d64, err := strconv.ParseFloat(vals[i], 32)
			require.NoError(t, err)

			require.Equalf(t, v, float32(d64),
				"expected f32s[%d] to be %f from GetFloat32Slice but got: %f", i, float32(d64), v,
			)
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--f32s=%s"
		in := []string{"1.0,2.0", "3.0"}
		f32s := make([]float32, 0, len(in))
		f := newFlag(&f32s)

		expected := []float32{1.0, 2.0, 3.0}
		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, f32s)
	})

	t.Run("as slice value", func(t *testing.T) {
		const argfmt = "--f32s=%s"
		in := []string{"1.0", "2.0"}
		f32s := make([]float32, 0, len(in))
		f := newFlag(&f32s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"3.1"}))
			}
		})

		require.Equalf(t, []float32{3.1}, f32s,
			"expected ss to be overwritten with '3.1', but got: %v", f32s,
		)
	})
}
