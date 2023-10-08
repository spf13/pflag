// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// This value can be a boolean ("true", "false") or "maybe"
type triStateValue int

const (
	triStateFalse triStateValue = 0
	triStateTrue  triStateValue = 1
	triStateMaybe triStateValue = 2
)

const strTriStateMaybe = "maybe"

func (v *triStateValue) IsBoolFlag() bool {
	return true
}

func (v *triStateValue) Get() interface{} {
	return *v
}

func (v *triStateValue) Set(s string) error {
	if s == strTriStateMaybe {
		*v = triStateMaybe
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = triStateTrue
	} else {
		*v = triStateFalse
	}

	return err
}

func (v *triStateValue) String() string {
	if *v == triStateMaybe {
		return strTriStateMaybe
	}
	return strconv.FormatBool(*v == triStateTrue)
}

// The type of the flag as required by the pflag.Value interface
func (v *triStateValue) Type() string {
	return "version"
}

func setUpFlagSet(tristate *triStateValue) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	*tristate = triStateFalse
	flag := f.VarPF(tristate, "tristate", "t", "tristate value (true, maybe or false)")
	flag.NoOptDefVal = "true"

	return f
}

func TestBool(t *testing.T) {
	t.Parallel()

	t.Run("with explicit true", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		require.NoError(t, f.Parse([]string{"--tristate=true"}))
		require.Equalf(t, triStateTrue, triState,
			"expected", triStateTrue, "(triStateTrue) but got", triState, "instead",
		)
	})

	t.Run("with implicit true", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		require.NoError(t, f.Parse([]string{"--tristate"}))
		require.Equalf(t, triStateTrue, triState,
			"expected", triStateTrue, "(triStateTrue) but got", triState, "instead",
		)
	})

	t.Run("with short flag", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		require.NoError(t, f.Parse([]string{"-t"}))
		require.Equalf(t, triStateTrue, triState,
			"expected", triStateTrue, "(triStateTrue) but got", triState, "instead",
		)
	})

	t.Run("with short flag extra argument", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		// The"maybe"turns into an arg, since short boolean options will only do true/false
		require.NoError(t, f.Parse([]string{"-t", "maybe"}))
		require.Equalf(t, triStateTrue, triState,
			"expected", triStateTrue, "(triStateTrue) but got", triState, "instead",
		)
		args := f.Args()
		require.Len(t, args, 1)
		require.Equalf(t, "maybe", args[0],
			"expected an extra 'maybe' argument to stick around",
		)
	})

	t.Run("with explicit maybe", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		require.NoError(t, f.Parse([]string{"--tristate=maybe"}))
		require.Equalf(t, triStateMaybe, triState,
			"expected", triStateMaybe, "(triStateMaybe) but got", triState, "instead",
		)
	})

	t.Run("with explicit false", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		require.NoError(t, f.Parse([]string{"--tristate=false"}))
		require.Equalf(t, triStateFalse, triState,
			"expected", triStateFalse, "(triStateFalse) but got", triState, "instead",
		)
	})

	t.Run("with implicit false", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		require.NoError(t, f.Parse([]string{}))
		require.Equalf(t, triStateFalse, triState,
			"expected", triStateFalse, "(triStateFalse) but got", triState, "instead",
		)
	})

	t.Run("with invalid value", func(t *testing.T) {
		var triState triStateValue
		f := setUpFlagSet(&triState)
		var buf bytes.Buffer
		f.SetOutput(&buf)
		require.Errorf(t, f.Parse([]string{"--tristate=invalid"}),
			"expected an error but did not get any, tristate has value", triState,
		)
	})

	t.Run("with BoolP", func(t *testing.T) {
		b := BoolP("bool", "b", false, "bool value in CommandLine")
		c := BoolP("c", "c", false, "other bool value")
		args := []string{"--bool"}
		require.NoError(t, CommandLine.Parse(args))
		require.Truef(t, *b,
			"expected b=true got b=%v", *b,
		)
		require.Falsef(t, *c,
			"expect c=false got c=%v", *c,
		)
	})
}
