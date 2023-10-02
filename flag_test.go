// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testBool        *bool
	testInt         *int
	testInt64       *int64
	testUint        *uint
	testUint64      *uint64
	testString      *string
	testFloat       *float64
	testDuration    *time.Duration
	testOptionalInt *int

	normalizeFlagNameInvocations = 0
)

func init() {
	testBool = Bool("test_bool", false, "bool value")
	testInt = Int("test_int", 0, "int value")
	testInt64 = Int64("test_int64", 0, "int64 value")
	testUint = Uint("test_uint", 0, "uint value")
	testUint64 = Uint64("test_uint64", 0, "uint64 value")
	testString = String("test_string", "0", "string value")
	testFloat = Float64("test_float64", 0, "float64 value")
	testDuration = Duration("test_duration", 0, "time.Duration value")
	testOptionalInt = Int("test_optional_int", 0, "optional int value")
}

func TestVisit(t *testing.T) {
	boolString := func(s string) string {
		if s == "0" {
			return "false"
		}
		return "true"
	}

	visitor := func(desired string, m map[string]*Flag) func(*Flag) {
		return func(f *Flag) {
			if len(f.Name) <= 5 || f.Name[0:5] != "test_" {
				return
			}

			m[f.Name] = f
			ok := false

			switch {
			case f.Value.String() == desired:
				ok = true
			case f.Name == "test_bool" && f.Value.String() == boolString(desired):
				ok = true
			case f.Name == "test_duration" && f.Value.String() == desired+"s":
				ok = true
			}
			require.Truef(t, ok,
				"visit: bad value", f.Value.String(), "for", f.Name,
			)
		}
	}

	printMap := func(m map[string]*Flag) {
		for k, v := range m {
			t.Log(k, *v)
		}
	}

	t.Run("with VisitAll", func(t *testing.T) {
		const desired = "0"
		m := make(map[string]*Flag)

		VisitAll(visitor(desired, m))
		if !assert.Lenf(t, m, 9, "VisitAll misses some flags") {
			printMap(m)
		}
	})

	t.Run("with Visit", func(t *testing.T) {
		const desired = "0"
		m := make(map[string]*Flag)

		Visit(visitor(desired, m))
		if !assert.Lenf(t, m, 0, "Visit sees unset flags") {
			printMap(m)
		}
	})

	t.Run("with all flags set", func(t *testing.T) {
		const desired = "1"
		m := make(map[string]*Flag)

		require.NoError(t, Set("test_bool", "true"))
		require.NoError(t, Set("test_int", "1"))
		require.NoError(t, Set("test_int64", "1"))
		require.NoError(t, Set("test_uint", "1"))
		require.NoError(t, Set("test_uint64", "1"))
		require.NoError(t, Set("test_string", "1"))
		require.NoError(t, Set("test_float64", "1"))
		require.NoError(t, Set("test_duration", "1s"))
		require.NoError(t, Set("test_optional_int", "1"))

		Visit(visitor(desired, m))
		if !assert.Lenf(t, m, 9, "Visit fails after set") {
			printMap(m)
		}
	})

	t.Run("visit in sorted order", func(t *testing.T) {
		var flagNames []string
		Visit(func(f *Flag) { flagNames = append(flagNames, f.Name) })
		require.Truef(t, sort.StringsAreSorted(flagNames),
			"flag names not sorted: %v", flagNames,
		)
	})
}

func TestUsage(t *testing.T) {
	called := false
	ResetForTesting(func() { called = true })

	require.NotNilf(t, GetCommandLine().Parse([]string{"--x"}),
		"parse did not fail for unknown flag",
	)
	require.Falsef(t, called,
		"did call Usage while using ContinueOnError",
	)
}

func TestAddFlagSet(t *testing.T) {
	oldSet := NewFlagSet("old", ContinueOnError)
	newSet := NewFlagSet("new", ContinueOnError)

	oldSet.String("flag1", "flag1", "flag1")
	oldSet.String("flag2", "flag2", "flag2")

	newSet.String("flag2", "flag2", "flag2")
	newSet.String("flag3", "flag3", "flag3")

	oldSet.AddFlagSet(newSet)

	require.Lenf(t, oldSet.formal, 3,
		"unexpected result adding a FlagSet to a FlagSet %v", oldSet,
	)
}

func TestAnnotation(t *testing.T) {
	f := NewFlagSet("shorthand", ContinueOnError)

	require.Errorf(t, f.SetAnnotation("missing-flag", "key", nil),
		"expected error setting annotation on non-existent flag",
	)

	f.StringP("stringa", "a", "", "string value")
	require.NoErrorf(t, f.SetAnnotation("stringa", "key", nil),
		"unexpected error setting new nil annotation",
	)
	require.Nil(t, f.Lookup("stringa").Annotations["key"],
		"unexpected annotation",
	)

	f.StringP("stringb", "b", "", "string2 value")
	require.NoErrorf(t, f.SetAnnotation("stringb", "key", []string{"value1"}),
		"unexpected error setting new annotation",
	)

	annotation := f.Lookup("stringb").Annotations["key"]
	require.EqualValuesf(t, []string{"value1"}, annotation,
		"unexpected annotation: %v", annotation,
	)

	require.NoErrorf(t, f.SetAnnotation("stringb", "key", []string{"value2"}),
		"unexpected error updating annotation",
	)
	annotation = f.Lookup("stringb").Annotations["key"]
	require.EqualValuesf(t, []string{"value2"}, annotation,
		"unexpected annotation: %v", annotation,
	)
}

func TestName(t *testing.T) {
	const flagSetName = "bob"
	f := NewFlagSet(flagSetName, ContinueOnError)

	givenName := f.Name()
	require.Equalf(t, flagSetName, givenName,
		"unexpected result when retrieving a FlagSet's name: expected %s, but found %s",
		flagSetName, givenName,
	)
}

