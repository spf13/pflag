// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestFlagSet_HasFlags_And_HasAvailableFlags(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	if f.HasFlags() {
		t.Errorf("expected HasFlags to be false on empty FlagSet")
	}
	if f.HasAvailableFlags() {
		t.Errorf("expected HasAvailableFlags to be false on empty FlagSet")
	}

	f.String("hidden", "default", "hidden flag")
	if err := f.MarkHidden("hidden"); err != nil {
		t.Fatalf("unexpected error marking hidden: %v", err)
	}

	if !f.HasFlags() {
		t.Errorf("expected HasFlags to be true after adding hidden flag")
	}
	if f.HasAvailableFlags() {
		t.Errorf("expected HasAvailableFlags to be false when only hidden flag exists")
	}

	f.String("visible", "default", "visible flag")
	if !f.HasAvailableFlags() {
		t.Errorf("expected HasAvailableFlags to be true after adding visible flag")
	}
}

func TestFlagSet_MarkDeprecated_Errors(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	if err := f.MarkDeprecated("nonexistent", "deprecated message"); err == nil {
		t.Errorf("expected error marking non-existent flag as deprecated")
	}
	if err := f.MarkShorthandDeprecated("nonexistent", "deprecated shorthand"); err == nil {
		t.Errorf("expected error marking non-existent shorthand as deprecated")
	}
	if err := f.MarkHidden("nonexistent"); err == nil {
		t.Errorf("expected error marking non-existent flag as hidden")
	}
}

func TestFlagSet_AddFlag_Errors(t *testing.T) {
	// Test duplicate flag panic
	assertPanic := func(name string, f func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s: expected panic, but did not panic", name)
			}
		}()
		f()
	}

	assertPanic("duplicate flag", func() {
		f := NewFlagSet("test", ContinueOnError)
		f.SetOutput(io.Discard)
		f.String("testflag", "val", "usage")
		fl := &Flag{
			Name:  "testflag",
			Usage: "duplicate",
			Value: newStringValue("val2", new(string)),
		}
		f.AddFlag(fl)
	})

	assertPanic("multichar shorthand", func() {
		f := NewFlagSet("test", ContinueOnError)
		f.SetOutput(io.Discard)
		fl := &Flag{
			Name:      "testflag2",
			Shorthand: "ab",
			Usage:     "usage",
			Value:     newStringValue("val", new(string)),
		}
		f.AddFlag(fl)
	})

	assertPanic("duplicate shorthand", func() {
		f := NewFlagSet("test", ContinueOnError)
		f.SetOutput(io.Discard)
		f.StringVarP(new(string), "flagA", "a", "val", "usage")
		f.StringVarP(new(string), "flagB", "a", "val", "usage")
	})
}

func TestFlagSet_AddFlagSet_Duplicate(t *testing.T) {
	f1 := NewFlagSet("f1", ContinueOnError)
	f1.String("shared", "val1", "desc1")
	f1.AddFlagSet(nil) // should safely return

	f2 := NewFlagSet("f2", ContinueOnError)
	f2.String("shared", "val2", "desc2")
	f2.String("unique", "val3", "desc3")

	f1.AddFlagSet(f2) // shared is ignored, unique is added
	if f1.Lookup("unique") == nil {
		t.Errorf("expected unique flag to be added from f2")
	}
}

func TestFlagSet_Lookup_And_Set_Errors(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.String("foo", "bar", "usage")

	if flag := f.Lookup("unknown"); flag != nil {
		t.Errorf("expected nil for unknown flag lookup")
	}
	if flag := f.ShorthandLookup("u"); flag != nil {
		t.Errorf("expected nil for unknown shorthand lookup")
	}

	if err := f.Set("unknown", "value"); err == nil {
		t.Errorf("expected error setting unknown flag")
	}

	// getFlagType on non-existent flag
	_, err := f.getFlagType("unknown", "string", func(s string) (interface{}, error) { return s, nil })
	if err == nil {
		t.Errorf("expected error for getFlagType on non-existent flag")
	}

	// getFlagType with type mismatch
	f.Int("intflag", 10, "int flag")
	_, err = f.getFlagType("intflag", "string", func(s string) (interface{}, error) { return s, nil })
	if err == nil {
		t.Errorf("expected error for getFlagType with type mismatch")
	}

	// getFlagType with conversion error
	_, err = f.getFlagType("foo", "int", func(s string) (interface{}, error) {
		return nil, errors.New("parse error")
	})
	if err == nil {
		t.Errorf("expected error for getFlagType when conversion fails")
	}
}

func TestFlagSet_UnquoteUsage_Variations(t *testing.T) {
	dummyVal := newStringValue("val", new(string))
	tests := []struct {
		usage     string
		wantName  string
		wantUsage string
	}{
		{"`filepath` path to file", "filepath", "filepath path to file"},
		{"unquoted usage", "string", "unquoted usage"},
		{"`invalid backtick without closing", "string", "`invalid backtick without closing"},
		{"`name`", "name", "name"},
	}

	for _, tc := range tests {
		name, usage := UnquoteUsage(&Flag{Usage: tc.usage, Value: dummyVal})
		if name != tc.wantName || usage != tc.wantUsage {
			t.Errorf("UnquoteUsage(%q) = (%q, %q); want (%q, %q)", tc.usage, name, usage, tc.wantName, tc.wantUsage)
		}
	}
}

func TestFlagSet_Wrap_And_Usage_EdgeCases(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	buf := new(bytes.Buffer)
	f.SetOutput(buf)

	f.StringP("longflag", "l", "default", "This is a very long usage description that should wrap around multiple lines properly when printed to the output stream.")
	f.StringP("second", "s", "val", "`type` short usage")
	f.PrintDefaults()

	if buf.Len() == 0 {
		t.Errorf("expected non-empty output from PrintDefaults")
	}

	// FlagUsages
	usages := f.FlagUsages()
	if len(usages) == 0 {
		t.Errorf("expected non-empty FlagUsages")
	}
}

