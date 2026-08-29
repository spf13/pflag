//go:build conformance

// Package testenv is a minimal stand-in for the standard library's
// internal/testenv package, providing just the helpers used by the vendored
// stdlib flag test suite (see ../../flag_go*_test.go).
//
// The stdlib test suite imports "internal/testenv", which is not importable
// from outside the Go tree. The conformance sync script rewrites that import
// to point here instead.
package testenv

import (
	"os"
	"testing"
)

// MustHaveExec checks that the current system can start new processes. The
// stdlib version skips on platforms (js/wasm, ios) where exec is unavailable;
// for our purposes a best-effort check is enough.
func MustHaveExec(t testing.TB) {
	if _, err := os.Executable(); err != nil {
		t.Skipf("skipping test: cannot determine executable: %v", err)
	}
}

// Executable returns the path to the current test binary, used by TestExitCode
// to re-exec itself as a child process.
func Executable(t testing.TB) string {
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable failed: %v", err)
	}
	return exe
}
