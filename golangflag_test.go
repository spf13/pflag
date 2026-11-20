// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	goflag "flag"
	"testing"
	"time"
)

func TestGoflags(t *testing.T) {
	goflag.String("stringFlag", "stringFlag", "stringFlag")
	goflag.Bool("boolFlag", false, "boolFlag")
	var testxxxValue string
	goflag.StringVar(&testxxxValue, "test.xxx", "test.xxx", "it is a test flag")

	f := NewFlagSet("test", ContinueOnError)

	f.AddGoFlagSet(goflag.CommandLine)
	args := []string{"--stringFlag=bob", "--boolFlag", "-test.xxx=testvalue"}
	err := f.Parse(args)
	if err != nil {
		t.Fatal("expected no error; get", err)
	}

	getString, err := f.GetString("stringFlag")
	if err != nil {
		t.Fatal("expected no error; get", err)
	}
	if getString != "bob" {
		t.Fatalf("expected getString=bob but got getString=%s", getString)
	}

	getBool, err := f.GetBool("boolFlag")
	if err != nil {
		t.Fatal("expected no error; get", err)
	}
	if getBool != true {
		t.Fatalf("expected getBool=true but got getBool=%v", getBool)
	}
	if !f.Parsed() {
		t.Fatal("f.Parsed() return false after f.Parse() called")
	}

	if testxxxValue != "test.xxx" {
		t.Fatalf("expected testxxxValue to be test.xxx but got %v", testxxxValue)
	}
	err = ParseSkippedFlags(args, goflag.CommandLine)
	if err != nil {
		t.Fatal("expected no error; ParseSkippedFlags", err)
	}
	if testxxxValue != "testvalue" {
		t.Fatalf("expected testxxxValue to be testvalue but got %v", testxxxValue)
	}

	// in fact it is useless. because `go test` called flag.Parse()
	if !goflag.CommandLine.Parsed() {
		t.Fatal("goflag.CommandLine.Parsed() return false after f.Parse() called")
	}
}

func TestToGoflags(t *testing.T) {
	pfs := FlagSet{}
	gfs := goflag.FlagSet{}
	pfs.String("StringFlag", "String value", "String flag usage")
	pfs.Int("IntFlag", 1, "Int flag usage")
	pfs.Uint("UintFlag", 2, "Uint flag usage")
	pfs.Int64("Int64Flag", 3, "Int64 flag usage")
	pfs.Uint64("Uint64Flag", 4, "Uint64 flag usage")
	pfs.Int8("Int8Flag", 5, "Int8 flag usage")
	pfs.Float64("Float64Flag", 6.0, "Float64 flag usage")
	pfs.Duration("DurationFlag", time.Second, "Duration flag usage")
	pfs.Bool("BoolFlag", true, "Bool flag usage")
	pfs.String("deprecated", "Deprecated value", "Deprecated flag usage")
	_ = pfs.MarkDeprecated("deprecated", "obsolete")

	pfs.CopyToGoFlagSet(&gfs)

	// Modify via pfs. Should be visible via gfs because both share the
	// same values.
	for name, value := range map[string]string{
		"StringFlag":  "Modified String value",
		"IntFlag":     "11",
		"UintFlag":    "12",
		"Int64Flag":   "13",
		"Uint64Flag":  "14",
		"Int8Flag":    "15",
		"Float64Flag": "16.0",
		"BoolFlag":    "false",
	} {
		pf := pfs.Lookup(name)
		if pf == nil {
			t.Errorf("%s: not found in pflag flag set", name)
			continue
		}
		if err := pf.Value.Set(value); err != nil {
			t.Errorf("error setting %s = %s: %v", name, value, err)
		}
	}

	// Check that all flags were added and share the same value.
	pfs.VisitAll(func(pf *Flag) {
		gf := gfs.Lookup(pf.Name)
		if gf == nil {
			t.Errorf("%s: not found in Go flag set", pf.Name)
			return
		}
		if gf.Value.String() != pf.Value.String() {
			t.Errorf("%s: expected value %v from Go flag set, got %v",
				pf.Name, pf.Value, gf.Value)
			return
		}
	})

	// Check for unexpected additional flags.
	gfs.VisitAll(func(gf *goflag.Flag) {
		pf := pfs.Lookup(gf.Name)
		if pf == nil {
			t.Errorf("%s: not found in pflag flag set", gf.Name)
			return
		}
	})

	deprecated := gfs.Lookup("deprecated")
	if deprecated == nil {
		t.Error("deprecated: not found in Go flag set")
	} else {
		expectedUsage := "Deprecated flag usage (DEPRECATED: obsolete)"
		if deprecated.Usage != expectedUsage {
			t.Errorf("deprecation remark not added, expected usage %q, got %q", expectedUsage, deprecated.Usage)
		}
	}
}

