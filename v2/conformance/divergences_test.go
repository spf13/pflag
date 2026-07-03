package conformance

// This is an untagged, in-package test so it runs in the normal `go test ./...`
// (it does not need v2 to compile). It validates the divergence catalogue
// (divergences.json) and prints it under -v. The catalogue is the single source
// of truth shared with the conformance oracle (hack/oracle); the oracle uses it
// to treat exactly these test failures as expected when gating CI.

import "testing"

func TestDivergenceManifest(t *testing.T) {
	m, err := Manifest()
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("divergences.json is invalid: %v", err)
	}
	for _, d := range m.Divergences {
		t.Logf("[%s] %s\n  stdlib: %s\n  pflag : %s\n  refs  : %s\n  affected: %v",
			d.Category, d.Topic, d.Stdlib, d.Pflag, d.Refs, d.AffectedTests)
	}
}
