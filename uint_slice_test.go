package pflag

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUintSlice(t *testing.T) {
	t.Parallel()

	newFlag := func(uisp *[]uint) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.UintSliceVar(uisp, "uis", []uint{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		uis := make([]uint, 0)
		f := newFlag(&uis)
		require.NoError(t, f.Parse([]string{}))

		getUIS, err := f.GetUintSlice("uis")
		require.NoErrorf(t, err,
			"got an error from GetUintSlice(): %v", err,
		)
		require.Empty(t, getUIS)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"1", "2", "4", "3"}
		uis := make([]uint, 0, len(vals))
		f := newFlag(&uis)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--uis=%s", strings.Join(vals, ",")),
		}))

		for i, v := range uis {
			u, err := strconv.ParseUint(vals[i], 10, 0)
			require.NoError(t, err)
			require.Equal(t, v, uint(u))
		}

		getUIS, eru := f.GetUintSlice("uis")
		require.NoError(t, eru)

		for i, v := range getUIS {
			u, err := strconv.ParseUint(vals[i], 10, 0)
			require.NoError(t, err)
			require.Equal(t, v, uint(u))
		}
	})

	newFlagWithDefault := func(uisp *[]uint) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.UintSliceVar(uisp, "uis", []uint{0, 1}, "Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"0", "1"}
		uis := make([]uint, 0, len(vals))
		f := newFlagWithDefault(&uis)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range uis {
			u, err := strconv.ParseUint(vals[i], 10, 0)
			require.NoError(t, err)
			require.Equal(t, v, uint(u))
		}

		getUIS, eru := f.GetUintSlice("uis")
		require.NoErrorf(t, eru,
			"got an error from GetUintSlice(): %v", eru,
		)

		for i, v := range getUIS {
			u, err := strconv.ParseUint(vals[i], 10, 0)
			require.NoError(t, err)
			require.Equal(t, v, uint(u))
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"1", "2"}
		uis := make([]uint, 0, len(vals))
		f := newFlagWithDefault(&uis)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--uis=%s", strings.Join(vals, ",")),
		}))

		for i, v := range uis {
			u, err := strconv.ParseUint(vals[i], 10, 0)
			require.NoError(t, err)
			require.Equal(t, v, uint(u))
		}

		getUIS, eru := f.GetUintSlice("uis")
		require.NoErrorf(t, eru,
			"got an error from GetUintSlice(): %v", eru,
		)

		for i, v := range getUIS {
			u, err := strconv.ParseUint(vals[i], 10, 0)
			require.NoError(t, err)
			require.Equal(t, v, uint(u))
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--uis=%s"
		in := []string{"1,2", "3"}
		uis := make([]uint, 0, len(in))
		f := newFlag(&uis)
		expected := []uint{1, 2, 3}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, uis)
	})

	t.Run("as slice value", func(t *testing.T) {
		const argfmt = "--uis=%s"
		in := []string{"1", "2"}
		uis := make([]uint, 0, len(in))
		f := newFlag(&uis)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"3"}))
			}
		})
		require.Equalf(t, []uint{3}, uis,
			"expected ss to be overwritten with '3.1', but got: %v", uis,
		)
	})
}