func TestCopyToGoFlagSetPreservesCurrentStringValue(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{name: "default", want: "default"},
		{name: "parsed", args: []string{"--name=parsed"}, want: "parsed"},
		{name: "empty", args: []string{"--name="}, want: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pfs := NewFlagSet("source", ContinueOnError)
			value := pfs.String("name", "default", "a string flag")
			if err := pfs.Parse(tc.args); err != nil {
				t.Fatal(err)
			}
			gfs := goflag.NewFlagSet("destination", goflag.ContinueOnError)
			pfs.CopyToGoFlagSet(gfs)
			if *value != tc.want {
				t.Fatalf("copy changed source value: got %q, want %q", *value, tc.want)
			}
			if got := gfs.Lookup("name").Value.String(); got != tc.want {
				t.Errorf("copied value = %q, want %q", got, tc.want)
			}
			if got := pfs.Lookup("name").DefValue; got != "default" {
				t.Errorf("source default = %q, want %q", got, "default")
			}
			if got := gfs.Lookup("name").DefValue; got != tc.want {
				t.Errorf("copied default = %q, want current value %q", got, tc.want)
			}
			if err := gfs.Set("name", "from-go"); err != nil {
				t.Fatal(err)
			}
			if *value != "from-go" {
				t.Errorf("Go flag update not shared with source: got %q", *value)
			}
			if err := pfs.Set("name", "from-pflag"); err != nil {
				t.Fatal(err)
			}
			if got := gfs.Lookup("name").Value.String(); got != "from-pflag" {
				t.Errorf("pflag update not shared with copy: got %q", got)
			}
		})
	}
}

func TestParseSkippedFlagsSeparateValues(t *testing.T) {
	for _, tc := range []struct {
		name  string
		args  []string
		run   string
		count int
		v     bool
	}{
		{
			name: "separate string value",
			args: []string{"-test.run", "TestExample"},
			run:  "TestExample",
		},
		{
			name:  "separate integer value",
			args:  []string{"-test.count", "2"},
			count: 2,
		},
		{
			name:  "negative integer value",
			args:  []string{"-test.count", "-2"},
			count: -2,
		},
		{
			name: "repeated separate values",
			args: []string{"-test.run", "first", "-test.run", "second"},
			run:  "second",
		},
		{
			name: "value is argument terminator",
			args: []string{"-test.run", "--", "-test.v"},
			run:  "--",
			v:    true,
		},
		{
			name: "empty value",
			args: []string{"-test.run", ""},
		},
		{
			name: "value looks like a test flag",
			args: []string{"-test.run", "-test.v"},
			run:  "-test.v",
		},
		{
			name: "value looks like another flag",
			args: []string{"-test.run", "--custom"},
			run:  "--custom",
		},
		{
			name:  "mixed flag forms",
			args:  []string{"--custom=value", "-test.run", "TestExample", "-test.v", "-test.count=2"},
			run:   "TestExample",
			count: 2,
			v:     true,
		},
		{
			name: "boolean does not consume next flag",
			args: []string{"-test.v", "-test.run", "TestExample"},
			run:  "TestExample",
			v:    true,
		},
		{
			name: "equals value does not consume next flag",
			args: []string{"-test.run=TestExample", "-test.v"},
			run:  "TestExample",
			v:    true,
		},
		{
			name: "explicit boolean value",
			args: []string{"-test.v=false", "-test.run", "TestExample"},
			run:  "TestExample",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
			run := fs.String("test.run", "", "")
			count := fs.Int("test.count", 0, "")
			verbose := fs.Bool("test.v", false, "")
			if err := ParseSkippedFlags(tc.args, fs); err != nil {
				t.Fatal(err)
			}
			if *run != tc.run || *count != tc.count || *verbose != tc.v {
				t.Errorf("got run=%q, count=%d, verbose=%v; want %q, %d, %v", *run, *count, *verbose, tc.run, tc.count, tc.v)
			}
		})
	}
}

func TestParseSkippedFlagsMissingOrUnknownFlag(t *testing.T) {
	for _, arg := range []string{"-test.run", "-test.unknown"} {
		t.Run(arg, func(t *testing.T) {
			fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
			fs.String("test.run", "", "")
			if err := ParseSkippedFlags([]string{arg}, fs); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

type skippedFlagValue struct {
	value   string
	boolean bool
}

func (v *skippedFlagValue) String() string         { return v.value }
func (v *skippedFlagValue) Set(value string) error { v.value = value; return nil }
func (v *skippedFlagValue) IsBoolFlag() bool       { return v.boolean }

func TestParseSkippedFlagsCustomBoolFlag(t *testing.T) {
	for _, boolean := range []bool{false, true} {
		fs := goflag.NewFlagSet("test", goflag.ContinueOnError)
		value := &skippedFlagValue{boolean: boolean}
		fs.Var(value, "test.custom", "")
		args := []string{"-test.custom", "separate"}
		want := "separate"
		if boolean {
			want = "true"
		}
		if err := ParseSkippedFlags(args, fs); err != nil {
			t.Fatal(err)
		}
		if value.value != want {
			t.Errorf("IsBoolFlag=%v: got %q, want %q", boolean, value.value, want)
		}
	}
}
