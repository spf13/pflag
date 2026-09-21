package pflag

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

var (
	flagName = "uuid"
)

func setUpUUID(u *uuid.UUID) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.UUIDVar(u, flagName, uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"), "Common UUID")
	return f
}

func TestUUID(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		{"00000000-0000-0000-0000-000000000000", true, "00000000-0000-0000-0000-000000000000"},
		{"6ba7b810-9dad-11d1-80b4-00c04fd430c8", true, "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
		{"deadbeef-cafe-c0de-f00d-badfacecab42", true, "deadbeef-cafe-c0de-f00d-badfacecab42"},
		{" 00000000-0000-0000-0000-000000000000 ", true, "00000000-0000-0000-0000-000000000000"},
		{"00000000-0000-0000-0000-00000000000t", false, ""},
		{"00000000-0000-0000-0000-0000000000000", false, ""},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		var u uuid.UUID
		f := setUpUUID(&u)

		tc := &testCases[i]

		arg := fmt.Sprintf("--%s=%s", flagName, tc.input)
		err := f.Parse([]string{arg})
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure")
			continue
		} else if tc.success {
			u, err := f.GetUUID(flagName)
			if err != nil {
				t.Errorf("Got error trying to fetch the IP flag: %v", err)
			}
			if u.String() != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, u.String())
			}
		}
	}
}
