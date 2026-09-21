package pflag

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/go-units"
)

func TestMarshalByteSize(t *testing.T) {
	v, err := byteSizeValue(1024).MarshalFlag()
	if err != nil {
		t.Errorf("expected success, got %q", err)
	}
	expected := "1.024kB"
	if v != expected {
		t.Errorf("expected value to be %q, got %q", expected, v)
	}
}

func TestStringByteSize(t *testing.T) {
	v := byteSizeValue(2048).String()
	expected := "2.048kB"
	if v != expected {
		t.Errorf("expected value to be %q, got %q", expected, v)
	}
}

func TestUnmarshalByteSize(t *testing.T) {
	var b byteSizeValue
	err := b.UnmarshalFlag("notASize")
	if err == nil {
		t.Errorf("expected failure, got nil")
	}

	err = b.UnmarshalFlag("1MB")
	if err != nil {
		t.Errorf("expected success, got %q", err)
	}
	expected := byteSizeValue(1000000)
	if b != expected {
		t.Errorf("expected value to be %d, got %d", expected, b)
	}
}

func TestSetByteSize(t *testing.T) {
	var b byteSizeValue
	err := b.Set("notASize")
	if err == nil {
		t.Errorf("expected failure, got nil")
	}

	err = b.Set("2MB")
	if err != nil {
		t.Errorf("expected success, got %q", err)
	}
	expected := byteSizeValue(2000000)
	if b != expected {
		t.Errorf("expected value to be %d, got %d", expected, b)
	}
}

func TestTypeByteSize(t *testing.T) {
	var b byteSizeValue
	v := b.Type()
	expected := "byte-size"
	if v != expected {
		t.Errorf("expected value to be %q, got %q", expected, v)
	}
}

func setUpByteSize(value *uint64) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.ByteSizeVar(value, "size", 1*units.MiB, "Size")
	return f
}

func TestByteSize(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected uint64
	}{
		{"1KB", true, 1000},
		{"1MB", true, 1000000},
		{"1kb", true, 1000},
		{"zzz", false, 0},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		var addr uint64
		f := setUpByteSize(&addr)

		tc := &testCases[i]

		arg := fmt.Sprintf("--size=%s", tc.input)
		err := f.Parse([]string{arg})
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure")
			continue
		} else if tc.success {
			size, err := f.GetByteSize("size")
			if err != nil {
				t.Errorf("Got error trying to fetch the IP flag: %v", err)
			}
			if size != tc.expected {
				t.Errorf("for input %q, expected %d, got %d", tc.input, tc.expected, size)
			}
		}
	}
}
