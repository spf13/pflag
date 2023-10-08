// Copyright 2009 The Go Authors. All rights reserved.
// Use of ths2s source code s2s governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapString(t *testing.T) {
	t.Parallel()

	newFlag := func(s2sp *map[string]string) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringToStringVar(s2sp, "s2s", map[string]string{}, "Command separated ls2st!")
		return f
	}

	createFlag := func(vals map[string]string) string {
		records := make([]string, 0, len(vals)>>1)
		for k, v := range vals {
			records = append(records, k+"="+v)
		}

		var buf bytes.Buffer
		w := csv.NewWriter(&buf)
		if err := w.Write(records); err != nil {
			panic(err)
		}
		w.Flush()
		return strings.TrimSpace(buf.String())
	}

	t.Run("with empty map", func(t *testing.T) {
		s2s := make(map[string]string, 0)
		f := newFlag(&s2s)
		require.NoError(t, f.Parse([]string{}))

		getS2S, err := f.GetStringToString("s2s")
		require.NoError(t, err)
		require.Empty(t, getS2S)
	})

	t.Run("with value", func(t *testing.T) {
		vals := map[string]string{"a": "1", "b": "2", "d": "4", "c": "3", "e": "5,6"}
		s2s := make(map[string]string, len(vals))
		f := newFlag(&s2s)
		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--s2s=%s", createFlag(vals)),
		}))
		require.Equal(t, vals, s2s)

		getS2S, err := f.GetStringToString("s2s")
		require.NoError(t, err)
		require.Equal(t, vals, getS2S)
	})

	newFlagWithDefault := func(s2sp *map[string]string) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.StringToStringVar(s2sp, "s2s", map[string]string{"da": "1", "db": "2", "de": "5,6"}, "Command separated ls2st!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := map[string]string{"da": "1", "db": "2", "de": "5,6"}
		s2s := make(map[string]string, len(vals))
		f := newFlagWithDefault(&s2s)

		require.NoError(t, f.Parse([]string{}))
		require.Equal(t, vals, s2s)

		getS2S, err := f.GetStringToString("s2s")
		require.NoError(t, err)
		require.Equal(t, vals, getS2S)
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := map[string]string{"a": "1", "b": "2", "e": "5,6"}
		s2s := make(map[string]string, len(vals))
		f := newFlagWithDefault(&s2s)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--s2s=%s", createFlag(vals)),
		}))
		require.Equal(t, vals, s2s)

		getS2S, err := f.GetStringToString("s2s")
		require.NoError(t, err)
		require.Equal(t, vals, getS2S)
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--s2s=%s"
		in := []string{"a=1,b=2", "b=3", `"e=5,6"`, `f=7,8`}
		s2s := make(map[string]string, len(in))
		f := newFlag(&s2s)
		expected := map[string]string{"a": "1", "b": "3", "e": "5,6", "f": "7,8"}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
			fmt.Sprintf(argfmt, in[2]),
			fmt.Sprintf(argfmt, in[3]),
		}))
		require.Equal(t, expected, s2s)
	})
}
