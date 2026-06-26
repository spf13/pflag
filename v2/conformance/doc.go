// Package conformance runs the Go standard library's own flag package test
// suite against pflag v2 to verify the drop-in-replacement promise.
//
// The actual tests are vendored, build-tagged copies of the stdlib
// flag_test.go (one per supported Go version) and only compile under the
// "conformance" build tag, so they do not affect the normal `go test ./...`
// run while v2 is still being built out. See README.md for details and run
// them with:
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
