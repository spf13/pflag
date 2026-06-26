// Package divergence defines the catalogue of intentional differences between
// pflag v2's two conformance suites and the standard library flag package, plus
// the helpers the oracle uses to gate on them.
//
// It carries no build tag and has no pflag dependency, so the oracle can use it
// even while v2 itself does not yet compile under the "conformance" tag.
package divergence

import (
	"encoding/json"
	"fmt"
)

// Category classifies why a difference exists.
type Category string

const (
	// POSIXOverridesStdlib marks a place where pflag follows POSIX/GNU and so
	// breaks a stdlib-flag expectation. POSIX wins. Permanent.
	POSIXOverridesStdlib Category = "posix-overrides-stdlib"
	// PflagDesignDiffers marks a message/usage/behavior difference that is
	// pflag's own design choice, not mandated by POSIX. Permanent.
	PflagDesignDiffers Category = "pflag-design-differs"
	// PflagOmitsGNUFeature marks a GNU extension pflag does not implement.
	// Permanent (unless a maintainer decides to implement it).
	PflagOmitsGNUFeature Category = "pflag-omits-gnu-feature"
	// NotImplementedYet marks a test that fails only because the relevant v2 API
	// is not built yet. Temporary: it is tolerated leniently (fail/skip/never-run
	// are all fine) and should be burned down as v2 is implemented. Unlike the
	// permanent categories, a NotImplementedYet test that *passes* is not an error
	// — it just means the entry can be retired.
	NotImplementedYet Category = "not-implemented-yet"
)

// Permanent reports whether a category describes a lasting divergence (as
// opposed to NotImplementedYet, which is a temporary build-out concession).
func (c Category) Permanent() bool { return c != NotImplementedYet }

// Manifest enforcement status.
const (
	// StatusBootstrapping tolerates the conformance suite failing to build at all
	// (v2 does not implement the API yet), so CI can be green during build-out.
	StatusBootstrapping = "bootstrapping"
	// StatusEnforcing requires the suite to build and run; a build failure is red.
	StatusEnforcing = "enforcing"
)

// Divergence is one documented difference and the conformance tests it is
// expected to make fail.
type Divergence struct {
	Category      Category `json:"category"`
	Topic         string   `json:"topic"`
	Stdlib        string   `json:"stdlib"`
	Pflag         string   `json:"pflag"`
	Refs          string   `json:"refs"`
	AffectedTests []string `json:"affectedTests"`
}

// Manifest is the whole catalogue.
type Manifest struct {
	// Status gates whether a build failure is tolerated (see StatusBootstrapping).
	Status      string       `json:"status"`
	Divergences []Divergence `json:"divergences"`
}

// Parse decodes a manifest from its JSON encoding. Unknown fields (such as the
// leading "_comment") are ignored.
func Parse(data []byte) (Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parsing divergence manifest: %w", err)
	}
	return m, nil
}

// BuildRequired reports whether a build failure should be treated as red.
func (m Manifest) BuildRequired() bool { return m.Status == StatusEnforcing }

// Validate checks structural invariants of the catalogue.
func (m Manifest) Validate() error {
	switch m.Status {
	case StatusBootstrapping, StatusEnforcing:
	default:
		return fmt.Errorf("manifest status %q must be %q or %q", m.Status, StatusBootstrapping, StatusEnforcing)
	}
	seen := map[string]string{}
	for i, d := range m.Divergences {
		switch d.Category {
		case POSIXOverridesStdlib, PflagDesignDiffers, PflagOmitsGNUFeature, NotImplementedYet:
		default:
			return fmt.Errorf("divergence %d (%q): unknown category %q", i, d.Topic, d.Category)
		}
		if d.Topic == "" {
			return fmt.Errorf("divergence %d: missing topic", i)
		}
		// Permanent divergences must explain themselves; not-implemented entries
		// only need the test list (they are self-explanatory and temporary).
		if d.Category.Permanent() && (d.Pflag == "" || d.Refs == "") {
			return fmt.Errorf("divergence %d (%q): missing required field (pflag/refs)", i, d.Topic)
		}
		if len(d.AffectedTests) == 0 {
			return fmt.Errorf("divergence %d (%q): lists no affected tests", i, d.Topic)
		}
		for _, t := range d.AffectedTests {
			if prev, ok := seen[t]; ok {
				return fmt.Errorf("test %q listed under two divergences (%q and %q)", t, prev, d.Topic)
			}
			seen[t] = d.Topic
		}
	}
	return nil
}

// Expected maps each affected test name to the category that explains it.
func (m Manifest) Expected() map[string]Category {
	out := make(map[string]Category)
	for _, d := range m.Divergences {
		for _, t := range d.AffectedTests {
			out[t] = d.Category
		}
	}
	return out
}
