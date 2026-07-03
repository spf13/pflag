// Package conformance verifies pflag v2's two compatibility promises:
//
//   - Drop-in replacement for the Go standard library flag package. Tested by
//     vendored, build-tagged copies of the stdlib flag_test.go, one per
//     supported Go version (flag_go*_test.go).
//   - POSIX/GNU command-line argument syntax. Tested by a hand-written suite
//     keyed to the POSIX Utility Syntax Guidelines and GNU extensions
//     (gnuposix_test.go).
//
// The two specs genuinely conflict in a few places — that is why pflag exists.
// POSIX wins; those conflicts are catalogued in divergences.json and enforced by
// the oracle (hack/oracle), which gates CI: the suite is green iff the failing
// tests match the catalogue exactly.
//
// Everything here only compiles under the "conformance" build tag, so it does
// not affect the normal `go test ./...` run while v2 is still being built out.
// See README.md for details and run them with:
//
//	go test -tags conformance ./conformance/...
//
// This file exists so the directory always contains at least one buildable Go
// file regardless of build tags.
//
// Run `go generate ./conformance/...` to re-sync the vendored stdlib flag tests
// (it refreshes every version already vendored here; pass explicit versions to
// hack/sync.sh to add or drop one).
//
//go:generate ./hack/sync.sh
package conformance
