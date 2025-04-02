package pflag

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func setUpTime(t *time.Time, formats []string) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.TimeVar(t, "time", time.Time{}, formats, "Time")
	return f
}

func TestTime(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		{"2022-01-01T01:01:01+00:00", true, "2022-01-01T01:01:01Z"},
		{" 2022-01-01T01:01:01+00:00", true, "2022-01-01T01:01:01Z"},
		{"2022-01-01T01:01:01+00:00 ", true, "2022-01-01T01:01:01Z"},
		{"2022-01-01T01:01:01+02:00", true, "2022-01-01T01:01:01+02:00"},
		{"2022-01-01T01:01:01.01+02:00", true, "2022-01-01T01:01:01.01+02:00"},
		{"Sat, 01 Jan 2022 01:01:01 +0000", true, "2022-01-01T01:01:01Z"},
		{"Sat, 01 Jan 2022 01:01:01 +0200", true, "2022-01-01T01:01:01+02:00"},
		{"Sat, 01 Jan 2022 01:01:01 +0000", true, "2022-01-01T01:01:01Z"},
		{"", false, ""},
		{"not a date", false, ""},
		{"2022-01-01 01:01:01", false, ""},
		{"2022-01-01T01:01:01", false, ""},
		{"01 Jan 2022 01:01:01 +0000", false, ""},
		{"Sat, 01 Jan 2022 01:01:01", false, ""},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		var timeVar time.Time
		formats := []string{time.RFC3339Nano, time.RFC1123Z}
		f := setUpTime(&timeVar, formats)

		tc := &testCases[i]

		arg := fmt.Sprintf("--time=%s", tc.input)
		err := f.Parse([]string{arg})
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure")
			continue
		} else if tc.success {
			timeResult, err := f.GetTime("time")
			if err != nil {
				t.Errorf("Got error trying to fetch the Time flag: %v", err)
			}
			if timeResult.Format(time.RFC3339Nano) != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, timeVar.Format(time.RFC3339Nano))
			}
		}
	}
}