func testParse(f *FlagSet, t *testing.T) {
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")

	boolFlag := f.Bool("bool", false, "bool value")
	bool2Flag := f.Bool("bool2", false, "bool2 value")
	bool3Flag := f.Bool("bool3", false, "bool3 value")
	intFlag := f.Int("int", 0, "int value")
	int8Flag := f.Int8("int8", 0, "int value")
	int16Flag := f.Int16("int16", 0, "int value")
	int32Flag := f.Int32("int32", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint8Flag := f.Uint8("uint8", 0, "uint value")
	uint16Flag := f.Uint16("uint16", 0, "uint value")
	uint32Flag := f.Uint32("uint32", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	float32Flag := f.Float32("float32", 0, "float32 value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	ipFlag := f.IP("ip", net.ParseIP("127.0.0.1"), "ip value")
	maskFlag := f.IPMask("mask", ParseIPv4Mask("0.0.0.0"), "mask value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")
	optionalIntNoValueFlag := f.Int("optional-int-no-value", 0, "int value")

	f.Lookup("optional-int-no-value").NoOptDefVal = "9"
	optionalIntWithValueFlag := f.Int("optional-int-with-value", 0, "int value")
	f.Lookup("optional-int-no-value").NoOptDefVal = "9"

	const extra = "one-extra-argument"

	t.Run("parse args", func(t *testing.T) {
		args := []string{
			"--bool",
			"--bool2=true",
			"--bool3=false",
			"--int=22",
			"--int8=-8",
			"--int16=-16",
			"--int32=-32",
			"--int64=0x23",
			"--uint", "24",
			"--uint8=8",
			"--uint16=16",
			"--uint32=32",
			"--uint64=25",
			"--string=hello",
			"--float32=-172e12",
			"--float64=2718e28",
			"--ip=10.11.12.13",
			"--mask=255.255.255.0",
			"--duration=2m",
			"--optional-int-no-value",
			"--optional-int-with-value=42",
			extra,
		}
		require.NoError(t, f.Parse(args))
		require.Truef(t, f.Parsed(), "f.Parse() = false after Parse")
	})

	t.Run("with bool flags", func(t *testing.T) {
		require.Truef(t, *boolFlag,
			"bool flag should be true, is ", *boolFlag,
		)

		v, err := f.GetBool("bool")
		require.NoError(t, err)
		require.Equalf(t, *boolFlag, v, "GetBool does not work")
		require.Truef(t, *bool2Flag,
			"bool2 flag should be true, is ", *bool2Flag,
		)
		require.Falsef(t, *bool3Flag,
			"bool3 flag should be false, is ", *bool2Flag,
		)
	})

	t.Run("with integer flags", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			require.Equalf(t, 22, *intFlag,
				"int flag should be 22, is ", *intFlag,
			)
			v, err := f.GetInt("int")
			require.NoError(t, err)
			require.Equalf(t, *intFlag, v, "GetInt does not work")
		})

		t.Run("int8", func(t *testing.T) {
			require.Equalf(t, int8(-8), *int8Flag,
				"int8 flag should be 0x23, is ", *int8Flag,
			)
			v, err := f.GetInt8("int8")
			require.NoError(t, err)
			require.Equalf(t, *int8Flag, v, "GetInt8 does not work")
		})

		t.Run("int16", func(t *testing.T) {
			require.Equalf(t, int16(-16), *int16Flag,
				"int16 flag should be -16, is ", *int16Flag,
			)
			v, err := f.GetInt16("int16")
			require.NoError(t, err)
			require.Equalf(t, *int16Flag, v, "GetInt16 does not work")
		})

		t.Run("int32", func(t *testing.T) {
			require.Equalf(t, int32(-32), *int32Flag,
				"int32 flag should be 0x23, is ", *int32Flag,
			)
			v, err := f.GetInt32("int32")
			require.NoError(t, err)
			require.Equalf(t, *int32Flag, v, "GetInt32 does not work")
		})

		t.Run("int64", func(t *testing.T) {
			require.Equalf(t, int64(0x23), *int64Flag,
				"int64 flag should be 0x23, is ", *int64Flag,
			)
			v, err := f.GetInt64("int64")
			require.NoError(t, err)
			require.Equalf(t, *int64Flag, v, "GetInt64 does not work")
		})

		t.Run("uint", func(t *testing.T) {
			require.Equalf(t, uint(24), *uintFlag,
				"uint flag should be 24, is ", *uintFlag,
			)
			v, err := f.GetUint("uint")
			require.NoError(t, err)
			require.Equalf(t, *uintFlag, v, "GetUint does not work")
		})

		t.Run("uint8", func(t *testing.T) {
			require.Equalf(t, uint8(8), *uint8Flag,
				"uint8 flag should be 8, is ", *uint8Flag,
			)
			v, err := f.GetUint8("uint8")
			require.NoError(t, err)
			require.Equalf(t, *uint8Flag, v, "GetUint8 does not work")
		})

		t.Run("uint16", func(t *testing.T) {
			require.Equalf(t, uint16(16), *uint16Flag,
				"uint16 flag should be 16, is ", *uint16Flag,
			)
			v, err := f.GetUint16("uint16")
			require.NoError(t, err)
			require.Equalf(t, *uint16Flag, v, "GetUint16 does not work")
		})

		t.Run("uint32", func(t *testing.T) {
			require.Equalf(t, uint32(32), *uint32Flag,
				"uint32 flag should be 32, is ", *uint32Flag,
			)
			v, err := f.GetUint32("uint32")
			require.NoError(t, err)
			require.Equalf(t, *uint32Flag, v, "GetUint32 does not work")
		})

		t.Run("uint64", func(t *testing.T) {
			require.Equalf(t, uint64(25), *uint64Flag,
				"uint64 flag should be 25, is ", *uint64Flag,
			)
			v, err := f.GetUint64("uint64")
			require.NoError(t, err)
			require.Equalf(t, *uint64Flag, v, "GetUint64 does not work")
		})
	})

	t.Run("with string flags", func(t *testing.T) {
		require.Equalf(t, "hello", *stringFlag,
			"string flag should be `hello`, is ", *stringFlag,
		)
		v, err := f.GetString("string")
		require.NoError(t, err)
		require.Equalf(t, *stringFlag, v, "GetString does not work")
	})

	t.Run("with float flags", func(t *testing.T) {
		t.Run("float32", func(t *testing.T) {
			require.Equalf(t, float32(-172e12), *float32Flag,
				"float32 flag should be -172e12, is ", *float32Flag,
			)
			v, err := f.GetFloat32("float32")
			require.NoError(t, err)
			require.Equalf(t, *float32Flag, v, "GetFloat32 returned %v but float32Flag was %v", v, *float32Flag)
		})

		t.Run("float64", func(t *testing.T) {
			require.Equalf(t, 2718e28, *float64Flag,
				"float64 flag should be 2718e28, is ", *float64Flag,
			)
			v, err := f.GetFloat64("float64")
			require.NoError(t, err)
			require.Equalf(t, *float64Flag, v, "GetFloat64 returned %v but float64Flag was %v", v, *float64Flag)
		})
	})

	t.Run("with IP address flags", func(t *testing.T) {
		t.Run("IP", func(t *testing.T) {
			require.True(t, ipFlag.Equal(net.ParseIP("10.11.12.13")),
				"ip flag should be 10.11.12.13, is ", *ipFlag,
			)
			v, err := f.GetIP("ip")
			require.NoError(t, err)
			require.True(t, v.Equal(*ipFlag),
				"GetIP returned %v but ipFlag was %v", v, *ipFlag,
			)
		})

		t.Run("IPv4Mask", func(t *testing.T) {
			require.Equal(t, ParseIPv4Mask("255.255.255.0").String(), maskFlag.String(),
				"mask flag should be 255.255.255.0, is ", maskFlag.String(),
			)
			v, err := f.GetIPv4Mask("mask")
			require.NoError(t, err)
			require.Equal(t, maskFlag.String(), v.String(),
				"GetIP returned %v maskFlag was %v", v, *maskFlag,
			)
		})
	})

	t.Run("with duration flags", func(t *testing.T) {
		require.Equalf(t, 2*time.Minute, *durationFlag,
			"duration flag should be 2m, is ", *durationFlag,
		)
		v, err := f.GetDuration("duration")
		require.NoError(t, err)
		require.Equalf(t, *durationFlag, v, "GetDuration does not work")

		_, err = f.GetInt("duration")
		require.Errorf(t, err, "unexpectedly, GetInt parsed a time.Duration")
	})

	t.Run("flags with no-value defaults", func(t *testing.T) {
		require.Equalf(t, 9, *optionalIntNoValueFlag,
			"optional int flag should be the default value, is ", *optionalIntNoValueFlag,
		)
		require.Equalf(t, 42, *optionalIntWithValueFlag,
			"optional int flag should be 42, is ", *optionalIntWithValueFlag,
		)
	})

	t.Run("with non-flag argument", func(t *testing.T) {
		require.Lenf(t, f.Args(), 1,
			"expected one argument, got", len(f.Args()),
		)
		require.Equalf(t, extra, f.Args()[0],
			"expected argument %q got %q", extra, f.Args()[0],
		)
	})
}

