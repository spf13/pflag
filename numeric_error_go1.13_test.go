//go:build go1.13
// +build go1.13

package pflag

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestNumericParseErrors(t *testing.T) {
	f := NewFlagSet("numbers", ContinueOnError)
	f.Int("int", 0, "")
	f.Int8("int8", 0, "")
	f.Int16("int16", 0, "")
	f.Int32("int32", 0, "")
	f.Int64("int64", 0, "")
	f.Uint("uint", 0, "")
	f.Uint8("uint8", 0, "")
	f.Uint16("uint16", 0, "")
	f.Uint32("uint32", 0, "")
	f.Uint64("uint64", 0, "")
	f.Float32("float32", 0, "")
	f.Float64("float64", 0, "")

	for _, tc := range []struct{ name, overflow string }{
		{"int", "9223372036854775808"}, {"int8", "-129"}, {"int16", "32768"},
		{"int32", "2147483648"}, {"int64", "9223372036854775808"},
		{"uint", "18446744073709551616"}, {"uint8", "256"}, {"uint16", "65536"},
		{"uint32", "4294967296"}, {"uint64", "18446744073709551616"},
		{"float32", "1e39"}, {"float64", "1e309"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, input := range []struct {
				value string
				cause error
			}{{tc.overflow, strconv.ErrRange}, {"abc", strconv.ErrSyntax}} {
				err := f.Parse([]string{"--" + tc.name + "=" + input.value})
				var numErr *strconv.NumError
				if !errors.As(err, &numErr) || numErr.Num != input.value || !errors.Is(err, input.cause) {
					t.Errorf("Parse(%q) = %v, want a NumError wrapping %v", input.value, err, input.cause)
				}
				if errors.Is(input.cause, strconv.ErrRange) && (err == nil || !strings.HasSuffix(err.Error(), ": value out of range")) {
					t.Errorf("Parse(%q) = %v, want a range diagnostic", input.value, err)
				}
			}
		})
	}
}
