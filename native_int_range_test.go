package pflag

import (
	"strconv"
	"testing"
)

func TestNativeIntegerFlagRange(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	t.Run("int", func(t *testing.T) {
		f := NewFlagSet("int-range", ContinueOnError)
		value := f.Int("value", 7, "")
		for _, valid := range []int{minInt, maxInt} {
			if err := f.Set("value", strconv.FormatInt(int64(valid), 10)); err != nil {
				t.Fatalf("valid native int rejected: %v", err)
			}
			if *value != valid {
				t.Fatalf("got %d, want %d", *value, valid)
			}
		}
		for _, invalid := range []string{
			strconv.FormatUint(uint64(maxInt)+1, 10),
			"-" + strconv.FormatUint(uint64(maxInt)+2, 10),
		} {
			if err := f.Set("value", invalid); err == nil {
				t.Errorf("accepted out-of-range int %q as %d", invalid, *value)
			}
			if *value != maxInt {
				t.Errorf("invalid input changed value: got %d, want %d", *value, maxInt)
			}
		}
	})
	t.Run("uint", func(t *testing.T) {
		f := NewFlagSet("uint-range", ContinueOnError)
		maxUint := ^uint(0)
		value := f.Uint("value", 7, "")
		if err := f.Set("value", strconv.FormatUint(uint64(maxUint), 10)); err != nil {
			t.Fatalf("valid native uint rejected: %v", err)
		}
		if *value != maxUint {
			t.Fatalf("got %d, want %d", *value, maxUint)
		}
		invalid := "18446744073709551616"
		if strconv.IntSize == 32 {
			invalid = "4294967296"
		}
		if err := f.Set("value", invalid); err == nil {
			t.Errorf("accepted out-of-range uint %q as %d", invalid, *value)
		}
		if *value != maxUint {
			t.Errorf("invalid input changed value: got %d, want %d", *value, maxUint)
		}
	})
}
