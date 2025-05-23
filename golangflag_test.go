// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	goflag "flag"
	"testing"
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