func TestFlagSet_ParseAll_And_Interspersed(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.SetInterspersed(false)

	flagA := f.BoolP("alpha", "a", false, "alpha flag")
	err := f.Parse([]string{"arg1", "-a", "arg2"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if *flagA {
		t.Errorf("expected flagA to be false when interspersed is false and argument comes first")
	}
	if len(f.Args()) != 3 {
		t.Errorf("expected 3 non-flag args with interspersed=false, got %v", f.Args())
	}

	// ParseAll with custom fn
	f2 := NewFlagSet("test2", ContinueOnError)
	f2.String("valid", "default", "usage")
	var parsedFlags []string
	err = f2.ParseAll([]string{"--valid=customval"}, func(flag *Flag, value string) error {
		parsedFlags = append(parsedFlags, flag.Name)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected ParseAll error: %v", err)
	}
	if len(parsedFlags) != 1 || parsedFlags[0] != "valid" {
		t.Errorf("expected parsedFlags = [valid], got %v", parsedFlags)
	}

	// ParseAll returning error
	f3 := NewFlagSet("test3", ContinueOnError)
	f3.SetOutput(io.Discard)
	f3.String("errflag", "default", "usage")
	err = f3.ParseAll([]string{"--errflag=val"}, func(flag *Flag, value string) error {
		return errors.New("custom error")
	})
	if err == nil {
		t.Errorf("expected error from ParseAll when fn returns error")
	}
}

func TestPackageLevel_CommandLine_Helpers(t *testing.T) {
	// Reset CommandLine for testing package-level functions
	oldCommandLine := CommandLine
	defer func() { CommandLine = oldCommandLine }()

	CommandLine = NewFlagSet("test_cmd", ContinueOnError)
	CommandLine.SetOutput(io.Discard)

	if CommandLine.HasFlags() {
		t.Errorf("expected CommandLine.HasFlags() to be false on empty CommandLine")
	}
	if CommandLine.HasAvailableFlags() {
		t.Errorf("expected CommandLine.HasAvailableFlags() to be false on empty CommandLine")
	}

	var dummy string
	StringVar(&dummy, "dummy", "default", "dummy usage")
	StringVarP(&dummy, "dummy2", "d", "default", "dummy usage")

	if !CommandLine.HasFlags() {
		t.Errorf("expected CommandLine.HasFlags() to be true")
	}
	if !CommandLine.HasAvailableFlags() {
		t.Errorf("expected CommandLine.HasAvailableFlags() to be true")
	}

	if f := Lookup("dummy"); f == nil {
		t.Errorf("Lookup('dummy') returned nil")
	}
	if f := ShorthandLookup("d"); f == nil {
		t.Errorf("ShorthandLookup('d') returned nil")
	}

	SetInterspersed(true)

	if Parsed() {
		t.Errorf("expected Parsed() to be false before parse")
	}

	// Var and VarPF on CommandLine
	var customVal string
	Var(newStringValue("init", &customVal), "custom", "custom usage")
	VarP(newStringValue("init2", &customVal), "custom2", "c", "custom usage")

	err := CommandLine.Parse([]string{"--dummy=foo", "pos1", "pos2"})
	if err != nil {
		t.Fatalf("CommandLine.Parse failed: %v", err)
	}

	if !Parsed() {
		t.Errorf("expected Parsed() to be true after parse")
	}

	if NFlag() != 1 {
		t.Errorf("NFlag() = %d, want 1", NFlag())
	}
	if NArg() != 2 {
		t.Errorf("NArg() = %d, want 2", NArg())
	}
	if len(Args()) != 2 {
		t.Errorf("Args() len = %d, want 2", len(Args()))
	}
	if Arg(0) != "pos1" || Arg(1) != "pos2" {
		t.Errorf("Arg(0)=%q, Arg(1)=%q, want pos1, pos2", Arg(0), Arg(1))
	}
	if Arg(99) != "" {
		t.Errorf("Arg(99) out of bounds want empty, got %q", Arg(99))
	}

	// PrintDefaults & defaultUsage
	buf := new(bytes.Buffer)
	CommandLine.SetOutput(buf)
	PrintDefaults()
	if buf.Len() == 0 {
		t.Errorf("PrintDefaults output empty")
	}

	// ParseAll on CommandLine
	CommandLine = NewFlagSet("test_cmd2", ContinueOnError)
	CommandLine.SetOutput(io.Discard)
	var flagVal string
	StringVar(&flagVal, "target", "default", "usage")
	var called bool
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--target=custom"}
	ParseAll(func(flag *Flag, value string) error {
		if flag.Name == "target" && value == "custom" {
			called = true
		}
		return nil
	})
	if !called {
		t.Errorf("ParseAll callback not called")
	}
}

func TestFlagSet_ArgsLenAtDash(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.String("flag", "val", "usage")

	err := f.Parse([]string{"--flag=foo", "arg1", "--", "arg2", "arg3"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if f.ArgsLenAtDash() != 1 {
		t.Errorf("ArgsLenAtDash() = %d, want 1", f.ArgsLenAtDash())
	}
	if len(f.Args()) != 3 {
		t.Errorf("Args() = %v, want 3 args", f.Args())
	}

	// Output & SetOutput
	var b bytes.Buffer
	f.SetOutput(&b)
	if f.Output() != &b {
		t.Errorf("Output() did not return buffer")
	}
	if f.Name() != "test" {
		t.Errorf("Name() = %q, want test", f.Name())
	}
}