func testParseAll(f *FlagSet, t *testing.T) {
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")

	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")
	f.BoolP("boolc", "c", false, "bool3 value")
	f.BoolP("boold", "d", false, "bool4 value")
	f.StringP("stringa", "s", "0", "string value")
	f.StringP("stringz", "z", "0", "string value")
	f.StringP("stringx", "x", "0", "string value")
	f.StringP("stringy", "y", "0", "string value")
	f.Lookup("stringx").NoOptDefVal = "1"

	args := []string{
		"-ab",
		"-cs=xx",
		"--stringz=something",
		"-d=true",
		"-x",
		"-y",
		"ee",
	}

	want := []string{
		"boola", "true",
		"boolb", "true",
		"boolc", "true",
		"stringa", "xx",
		"stringz", "something",
		"boold", "true",
		"stringx", "1",
		"stringy", "ee",
	}
	got := make([]string, 0, len(want))

	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}

	require.NoError(t, f.ParseAll(args, store))
	require.Truef(t, f.Parsed(), "f.Parse() = false after Parse")
	require.Equalf(t, want, got,
		"f.ParseAll() fail to restore the args. Got: %v, Want: %v",
		got, want,
	)
}

func testParseWithUnknownFlags(f *FlagSet, t *testing.T) {
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")
	f.ParseErrorsWhitelist.UnknownFlags = true

	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")
	f.BoolP("boolc", "c", false, "bool3 value")
	f.BoolP("boold", "d", false, "bool4 value")
	f.BoolP("boole", "e", false, "bool4 value")
	f.StringP("stringa", "s", "0", "string value")
	f.StringP("stringz", "z", "0", "string value")
	f.StringP("stringx", "x", "0", "string value")
	f.StringP("stringy", "y", "0", "string value")
	f.StringP("stringo", "o", "0", "string value")

	f.Lookup("stringx").NoOptDefVal = "1"

	args := []string{
		"-ab",
		"-cs=xx",
		"--stringz=something",
		"--unknown1",
		"unknown1Value",
		"-d=true",
		"-x",
		"--unknown2=unknown2Value",
		"-u=unknown3Value",
		"-p",
		"unknown4Value",
		"-q", // another unknown with bool value
		"-y",
		"ee",
		"--unknown7=unknown7value",
		"--stringo=ovalue",
		"--unknown8=unknown8value",
		"--boole",
		"--unknown6",
		"",
		"-uuuuu",
		"",
		"--unknown10",
		"--unknown11",
	}

	want := []string{
		"boola", "true",
		"boolb", "true",
		"boolc", "true",
		"stringa", "xx",
		"stringz", "something",
		"boold", "true",
		"stringx", "1",
		"stringy", "ee",
		"stringo", "ovalue",
		"boole", "true",
	}
	got := make([]string, 0, len(want))

	store := func(flag *Flag, value string) error {
		got = append(got, flag.Name)
		if len(value) > 0 {
			got = append(got, value)
		}
		return nil
	}

	require.NoError(t, f.ParseAll(args, store))
	require.Truef(t, f.Parsed(), "f.Parse() = false after Parse")
	require.Equalf(t, want, got,
		"f.ParseAll() fail to restore the args. Got: %v, Want: %v",
		got, want,
	)
}

