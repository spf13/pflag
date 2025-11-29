// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag_test

import (
	"fmt"

	"github.com/spf13/pflag"
)

func ExampleShorthandLookup() {
	name := "verbose"
	short := name[:1]

	pflag.BoolP(name, short, false, "verbose output")

	// len(short) must be == 1
	flag := pflag.ShorthandLookup(short)

	fmt.Println(flag.Name)
}

func ExampleFlagSet_ShorthandLookup() {
	name := "verbose"
	short := name[:1]

	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)
	fs.BoolP(name, short, false, "verbose output")

	// len(short) must be == 1
	flag := fs.ShorthandLookup(short)

	fmt.Println(flag.Name)
}

func ExampleFlagSet_StringToString() {
	args := []string{
		"--arg", "a=1,b=2",
		"--arg", "a=2",
		"--arg=d=4",
	}

	fs := pflag.NewFlagSet("Example", pflag.ContinueOnError)
	fs.StringToString("arg", make(map[string]string), "string-to-string arg accepting key=value pairs")

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	value, err := fs.GetStringToString("arg")
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
	// Output: map[a:2 b:2 d:4]
}
