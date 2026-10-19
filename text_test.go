package pflag

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type textString string

func (s *textString) UnmarshalText(text []byte) error {
	*s = textString(text)
	return nil
}

func TestTextType(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	var value time.Time
	f.TextVar(&value, "time", time.Time{}, "time stamp")
	if got := f.Lookup("time").Value.Type(); got != "*time.Time" {
		t.Errorf("expected type *time.Time, got %q", got)
	}
}

func TestGetTextRejectsDifferentType(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	var value time.Time
	f.TextVar(&value, "time", time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC), "time stamp")
	output := textString("unchanged")
	err := f.GetText("time", &output)
	if err == nil || !strings.Contains(err.Error(), "trying to get *pflag.textString value of flag of type *time.Time") {
		t.Errorf("expected type mismatch error, got %v", err)
	}
	if output != "unchanged" {
		t.Errorf("output was modified after type mismatch: %q", output)
	}
}

func setUpTime(t *time.Time) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.TextVar(t, "time", time.Now(), "time stamp")
	return f
}

func TestText(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected time.Time
	}{
		{"2003-01-02T15:04:05Z", true, time.Date(2003, 1, 2, 15, 04, 05, 0, time.UTC)},
		{"2003-01-02 15:05:01", false, time.Time{}}, // negative case, invalid layout
		{"2024-11-22T03:01:02Z", true, time.Date(2024, 11, 22, 3, 1, 02, 0, time.UTC)},
		{"2006-01-02T15:04:05+07:00", true, time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("UTC+7", 7*60*60))},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		var ts time.Time
		f := setUpTime(&ts)
		tc := &testCases[i]
		arg := fmt.Sprintf("--time=%s", tc.input)
		err := f.Parse([]string{arg})
		if err != nil {
			if tc.success {
				t.Errorf("expected parsing to succeed, but got %q", err)
			}
			continue
		}
		if !tc.success {
			t.Errorf("expected parsing failure, but parsing succeeded")
			continue
		}
		parsedT := new(time.Time)
		err = f.GetText("time", parsedT)
		if err != nil {
			t.Errorf("Got error trying to fetch the time flag: %v", err)
		}
		if !parsedT.Equal(tc.expected) {
			t.Errorf("expected %q, got %q", tc.expected, parsedT)
		}

	}
}