func TestShorthand(t *testing.T) {
	f := NewFlagSet("shorthand", ContinueOnError)
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")

	boolaFlag := f.BoolP("boola", "a", false, "bool value")
	boolbFlag := f.BoolP("boolb", "b", false, "bool2 value")
	boolcFlag := f.BoolP("boolc", "c", false, "bool3 value")
	booldFlag := f.BoolP("boold", "d", false, "bool4 value")
	stringaFlag := f.StringP("stringa", "s", "0", "string value")
	stringzFlag := f.StringP("stringz", "z", "0", "string value")

	const extra = "interspersed-argument"
	const notaflag = "--i-look-like-a-flag"

	args := []string{
		"-ab",
		extra,
		"-cs",
		"hello",
		"-z=something",
		"-d=true",
		"--",
		notaflag,
	}

	f.SetOutput(ioutil.Discard)
	require.NoError(t, f.Parse(args))
	require.Truef(t, f.Parsed(), "f.Parse() = false after Parse")

	require.Truef(t, *boolaFlag, "boola flag should be true, is ", *boolaFlag)
	require.Truef(t, *boolbFlag, "boolb flag should be true, is ", *boolbFlag)
	require.Truef(t, *boolcFlag, "boolc flag should be true, is ", *boolcFlag)
	require.Truef(t, *booldFlag, "boold flag should be true, is ", *booldFlag)
	require.Equalf(t, "hello", *stringaFlag, "stringa flag should be `hello`, is ", *stringaFlag)
	require.Equalf(t, "something", *stringzFlag, "stringz flag should be `something`, is ", *stringzFlag)

	require.Len(t, f.Args(), 2, "expected one argument, got", len(f.Args()))
	require.Equalf(t, extra, f.Args()[0], "expected argument %q got %q", extra, f.Args()[0])
	require.Equalf(t, notaflag, f.Args()[1], "expected argument %q got %q", notaflag, f.Args()[1])
	require.Equal(t, 1, f.ArgsLenAtDash(), "expected argsLenAtDash %d got %d", f.ArgsLenAtDash(), 1)
}

func TestShorthandLookup(t *testing.T) {
	f := NewFlagSet("shorthand", ContinueOnError)
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")

	f.BoolP("boola", "a", false, "bool value")
	f.BoolP("boolb", "b", false, "bool2 value")

	args := []string{
		"-ab",
	}

	f.SetOutput(ioutil.Discard)
	require.NoError(t, f.Parse(args))
	require.Truef(t, f.Parsed(), "f.Parse() = false after Parse")

	flag := f.ShorthandLookup("a")
	require.NotNil(t, flag, "f.ShorthandLookup(\"a\") returned nil")

	require.Equalf(t, "boola", flag.Name,
		"f.ShorthandLookup(\"a\") found %q instead of \"boola\"", flag.Name,
	)
	require.Nil(t, f.ShorthandLookup(""),
		"f.ShorthandLookup(\"\") did not return nil",
	)
	require.Panicsf(t, func() { _ = f.ShorthandLookup("ab") },
		"f.ShorthandLookup(\"ab\") did not panic",
	)
}

func TestParse(t *testing.T) {
	ResetForTesting(func() { t.Error("bad parse") })
	testParse(GetCommandLine(), t)
}

func TestFlagSetParse(t *testing.T) {
	testParse(NewFlagSet("test", ContinueOnError), t)
}

func TestParseAll(t *testing.T) {
	ResetForTesting(func() { t.Error("bad parse") })
	testParseAll(GetCommandLine(), t)
}

func TestIgnoreUnknownFlags(t *testing.T) {
	ResetForTesting(func() { t.Error("bad parse") })
	testParseWithUnknownFlags(GetCommandLine(), t)
}

func TestChangedHelper(t *testing.T) {
	f := NewFlagSet("changedtest", ContinueOnError)

	f.Bool("changed", false, "changed bool")
	f.Bool("settrue", true, "true to true")
	f.Bool("setfalse", false, "false to false")
	f.Bool("unchanged", false, "unchanged bool")

	args := []string{"--changed", "--settrue", "--setfalse=false"}

	require.NoError(t, f.Parse(args))
	require.Truef(t, f.Changed("changed"), "--changed wasn't changed!")
	require.Truef(t, f.Changed("settrue"), "--settrue wasn't changed!")
	require.Truef(t, f.Changed("setfalse"), "--setfalse wasn't changed!")
	require.Falsef(t, f.Changed("unchanged"), "--unchanged was changed!")
	require.Falsef(t, f.Changed("invalid"), "--invalid was changed!")

	require.Equalf(t, -1, f.ArgsLenAtDash(),
		"expected argsLenAtDash: %d but got %d", -1, f.ArgsLenAtDash(),
	)
}

