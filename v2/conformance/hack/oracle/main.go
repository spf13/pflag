// Command oracle reads `go test -json` output for the conformance suites on
// stdin and decides whether pflag v2 conforms, treating the documented
// divergences (conformance/divergences.json) as expected outcomes.
//
// When the catalogue status is "enforcing", it passes (exit 0) iff the suite
// built and ran, and:
//
//   - every failing top-level test is documented (a failure that is NOT in the
//     catalogue is a regression -> fail);
//   - every PERMANENT divergence test actually failed/skipped (one that passed is
//     a resolved divergence to retire -> fail; one that never ran names an
//     unknown test -> fail);
//   - "not-implemented-yet" tests are tolerated however they end up; if one now
//     passes the oracle nudges you to retire it but stays green.
//
// When the status is "bootstrapping", a build failure is tolerated (green) so CI
// can pass while v2 has no API yet.
//
// Usage:
//
//	go test -tags conformance -json ./conformance | go run ./conformance/hack/oracle
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/pflag/v2/conformance"
	"github.com/spf13/pflag/v2/conformance/internal/divergence"
)

// event is the subset of `go test -json` records we care about.
type event struct {
	Action string `json:"Action"`
	Test   string `json:"Test"`
}

func main() {
	m, err := conformance.Manifest()
	if err != nil {
		fmt.Fprintln(os.Stderr, "oracle: loading divergence catalogue:", err)
		os.Exit(2)
	}
	if err := m.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "oracle: invalid divergence catalogue:", err)
		os.Exit(2)
	}
	events, err := readEvents(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "oracle: reading go test -json:", err)
		os.Exit(2)
	}
	rep := evaluate(events, m.Expected(), m.BuildRequired())
	fmt.Print(rep.String())
	if !rep.OK {
		os.Exit(1)
	}
}

func readEvents(r io.Reader) ([]event, error) {
	var evs []event
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		var e event
		if err := json.Unmarshal(line, &e); err != nil {
			continue // non-event output line
		}
		evs = append(evs, e)
	}
	return evs, sc.Err()
}

type report struct {
	OK            bool
	Built         bool
	BuildRequired bool
	Ran           int

	Regressions     []string // failed but not documented
	Resolved        []string // permanent divergence that passed
	UnknownExpected []string // permanent divergence that never ran
	Retirable       []string // not-implemented-yet test that now passes (soft)
	Handled         []string // documented and failed/skipped (the good case)
}

// evaluate compares observed test outcomes against the catalogue. expected maps a
// top-level test name to its category; buildRequired reports whether a build
// failure should be red.
func evaluate(events []event, expected map[string]divergence.Category, buildRequired bool) report {
	const passAction, failAction, skipAction = "pass", "fail", "skip"

	outcome := make(map[string]string) // top-level test name -> last terminal action
	for _, e := range events {
		if e.Test == "" || strings.Contains(e.Test, "/") {
			continue // package-level event or subtest; track top-level only
		}
		switch e.Action {
		case passAction, failAction, skipAction:
			outcome[e.Test] = e.Action
		}
	}

	var rep report
	rep.Ran = len(outcome)
	rep.Built = rep.Ran > 0
	rep.BuildRequired = buildRequired

	if !rep.Built {
		rep.OK = !buildRequired
		return rep
	}

	for name, act := range outcome {
		cat, isExpected := expected[name]
		permanent := isExpected && cat.Permanent()
		switch {
		case act == failAction && !isExpected:
			rep.Regressions = append(rep.Regressions, name)
		case act == passAction && permanent:
			rep.Resolved = append(rep.Resolved, name)
		case act == passAction && isExpected: // not-implemented-yet that now passes
			rep.Retirable = append(rep.Retirable, name)
		case isExpected && (act == failAction || act == skipAction):
			rep.Handled = append(rep.Handled, name)
		}
	}
	// A permanent divergence that never ran means the catalogue names an unknown
	// test. (not-implemented-yet entries are allowed to never run.)
	for name, cat := range expected {
		if !cat.Permanent() {
			continue
		}
		if _, ran := outcome[name]; !ran {
			rep.UnknownExpected = append(rep.UnknownExpected, name)
		}
	}

	sort.Strings(rep.Regressions)
	sort.Strings(rep.Resolved)
	sort.Strings(rep.UnknownExpected)
	sort.Strings(rep.Retirable)
	sort.Strings(rep.Handled)

	rep.OK = len(rep.Regressions) == 0 &&
		len(rep.Resolved) == 0 &&
		len(rep.UnknownExpected) == 0

	return rep
}

func (r report) String() string {
	var b strings.Builder
	b.WriteString("== conformance oracle ==\n")

	if !r.Built {
		if r.BuildRequired {
			b.WriteString("FAIL: the conformance suite did not build or ran no tests.\n")
		} else {
			b.WriteString("OK (bootstrapping): the suite does not build yet; conformance not enforced.\n")
			b.WriteString("     Flip divergences.json status to \"enforcing\" once v2 compiles.\n")
		}
		return b.String()
	}

	fmt.Fprintf(&b, "ran %d top-level tests\n", r.Ran)
	fmt.Fprintf(&b, "documented outcomes handled (failed/skipped as expected): %d\n", len(r.Handled))

	writeList := func(title string, names []string, hint string) {
		if len(names) == 0 {
			return
		}
		fmt.Fprintf(&b, "\n%s:\n", title)
		for _, n := range names {
			fmt.Fprintf(&b, "  - %s\n", n)
		}
		if hint != "" {
			fmt.Fprintf(&b, "  -> %s\n", hint)
		}
	}
	writeList("REGRESSIONS (undocumented failures)", r.Regressions,
		"fix the code, or add it to divergences.json (category not-implemented-yet while building out)")
	writeList("RESOLVED DIVERGENCES (permanent divergence now passes)", r.Resolved,
		"remove it from divergences.json")
	writeList("UNKNOWN EXPECTED (catalogue names a test that never ran)", r.UnknownExpected,
		"fix the test name in divergences.json")
	// Retirable is informational only; it does not fail the run.
	writeList("RETIRABLE (not-implemented-yet test now passes)", r.Retirable,
		"implemented — remove it from divergences.json")

	b.WriteString("\n")
	if r.OK {
		b.WriteString("PASS: outcomes match the documented catalogue.\n")
	} else {
		b.WriteString("FAIL: see above.\n")
	}
	return b.String()
}
