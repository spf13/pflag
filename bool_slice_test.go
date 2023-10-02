package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolSlice(t *testing.T) {
	t.Parallel()

	newFlag := func(bsp *[]bool) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.BoolSliceVar(bsp, "bs", []bool{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		bs := make([]bool, 0)
		f := newFlag(&bs)

		require.NoError(t, f.Parse([]string{}))

		getBS, err := f.GetBoolSlice("bs")
		require.NoErrorf(t, err,
			"got an error from GetBoolSlice(): %v", err,
		)

		require.Empty(t, getBS)
	})

	t.Run("with truthy/falsy values", func(t *testing.T) {
		vals := []string{"1", "F", "TRUE", "0"}
		bs := make([]bool, 0, len(vals))
		f := newFlag(&bs)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--bs=%s", strings.Join(vals, ",")),
		}))

		for i, v := range bs {
			b, err := strconv.ParseBool(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, b,
				"expected is[%d] to be %s but got: %t", i, vals[i], v,
			)
		}

		getBS, erb := f.GetBoolSlice("bs")
		require.NoError(t, erb)

		for i, v := range getBS {
			b, err := strconv.ParseBool(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, b,
				"expected bs[%d] to be %s but got: %t from GetBoolSlice", i, vals[i], v,
			)
		}
	})

	newFlagWithDefault := func(bsp *[]bool) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.BoolSliceVar(bsp, "bs", []bool{false, true}, "Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"false", "T"}
		bs := make([]bool, 0, len(vals))
		f := newFlagWithDefault(&bs)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range bs {
			b, err := strconv.ParseBool(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, b,
				"expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v,
			)
		}

		getBS, erb := f.GetBoolSlice("bs")
		require.NoErrorf(t, erb,
			"got an error from GetBoolSlice(): %v", erb,
		)

		for i, v := range getBS {
			b, err := strconv.ParseBool(vals[i])
			require.NoErrorf(t, err,
				"got an error from GetBoolSlice(): %v", err,
			)
			require.Equalf(t, v, b,
				"expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v,
			)
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"FALSE", "1"}
		bs := make([]bool, 0, len(vals))
		f := newFlagWithDefault(&bs)

		arg := fmt.Sprintf("--bs=%s", strings.Join(vals, ","))
		require.NoError(t, f.Parse([]string{arg}))

		for i, v := range bs {
			b, err := strconv.ParseBool(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, b,
				"expected bs[%d] to be %t but got: %t", i, b, v,
			)
		}

		getBS, erb := f.GetBoolSlice("bs")
		require.NoErrorf(t, erb,
			"got an error from GetBoolSlice(): %v", erb,
		)

		for i, v := range getBS {
			b, err := strconv.ParseBool(vals[i])
			require.NoError(t, err)
			require.Equalf(t, v, b,
				"expected bs[%d] to be %t from GetBoolSlice but got: %t", i, b, v,
			)
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--bs=%s"
		in := []string{"T,F", "T"}
		bs := make([]bool, 0, len(in))
		f := newFlag(&bs)
		expected := []bool{true, false, true}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, bs)
	})

	t.Run("as slice value", func(t *testing.T) {
		const argfmt = "--bs=%s"
		in := []string{"true", "false"}
		bs := make([]bool, 0, len(in))
		f := newFlag(&bs)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"false"}))
			}
		})

		require.Equalf(t, []bool{false}, bs,
			"expected ss to be overwritten with 'false', but got: %v", bs,
		)
	})

	t.Run("with quoting", func(t *testing.T) {
		tests := []struct {
			Want    []bool
			FlagArg []string
		}{
			{
				Want:    []bool{true, false, true},
				FlagArg: []string{"1", "0", "true"},
			},
			{
				Want:    []bool{true, false},
				FlagArg: []string{"True", "F"},
			},
			{
				Want:    []bool{true, false},
				FlagArg: []string{"T", "0"},
			},
			{
				Want:    []bool{true, false},
				FlagArg: []string{"1", "0"},
			},
			{
				Want:    []bool{true, false, false},
				FlagArg: []string{"true,false", "false"},
			},
			{
				Want:    []bool{true, false, false, true, false, true, false},
				FlagArg: []string{`"true,false,false,1,0,     T"`, " false "},
			},
			{
				Want:    []bool{false, false, true, false, true, false, true},
				FlagArg: []string{`"0, False,  T,false  , true,F"`, "true"},
			},
		}

		for i, test := range tests {
			bs := make([]bool, 0, 7)
			f := newFlag(&bs)

			require.NoErrorf(t,
				f.Parse([]string{fmt.Sprintf("--bs=%s", strings.Join(test.FlagArg, ","))}),
				"flag parsing failed for test %d with error:\nparsing:\t%#vnwant:\t\t%#v",
				test.FlagArg, test.Want,
			)

			require.Equalf(t, test.Want, bs, "on test %d", i)
		}
	})
}