func replaceSeparators(name string, from []string, to string) string { //nolint: unparam
	result := name
	for _, sep := range from {
		result = strings.ReplaceAll(result, sep, to)
	}
	// Type convert to indicate normalization has been done.
	return result
}

func wordSepNormalizeFunc(_ *FlagSet, name string) NormalizedName {
	seps := []string{"-", "_"}
	name = replaceSeparators(name, seps, ".")
	normalizeFlagNameInvocations++

	return NormalizedName(name)
}

func testWordSepNormalizedNames(args []string, t *testing.T) {
	f := NewFlagSet("normalized", ContinueOnError)
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")

	withDashFlag := f.Bool("with-dash-flag", false, "bool value")
	// Set this after some flags have been added and before others.
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	withUnderFlag := f.Bool("with_under_flag", false, "bool value")
	withBothFlag := f.Bool("with-both_flag", false, "bool value")

	require.NoError(t, f.Parse(args))
	require.Truef(t, f.Parsed(), "f.Parse() = false after Parse")

	require.Truef(t, *withDashFlag, "withDashFlag flag should be true, is ", *withDashFlag)
	require.Truef(t, *withUnderFlag, "withUnderFlag flag should be true, is ", *withUnderFlag)
	require.Truef(t, *withBothFlag, "withBothFlag flag should be true, is ", *withBothFlag)
}

func TestWordSepNormalizedNames(t *testing.T) {
	t.Run("with dashes", func(t *testing.T) {
		args := []string{
			"--with-dash-flag",
			"--with-under-flag",
			"--with-both-flag",
		}
		testWordSepNormalizedNames(args, t)
	})

	t.Run("with underscores", func(t *testing.T) {
		args := []string{
			"--with_dash_flag",
			"--with_under_flag",
			"--with_both_flag",
		}
		testWordSepNormalizedNames(args, t)
	})

	t.Run("with dash and underscores", func(t *testing.T) {
		args := []string{
			"--with-dash_flag",
			"--with-under_flag",
			"--with-both_flag",
		}
		testWordSepNormalizedNames(args, t)
	})
}

func aliasAndWordSepFlagNames(_ *FlagSet, name string) NormalizedName {
	seps := []string{"-", "_"}

	oldName := replaceSeparators("old-valid_flag", seps, ".")
	newName := replaceSeparators("valid-flag", seps, ".")

	name = replaceSeparators(name, seps, ".")
	if name == oldName {
		name = newName
	}

	return NormalizedName(name)
}

func TestCustomNormalizedNames(t *testing.T) {
	f := NewFlagSet("normalized", ContinueOnError)
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")

	validFlag := f.Bool("valid-flag", false, "bool value")
	f.SetNormalizeFunc(aliasAndWordSepFlagNames)
	someOtherFlag := f.Bool("some-other-flag", false, "bool value")

	args := []string{"--old_valid_flag", "--some-other_flag"}
	require.NoError(t, f.Parse(args))

	require.Truef(t, *validFlag, "validFlag is %v even though we set the alias --old_valid_flag", *validFlag)
	require.Truef(t, *someOtherFlag, "someOtherFlag should be true, is ", *someOtherFlag)
}

// Every flag we add, the name (displayed also in usage) should be normalized
func TestNormalizationFuncShouldChangeFlagName(t *testing.T) {
	t.Run("with normalization after addition", func(t *testing.T) {
		f := NewFlagSet("normalized", ContinueOnError)

		f.Bool("valid_flag", false, "bool value")
		require.Equalf(t, "valid_flag", f.Lookup("valid_flag").Name,
			"the new flag should have the name 'valid_flag' instead of ", f.Lookup("valid_flag").Name,
		)

		f.SetNormalizeFunc(wordSepNormalizeFunc)
		require.Equalf(t, "valid.flag", f.Lookup("valid_flag").Name,
			"the new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name,
		)
	})

	t.Run("with normalization before addition", func(t *testing.T) {
		f := NewFlagSet("normalized", ContinueOnError)
		f.SetNormalizeFunc(wordSepNormalizeFunc)

		f.Bool("valid_flag", false, "bool value")
		require.Equalf(t, "valid.flag", f.Lookup("valid_flag").Name,
			"the new flag should have the name 'valid.flag' instead of ", f.Lookup("valid_flag").Name,
		)
	})
}

// Related to https://github.com/spf13/cobra/issues/521.
func TestNormalizationSharedFlags(t *testing.T) {
	f := NewFlagSet("set f", ContinueOnError)
	g := NewFlagSet("set g", ContinueOnError)

	const testName = "valid_flag"
	nfunc := wordSepNormalizeFunc
	normName := nfunc(nil, testName)
	require.NotEqualf(t, string(normName), testName,
		"TestNormalizationSharedFlags meaningless: the original and normalized flag names are identical:", testName,
	)

	f.Bool(testName, false, "bool value")
	g.AddFlagSet(f)

	f.SetNormalizeFunc(nfunc)
	g.SetNormalizeFunc(nfunc)

	require.Lenf(t, f.formal, 1,
		"normalizing flags should not result in duplications in the flag set:", f.formal,
	)
	require.Equalf(t, string(normName), f.orderedFormal[0].Name,
		"flag name not normalized",
	)

	for k := range f.formal {
		require.Equalf(t, "valid.flag", string(k),
			"the key in the flag map should have been normalized: wanted \"%s\", got \"%s\" instead", normName, k,
		)
	}

	require.Equalf(t, g.formal, f.formal,
		"two flag sets sharing the same flags should stay consistent after being normalized. Original set:",
		f.formal, "Duplicate set:", g.formal,
	)
	require.Equalf(t, g.orderedFormal, f.orderedFormal,
		"two ordered flag sets sharing the same flags should stay consistent after being normalized. Original set:",
		f.formal, "Duplicate set:", g.formal,
	)
}

