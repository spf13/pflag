# Standard-library `flag` conformance suite

pflag's headline promise is that it is a **drop-in replacement** for the Go
standard library's `flag` package. This directory runs `flag`'s *own* test suite
against pflag v2 to keep that promise honest: pflag v2 is dot-imported as `flag`,
and anything that fails to compile or fails to pass is, by definition, a gap in
drop-in compatibility.

## Running it

```sh
# from the v2 module root:
go test -tags conformance -v ./conformance/...

# re-sync every vendored version (refreshes whatever is already present):
go generate ./conformance/...

# add (or refresh) a specific version:
./conformance/hack/sync.sh 1.27

# drop a version — sync.sh never deletes, so remove its vendored file:
rm conformance/flag_go1.25_test.go
```

The set of `flag_go*_test.go` files is the source of truth for which versions
are vendored: `sync.sh <ver>` adds/refreshes, deleting the file drops, and
`sync.sh` with no args re-syncs whatever remains.

The suite only compiles under the `//go:build conformance` tag, so the normal
`go test ./...` run is unaffected while v2 is still being built out.

## Go-version awareness

A separate copy of upstream `flag_test.go` is vendored **per Go minor version**,
each pinned with a build constraint:

| File | Build constraint |
| --- | --- |
| `flag_go1.25_test.go` | `conformance && go1.25 && !go1.26` |
| `flag_go1.26_test.go` | `conformance && go1.26 && !go1.27` |

The toolchain running the tests therefore builds **exactly** the copy matching
its version. Two payoffs:

1. The CI matrix (`oldstable`, `stable`) tests each Go release against *that
   release's own* flag tests, so an upstream API or behavioral change in `flag`
   shows up as a separate pass/fail signal.
2. Re-running `sync.sh` and reviewing the diff makes upstream changes between
   versions visible in the repo.

`harness_test.go` contains a guard (`TestVendoredStdlibFlagTest`) that **fails
loudly** if the running toolchain has no vendored copy — so adding a new Go
version to the CI matrix without syncing is caught, rather than silently testing
nothing. Each vendored copy registers its version with that guard via a small
generated `init()` injected just below its imports.

## Layout

The tests live in the external test package `conformance_test` (black-box tests
of v2's *public* API, mirroring upstream's `flag_test`); `doc.go` provides the
primary `conformance` package.

| File | Origin | Edit policy |
| --- | --- | --- |
| `doc.go` | hand-written | package doc; keeps the dir buildable with no tags; holds the `//go:generate` directive |
| `flag_go<ver>_test.go` | generated from that version's `flag_test.go` | **do not edit** — regenerate with `sync.sh` |
| `harness_test.go` | hand-written | reimplements the stdlib's internal `export_test.go` helpers (`ResetForTesting`, `DefaultUsage`) against v2's public API, plus the version guard |
| `internal/testenv/testenv.go` | hand-written | minimal stand-in for the stdlib-internal `internal/testenv` (only `MustHaveExec` / `Executable`) |
| `hack/sync.sh` | hand-written | regenerates the vendored copies |

### Why a sync script instead of a verbatim copy?

External package conventions and stdlib-internal imports force a few mechanical
edits, which the script applies and nothing else:

- `package flag_test` → `package conformance_test`
- `. "flag"` → `. "github.com/spf13/pflag/v2"` (kept a dot import)
- `"internal/testenv"` → the local shim (not importable outside the Go tree)
- inject a one-line `init()` after the imports registering the version with the guard
- prepend the build tags + a "generated" banner

For a version not matching the local toolchain, `sync.sh` fetches `flag_test.go`
from the matching Go release branch on GitHub.

---

## The gap list (TODO for the v2 implementation)

This is the API the conformance suite requires. As of writing, v2 is an empty
package, so **all of it** is outstanding — that is intended; the suite is the
spec. Check items off as `go test -tags conformance ./conformance/...` gets
further.

### 1. Package-level API surface

Constructors / globals:

