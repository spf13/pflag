package main

import (
	"testing"

	"github.com/spf13/pflag/v2/conformance/internal/divergence"
)

func ev(action, test string) event { return event{Action: action, Test: test} }

func TestEvaluate(t *testing.T) {
	// Documented catalogue used by every case: TestParse is a permanent
	// divergence, TestNewThing is a temporary not-implemented-yet entry.
	expected := map[string]divergence.Category{
		"TestParse":    divergence.POSIXOverridesStdlib,
		"TestNewThing": divergence.NotImplementedYet,
	}

	tests := []struct {
		name            string
		events          []event
		buildRequired   bool
		wantOK          bool
		wantRegressions []string
		wantResolved    []string
		wantUnknown     []string
		wantRetirable   []string
	}{
		{
			name:          "build failure while bootstrapping is tolerated",
			events:        []event{{Action: "fail"}},
			buildRequired: false,
			wantOK:        true,
		},
		{
			name:          "build failure while enforcing is red",
			events:        []event{{Action: "fail"}},
			buildRequired: true,
			wantOK:        false,
		},
		{
			name: "perfect: permanent fails, not-implemented fails, others pass",
			events: []event{
				ev("fail", "TestParse"),
				ev("fail", "TestNewThing"),
				ev("pass", "TestShortOptions"),
			},
			buildRequired: true,
			wantOK:        true,
		},
		{
			name: "not-implemented test that passes is retirable but still green",
			events: []event{
				ev("fail", "TestParse"),
				ev("pass", "TestNewThing"),
				ev("pass", "TestShortOptions"),
			},
			buildRequired: true,
			wantOK:        true,
			wantRetirable: []string{"TestNewThing"},
		},
		{
			name: "not-implemented test that never ran is fine",
			events: []event{
				ev("fail", "TestParse"),
				ev("pass", "TestShortOptions"),
			},
			buildRequired: true,
			wantOK:        true,
		},
		{
			name: "regression: undocumented failure",
			events: []event{
				ev("fail", "TestParse"),
				ev("fail", "TestShortOptions"),
			},
			buildRequired:   true,
			wantOK:          false,
			wantRegressions: []string{"TestShortOptions"},
		},
		{
			name: "resolved permanent divergence is red",
			events: []event{
				ev("pass", "TestParse"),
			},
			buildRequired: true,
			wantOK:        false,
			wantResolved:  []string{"TestParse"},
		},
		{
			name: "permanent divergence that never ran is red",
			events: []event{
				ev("pass", "TestShortOptions"),
			},
			buildRequired: true,
			wantOK:        false,
			wantUnknown:   []string{"TestParse"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := evaluate(tc.events, expected, tc.buildRequired)
			if got.OK != tc.wantOK {
				t.Errorf("OK = %v, want %v\nreport:\n%s", got.OK, tc.wantOK, got.String())
			}
			if !equalStrings(got.Regressions, tc.wantRegressions) {
				t.Errorf("Regressions = %v, want %v", got.Regressions, tc.wantRegressions)
			}
			if !equalStrings(got.Resolved, tc.wantResolved) {
				t.Errorf("Resolved = %v, want %v", got.Resolved, tc.wantResolved)
			}
			if !equalStrings(got.UnknownExpected, tc.wantUnknown) {
				t.Errorf("UnknownExpected = %v, want %v", got.UnknownExpected, tc.wantUnknown)
			}
			if !equalStrings(got.Retirable, tc.wantRetirable) {
				t.Errorf("Retirable = %v, want %v", got.Retirable, tc.wantRetirable)
			}
		})
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