func TestNormalizationSetFlags(t *testing.T) {
	f := NewFlagSet("normalized", ContinueOnError)
	nfunc := wordSepNormalizeFunc
	const testName = "valid_flag"
	normName := nfunc(nil, testName)

	require.NotEqualf(t, string(normName), testName,
		"TestNormalizationSetFlags meaningless: the original and normalized flag names are identical:", testName,
	)

	f.Bool(testName, false, "bool value")
	require.NoError(t, f.Set(testName, "true"))
	f.SetNormalizeFunc(nfunc)

	require.Lenf(t, f.formal, 1,
		"normalizing flags should not result in duplications in the flag set:", f.formal,
	)
	require.Equalf(t, string(normName), f.orderedFormal[0].Name,
		"flag name not normalized",
	)

	for k := range f.formal {
		require.Equalf(t, "valid.flag", string(k),
			"the key in the flag map should have been normalized: wanted \"%s\", got \"%s\" instead", normName, k,
		)
	}

	require.Equalf(t, f.actual, f.formal,
		"the map of set flags should get normalized. Formal:",
		f.formal, "Actual:", f.actual,
	)
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *flagVar) Type() string {
	return "flagVar"
}

func TestUserDefined(t *testing.T) {
	var (
		flags FlagSet
		v     flagVar
	)

	flags.Init("test", ContinueOnError)
	flags.VarP(&v, "v", "v", "usage")

	require.NoError(t, flags.Parse([]string{"--v=1", "-v2", "-v", "3"}))
	require.Lenf(t, v, 3, "expected 3 args; got ", len(v))

	const expect = "[1 2 3]"
	require.Equalf(t, expect, v.String(),
		"expected value %q got %q", expect, v.String(),
	)
}

func TestSetOutput(t *testing.T) {
	t.Run("with ContinueOnError", func(t *testing.T) {
		var (
			flags FlagSet
			buf   bytes.Buffer
		)

		flags.SetOutput(&buf)
		flags.Init("test", ContinueOnError)
		err := flags.Parse([]string{"--unknown"})
		require.Error(t, err)

		out := buf.String()
		require.Emptyf(t, out, "expected no output, only error")
		require.Containsf(t, err.Error(), "--unknown",
			"expected output mentioning unknown; got %q", err,
		)
	})

	t.Run("with PanicOnError", func(t *testing.T) {
		// notice the behavior inconsistent with the above test. It is what it is...
		var (
			flags FlagSet
			buf   bytes.Buffer
		)

		flags.SetOutput(&buf)
		flags.Init("test", PanicOnError)
		require.PanicsWithError(t, "unknown flag: --unknown", func() {
			_ = flags.Parse([]string{"--unknown"})
		})

		out := buf.String()
		require.Containsf(t, out, "--unknown",
			"expected output mentioning unknown; got %q", out,
		)
	})
}

func TestOutput(t *testing.T) {
	var (
		flags FlagSet
		buf   bytes.Buffer
	)

	const expect = "an example string"
	flags.SetOutput(&buf)
	fmt.Fprint(flags.Output(), expect)
	out := buf.String()
	require.Containsf(t, out, expect,
		"expected output %q; got %q", expect, out,
	)
}

// This tests that one can reset the flags. This still works but not well, and is
// superseded by FlagSet.
//
// NOTE: this does not work well with parallel testing.
func TestChangingArgs(t *testing.T) {
	ResetForTesting(func() { t.Fatal("bad parse") })
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--before", "subcmd"}
	before := Bool("before", false, "")
	require.NoError(t, GetCommandLine().Parse(os.Args[1:]))

	cmd := Arg(0)
	os.Args = []string{"subcmd", "--after", "args"}
	after := Bool("after", false, "")
	Parse()
	args := Args()

	require.True(t, *before)
	require.Equal(t, "subcmd", cmd)
	require.True(t, *after)
	require.Len(t, args, 1)
	require.Equal(t, "args", args[0])
}

// Test that -help invokes the usage message and returns ErrHelp.
func TestHelp(t *testing.T) {
	var flag bool
	mockHelp := func(called *bool) func() {
		return func() {
			*called = true
		}
	}

	t.Run("not called, regular flag invocation should work", func(t *testing.T) {
		var helpCalled bool
		fs := NewFlagSet("help test", ContinueOnError)
		fs.Usage = mockHelp(&helpCalled)

		fs.BoolVar(&flag, "flag", false, "regular flag")
		require.NoError(t, fs.Parse([]string{"--flag=true"}))
		require.Truef(t, flag, "flag was not set by --flag")
		require.Falsef(t, helpCalled, "help called for regular flag")
	})

	t.Run("called, help flag should work", func(t *testing.T) {
		var helpCalled bool
		fs := NewFlagSet("help test", ContinueOnError)
		fs.Usage = mockHelp(&helpCalled)
		err := fs.Parse([]string{"--help"})
		require.Error(t, err)
		require.ErrorIsf(t, err, ErrHelp, "expected ErrHelp; got %v", err)
		require.Truef(t, helpCalled, "help was not called")
	})

	t.Run("with help flag override", func(t *testing.T) {
		var help, helpCalled bool
		fs := NewFlagSet("help test", ContinueOnError)
		fs.Usage = mockHelp(&helpCalled)
		fs.BoolVar(&help, "help", false, "help flag")
		require.NoErrorf(t, fs.Parse([]string{"--help"}), "expected no error for defined --help")
		require.Falsef(t, helpCalled,
			"help was called unexpectedly for a user-defined help flag",
		)
	})
}

