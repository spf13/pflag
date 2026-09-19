// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestSliceFlagsAcceptExplicitEmptyValue(t *testing.T) {
	_, defaultNetwork, err := net.ParseCIDR("192.0.2.0/24")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		add  func(*FlagSet) interface{}
	}{
		{"bool", func(f *FlagSet) interface{} { return f.BoolSlice("values", []bool{true}, "") }},
		{"duration", func(f *FlagSet) interface{} { return f.DurationSlice("values", []time.Duration{time.Second}, "") }},
		{"float32", func(f *FlagSet) interface{} { return f.Float32Slice("values", []float32{1}, "") }},
		{"float64", func(f *FlagSet) interface{} { return f.Float64Slice("values", []float64{1}, "") }},
		{"int", func(f *FlagSet) interface{} { return f.IntSlice("values", []int{1}, "") }},
		{"int32", func(f *FlagSet) interface{} { return f.Int32Slice("values", []int32{1}, "") }},
		{"int64", func(f *FlagSet) interface{} { return f.Int64Slice("values", []int64{1}, "") }},
		{"ip", func(f *FlagSet) interface{} { return f.IPSlice("values", []net.IP{net.ParseIP("192.0.2.1")}, "") }},
		{"ipnet", func(f *FlagSet) interface{} {
			return f.IPNetSlice("values", []net.IPNet{*defaultNetwork}, "")
		}},
		{"string", func(f *FlagSet) interface{} { return f.StringSlice("values", []string{"value"}, "") }},
		{"uint", func(f *FlagSet) interface{} { return f.UintSlice("values", []uint{1}, "") }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			value := reflect.ValueOf(test.add(f)).Elem()
			if err := f.Parse([]string{"--values="}); err != nil {
				t.Fatalf("expected no error; got %v", err)
			}
			if !f.Changed("values") {
				t.Fatal("expected --values= to mark the flag as changed")
			}
			if value.IsNil() || value.Len() != 0 {
				t.Fatalf("got %v, want a non-nil empty slice", value.Interface())
			}
		})
	}
}

func TestSliceDefaultPreservedWithoutFlag(t *testing.T) {
	tests := []struct {
		name         string
		defaultValue []int
	}{
		{"nil", nil},
		{"empty", []int{}},
		{"populated", []int{1, 2}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			value := f.IntSlice("values", test.defaultValue, "")
			if err := f.Parse(nil); err != nil {
				t.Fatalf("expected no error; got %v", err)
			}
			if !reflect.DeepEqual(*value, test.defaultValue) {
				t.Fatalf("got %v, want preserved default %v", *value, test.defaultValue)
			}
		})
	}
}