- [ ] `type FlagSet` (and a usable zero value — `var flags FlagSet` is used)
- [ ] `type Flag struct { Name string; Usage string; Value Value; DefValue string }`
- [ ] `type Value interface { String() string; Set(string) error }`
- [ ] `type Getter interface { Value; Get() any }`
- [ ] `type ErrorHandling int` with `ContinueOnError`, `ExitOnError` (and `PanicOnError`)
- [ ] `var CommandLine *FlagSet` (must be assignable — `ResetForTesting` replaces it)
- [ ] `var Usage func()` (assignable package var)
- [ ] `var ErrHelp error`
- [ ] `func NewFlagSet(name string, errorHandling ErrorHandling) *FlagSet`

Package-level convenience wrappers over `CommandLine`:

- [ ] `Bool` `Int` `Int64` `Uint` `Uint64` `String` `Float64` `Duration`
- [ ] `Func` `BoolFunc`
- [ ] `Set` `Parse` `Visit` `VisitAll` `Arg` `Args`

### 2. `*FlagSet` methods

- [ ] Definers: `Bool` `BoolVar` `Int` `IntVar` `Int64` `Uint` `Uint64` `String`
      `Float64` `Duration` `Var` `Func` `BoolFunc`
- [ ] Parsing: `Parse` `Parsed`
- [ ] Args: `Args` `Arg`
- [ ] Mutation/inspection: `Set` `Visit` `VisitAll`
- [ ] Config/identity: `Init` `Name` `ErrorHandling` `SetOutput` `Output`
- [ ] Usage: `PrintDefaults`, exported `Usage func()` field
- [ ] `Var` flags must populate `Flag.DefValue` from `Value.String()`

### 3. Behavioral conformance (won't show as compile errors — the hard part)

Once the API compiles, these stdlib tests assert behavior that differs from
pflag v1's POSIX/GNU conventions. Each needs a deliberate decision: match stdlib,
or document as an accepted divergence and adjust/skip the assertion.

- [ ] **Single-dash long flags.** `testParse` passes `-bool`, `-int`, `-uint`,
      `-string`, `-float64`, `-duration` (single dash, long name). pflag v1 reads
      `-int` as the shorthand cluster `-i -n -t`. This is the single biggest
      drop-in gap and affects `TestParse`, `TestFlagSetParse`, `TestUserDefined*`,
      `TestHelp`, `TestExitCode`.
- [ ] **`-flag=value` and `-flag value` forms** for all types (`-bool2=true`,
      `--int 22`).
- [ ] **Hex/base-0 integer parsing.** `testParse` sets `--int64 0x23` and expects
      `0x23`. Int/Int64/Uint/Uint64 must parse with base 0.
- [ ] **`PrintDefaults` exact format.** `TestPrintDefaults` and
      `TestUserDefinedBoolUsage` compare byte-for-byte against the stdlib layout
      (`  -A\tfor ...`, back-quoted `name` extraction, multiline indent, panic
      recovery for a `String()` that panics on a zero value).
- [ ] **Error message text.** `TestParseError` wants `invalid` + `parse error`;
      `TestRangeError` wants `invalid` + `value out of range`; `TestUsageOutput`
      wants exactly `flag provided but not defined: -i\nUsage of app:\n`.
- [ ] **`Var` validation panics.** `TestInvalidFlags`: name starting with `-`
      panics `flag "-foo" begins with -`; name containing `=` panics
      `flag "foo=bar" contains =`. `TestRedefinedFlags`: `flag redefined: foo`.
- [ ] **Define-after-set panic.** `TestDefineAfterSet` expects a panic matching
      `flag myFlag set at .*:.* before being defined`.
- [ ] **`Getter` semantics.** `TestGet`: every built-in value satisfies `Getter`
      and `Get()` returns the natural Go type (`bool`, `int`, `int64`, `uint`,
      `uint64`, `string`, `float64`, `time.Duration`). The `Func` value's
      `String()` returns `""` and it is **not** a `Getter`.
- [ ] **`IsBoolFlag` user types.** `TestUserDefinedBool` relies on a custom
      `Value` with `IsBoolFlag()` toggling between bool-like and value-like
      parsing mid-parse.
- [ ] **`-h`/`-help` ⇒ `ErrHelp`** when undefined, overridable by a defined flag
      (`TestHelp`), and exit codes 0/2/magic (`TestExitCode`, runs a child proc).
- [ ] **Int overflow** is a parse error (`TestIntFlagOverflow`, 32-bit only).
- [ ] **Visit ordering** is lexical/sorted (`TestEverything`).