func TestNoInterspersed(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetInterspersed(false)
	f.Bool("true", true, "always true")
	f.Bool("false", false, "always false")
	require.NoError(t, f.Parse([]string{"--true", "break", "--false"}))

	args := f.Args()
	require.Len(t, args, 2)
	require.Equal(t, "break", args[0])
	require.Equal(t, "--false", args[1])
}

func TestTermination(t *testing.T) {
	f := NewFlagSet("termination", ContinueOnError)
	boolFlag := f.BoolP("bool", "l", false, "bool value")
	require.Falsef(t, f.Parsed(), "f.Parse() = true before Parse")

	const (
		arg1 = "ls"
		arg2 = "-l"
	)
	args := []string{
		"--",
		arg1,
		arg2,
	}
	f.SetOutput(ioutil.Discard)
	require.NoError(t, f.Parse(args))
	require.Truef(t, f.Parsed(), "f.Parse() = false after Parse")
	require.Falsef(t, *boolFlag, "expected boolFlag=false, got true")
	require.Lenf(t, f.Args(), 2,
		"expected 2 arguments, got %d: %v", len(f.Args()), f.Args(),
	)
	require.Equalf(t, arg1, f.Args()[0],
		"expected argument %q got %q", arg1, f.Args()[0],
	)
	require.Equalf(t, arg2, f.Args()[1],
		"expected argument %q got %q", arg2, f.Args()[0],
	)
	require.Equalf(t, 0, f.ArgsLenAtDash(),
		"expected argsLenAtDash %d got %d", 0, f.ArgsLenAtDash(),
	)
}

func TestDeprecated(t *testing.T) {
	const (
		badFlag       = "badFlag"
		usageMsg      = "use --good-flag instead"
		shortHandName = "noshorthandflag"
		shortHandMsg  = "use --noshorthandflag instead"
	)

	newFlag := func() *FlagSet {
		f := NewFlagSet("bob", ContinueOnError)
		f.Bool(badFlag, true, "always true")
		_ = f.MarkDeprecated(badFlag, usageMsg)

		return f
	}

	t.Run("with flag in doc", func(t *testing.T) {
		f := newFlag()
		require.NotContainsf(t, printFlagDefaults(f), badFlag,
			"found deprecated flag in usage!",
		)
	})

	t.Run("with unhidden flag in doc", func(t *testing.T) {
		f := newFlag()
		flg := f.Lookup(badFlag)
		require.NotNilf(t, flg,
			"unable to lookup %q in flag doc", badFlag,
		)
		flg.Hidden = false
		defaults := printFlagDefaults(f)

		require.Containsf(t, defaults, badFlag,
			"did not find deprecated flag in usage!",
		)
		require.Containsf(t, defaults, usageMsg,
			"did not find %q in defaults", usageMsg,
		)
	})

	t.Run("with shorthand in doc", func(t *testing.T) {
		f := newFlag()
		f.BoolP(shortHandName, "n", true, "always true")
		require.NoError(t,
			f.MarkShorthandDeprecated("noshorthandflag", shortHandMsg),
		)

		require.NotContainsf(t, printFlagDefaults(f), "-n,",
			"found deprecated flag shorthand in usage!",
		)
	})

	t.Run("with usage", func(t *testing.T) {
		f := newFlag()
		f.Bool("badflag", true, "always true")
		usageMsg := "use --good-flag instead"
		require.NoError(t,
			f.MarkDeprecated("badflag", usageMsg),
		)

		args := []string{"--badflag"}
		out, err := parseReturnStderr(t, f, args)
		require.NoError(t, err)

		require.Containsf(t, out, usageMsg,
			"%q not printed when using a deprecated flag!", usageMsg,
		)
	})

	t.Run("with shorthand usage", func(t *testing.T) {
		f := newFlag()
		f.BoolP(shortHandName, "n", true, "always true")
		_ = f.MarkShorthandDeprecated(shortHandName, shortHandMsg)

		args := []string{"-n"}
		out, err := parseReturnStderr(t, f, args)
		require.NoError(t, err)

		require.Containsf(t, out, shortHandMsg,
			"%q not printed when using a deprecated flag!", shortHandMsg,
		)
	})

	t.Run("with usage normalized", func(t *testing.T) {
		f := newFlag()
		f.Bool("bad-double_flag", true, "always true")
		f.SetNormalizeFunc(wordSepNormalizeFunc)
		require.NoError(t, f.MarkDeprecated("bad_double-flag", usageMsg))

		args := []string{"--bad_double_flag"}
		out, err := parseReturnStderr(t, f, args)
		require.NoError(t, err)

		require.Containsf(t, out, usageMsg,
			"%q not printed when using a deprecated flag!", usageMsg,
		)
	})
}

// Name normalization function should be called only once on flag addition
func TestMultipleNormalizeFlagNameInvocations(t *testing.T) {
	normalizeFlagNameInvocations = 0

	f := NewFlagSet("normalized", ContinueOnError)
	f.SetNormalizeFunc(wordSepNormalizeFunc)
	f.Bool("with_under_flag", false, "bool value")

	require.Equalf(t, 1, normalizeFlagNameInvocations,
		"expected normalizeFlagNameInvocations to be 1; got ", normalizeFlagNameInvocations,
	)
}

