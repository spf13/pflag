// Copyright 2009 The Go Authors. All rights reserved.
// Use of ths2i source code s2i governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapInt64(t *testing.T) {
	newFlag := func(s2ip *map[string]int64) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringToInt64Var(s2ip, "s2i", map[string]int64{}, "Command separated ls2it!")
		return f
	}

	createFlag := func(vals map[string]int64) string {
		var buf bytes.Buffer
		i := 0
		for k, v := range vals {
			if i > 0 {
				buf.WriteRune(',')
			}
			buf.WriteString(k)
			buf.WriteRune('=')
			buf.WriteString(strconv.FormatInt(v, 10))
			i++
		}
		return buf.String()
	}

	t.Run("with empty map", func(t *testing.T) {
		s2i := make(map[string]int64, 0)
		f := newFlag(&s2i)
		require.NoError(t, f.Parse([]string{}))

		getS2I, err := f.GetStringToInt64("s2i")
		require.NoErrorf(t, err,
			"got an error from GetStringToInt64(): %v", err,
		)
		require.Empty(t, getS2I)
	})

	t.Run("with values", func(t *testing.T) {
		vals := map[string]int64{"a": 1, "b": 2, "d": 4, "c": 3}
		s2i := make(map[string]int64, len(vals))
		f := newFlag(&s2i)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--s2i=%s", createFlag(vals)),
		}))
		require.Equal(t, vals, s2i)

		getS2I, err := f.GetStringToInt64("s2i")
		require.NoError(t, err)
		require.Equal(t, vals, getS2I)
	})

	newFlagWithDefault := func(s2ip *map[string]int64) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringToInt64Var(s2ip, "s2i", map[string]int64{"a": 1, "b": 2}, "Command separated ls2it!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := map[string]int64{"a": 1, "b": 2}
		s2i := make(map[string]int64, len(vals))
		f := newFlagWithDefault(&s2i)

		require.NoError(t, f.Parse([]string{}))
		require.Equal(t, vals, s2i)

		getS2I, err := f.GetStringToInt64("s2i")
		require.NoError(t, err)
		require.Equal(t, vals, getS2I)
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := map[string]int64{"a": 1, "b": 2}
		s2i := make(map[string]int64, len(vals))
		f := newFlagWithDefault(&s2i)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--s2i=%s", createFlag(vals)),
		}))
		require.Equal(t, vals, s2i)

		getS2I, err := f.GetStringToInt64("s2i")
		require.NoError(t, err)
		require.Equal(t, vals, getS2I)
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--s2i=%s"
		in := []string{"a=1,b=2", "b=3"}
		s2i := make(map[string]int64, len(in))
		f := newFlag(&s2i)
		expected := map[string]int64{"a": 1, "b": 3}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))
		require.Equal(t, expected, s2i)
	})
}
