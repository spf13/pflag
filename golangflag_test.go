// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	goflag "flag"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoflags(t *testing.T) {
	goflag.String("stringFlag", "stringFlag", "stringFlag")
	goflag.Bool("boolFlag", false, "boolFlag")
	f := NewFlagSet("test", ContinueOnError)

	f.AddGoFlagSet(goflag.CommandLine)
	require.NoError(t, f.Parse([]string{"--stringFlag=bob", "--boolFlag"}))

	getString, err := f.GetString("stringFlag")
	require.NoError(t, err)
	require.Equal(t, "bob", getString)

	getBool, err := f.GetBool("boolFlag")
	require.NoError(t, err)

	require.True(t, getBool)
	require.Truef(t, f.Parsed(),
		"f.Parsed() return false after f.Parse() called",
	)
}
