//go:build conformance

package conformance_test

// This file is the hand-maintained glue for the vendored standard-library flag
// test suites (flag_go*_test.go). Keep it small: it is the only part of the
// conformance harness that is allowed to diverge from upstream.

import (
	"io"
	"os"
	"runtime"
	"strings"
	"testing"

	. "github.com/spf13/pflag/v2"
)

// vendoredVersion is set by the init() injected into the vendored flag test
// (flag_go*_test.go) whose build constraint matches the running toolchain. It
// stays empty when the suite is built with a Go version we have not vendored a
// copy for.
var vendoredVersion string

// TestVendoredStdlibFlagTest guards against silently testing nothing: if the
// toolchain has no matching vendored flag_test.go, no version-tagged file is
// compiled in and the rest of the suite would vacuously pass.
func TestVendoredStdlibFlagTest(t *testing.T) {
	if vendoredVersion == "" {
		t.Fatalf("no vendored stdlib flag test for %s; add one with: ./conformance/hack/sync.sh %s",
			runtime.Version(), minorOf(runtime.Version()))
	}
	if !strings.HasPrefix(runtime.Version(), vendoredVersion) {
		t.Fatalf("vendored copy %q does not match running toolchain %q; re-run ./conformance/hack/sync.sh",
			vendoredVersion, runtime.Version())
	}
}

// minorOf turns "go1.26.4" into "1.26" for the hint message.
func minorOf(v string) string {
	v = strings.TrimPrefix(v, "go")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) < 2 {
		return v
	}
	return parts[0] + "." + parts[1]
}

// DefaultUsage mirrors flag.DefaultUsage: a snapshot of the package-level Usage
// function captured before any test overrides it.
var DefaultUsage = Usage

// ResetForTesting mirrors flag.ResetForTesting: it replaces the global
// CommandLine with a fresh ContinueOnError FlagSet whose output is discarded,
// and installs the provided usage function.
func ResetForTesting(usage func()) {
	CommandLine = NewFlagSet(os.Args[0], ContinueOnError)
	CommandLine.SetOutput(io.Discard)
	CommandLine.Usage = func() { Usage() }
	Usage = usage
}