func TestHidden(t *testing.T) {
	t.Run("with doc", func(t *testing.T) {
		f := NewFlagSet("bob", ContinueOnError)
		f.Bool("secretFlag", true, "shhh")
		require.NoError(t,
			f.MarkHidden("secretFlag"),
		)

		require.NotContains(t, printFlagDefaults(f), "secretFlag",
			"found hidden flag in usage!",
		)
	})

	t.Run("with usage", func(t *testing.T) {
		f := NewFlagSet("bob", ContinueOnError)
		f.Bool("secretFlag", true, "shhh")
		require.NoError(t,
			f.MarkHidden("secretFlag"),
		)

		args := []string{"--secretFlag"}
		out, err := parseReturnStderr(t, f, args)
		require.NoError(t, err)

		require.NotContainsf(t, out, "shhh",
			"usage message printed when using a hidden flag!",
		)
	})
}

const defaultOutput = `      --A                         for bootstrapping, allow 'any' type
      --Alongflagname             disable bounds checking
  -C, --CCC                       a boolean defaulting to true (default true)
      --D path                    set relative path for local imports
  -E, --EEE num[=1234]            a num with NoOptDefVal (default 4321)
      --F number                  a non-zero number (default 2.7)
      --G float                   a float that defaults to zero
      --IP ip                     IP address with no default
      --IPMask ipMask             Netmask address with no default
      --IPNet ipNet               IP network with no default
      --Ints ints                 int slice with zero default
      --N int                     a non-zero int (default 27)
      --ND1 string[="bar"]        a string with NoOptDefVal (default "foo")
      --ND2 num[=4321]            a num with NoOptDefVal (default 1234)
      --StringArray stringArray   string array with zero default
      --StringSlice strings       string slice with zero default
      --Z int                     an int that defaults to zero
      --custom custom             custom Value implementation
      --customP custom            a VarP with default (default 10)
      --maxT timeout              set timeout for dial
  -v, --verbose count             verbosity
`

// Custom value that satisfies the Value interface.
type customValue int

func (cv *customValue) String() string { return fmt.Sprintf("%v", *cv) }

func (cv *customValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*cv = customValue(v)
	return err
}

func (cv *customValue) Type() string { return "custom" }

func TestPrintDefaults(t *testing.T) {
	fs := NewFlagSet("print defaults test", ContinueOnError)
	var buf bytes.Buffer

	fs.SetOutput(&buf)
	fs.Bool("A", false, "for bootstrapping, allow 'any' type")
	fs.Bool("Alongflagname", false, "disable bounds checking")
	fs.BoolP("CCC", "C", true, "a boolean defaulting to true")
	fs.String("D", "", "set relative `path` for local imports")
	fs.Float64("F", 2.7, "a non-zero `number`")
	fs.Float64("G", 0, "a float that defaults to zero")
	fs.Int("N", 27, "a non-zero int")
	fs.IntSlice("Ints", []int{}, "int slice with zero default")
	fs.IP("IP", nil, "IP address with no default")
	fs.IPMask("IPMask", nil, "Netmask address with no default")
	fs.IPNet("IPNet", net.IPNet{}, "IP network with no default")
	fs.Int("Z", 0, "an int that defaults to zero")
	fs.Duration("maxT", 0, "set `timeout` for dial")
	fs.String("ND1", "foo", "a string with NoOptDefVal")
	fs.Lookup("ND1").NoOptDefVal = "bar"
	fs.Int("ND2", 1234, "a `num` with NoOptDefVal")
	fs.Lookup("ND2").NoOptDefVal = "4321"
	fs.IntP("EEE", "E", 4321, "a `num` with NoOptDefVal")
	fs.ShorthandLookup("E").NoOptDefVal = "1234"
	fs.StringSlice("StringSlice", []string{}, "string slice with zero default")
	fs.StringArray("StringArray", []string{}, "string array with zero default")
	fs.CountP("verbose", "v", "verbosity")

	var cv customValue
	fs.Var(&cv, "custom", "custom Value implementation")

	cv2 := customValue(10)
	fs.VarP(&cv2, "customP", "", "a VarP with default")

	fs.PrintDefaults()
	got := buf.String()
	require.Equalf(t, defaultOutput, got,
		"got:\n%q\nwant:\n%q", got, defaultOutput,
	)
}

func TestVisitAllFlagOrder(t *testing.T) {
	fs := NewFlagSet("TestVisitAllFlagOrder", ContinueOnError)
	fs.SortFlags = false
	// https://github.com/spf13/pflag/issues/120
	fs.SetNormalizeFunc(func(f *FlagSet, name string) NormalizedName {
		return NormalizedName(name)
	})

	names := []string{"C", "B", "A", "D"}
	for _, name := range names {
		fs.Bool(name, false, "")
	}

	i := 0
	fs.VisitAll(func(f *Flag) {
		require.Equalf(t, f.Name, names[i],
			"incorrect order. Expected %v, got %v", names[i], f.Name,
		)
		i++
	})
}

func TestVisitFlagOrder(t *testing.T) {
	fs := NewFlagSet("TestVisitFlagOrder", ContinueOnError)
	fs.SortFlags = false
	names := []string{"C", "B", "A", "D"}
	for _, name := range names {
		fs.Bool(name, false, "")
		_ = fs.Set(name, "true")
	}

	i := 0
	fs.Visit(func(f *Flag) {
		require.Equalf(t, f.Name, names[i],
			"incorrect order. Expected %v, got %v", names[i], f.Name,
		)
		i++
	})
}

func parseReturnStderr(_ *testing.T, f *FlagSet, args []string) (string, error) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := f.Parse(args)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stderr = oldStderr
	out := <-outC

	return out, err
}

func printFlagDefaults(f *FlagSet) string {
	out := new(bytes.Buffer)
	f.SetOutput(out)
	f.PrintDefaults()

	return out.String()
}
