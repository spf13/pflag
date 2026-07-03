# Conformance suites

pflag makes two compatibility promises, and this directory keeps both honest
with executable specs that run against pflag v2:

1. **Drop-in replacement for the standard library `flag` package.** We run
   `flag`'s *own* test suite against pflag v2 (dot-imported as `flag`) —
   `flag_go*_test.go`. See [Standard-library `flag`](#standard-library-flag-suite).
2. **POSIX/GNU command-line argument syntax.** A hand-written suite keyed to the
   POSIX Utility Syntax Guidelines and GNU extensions — `gnuposix_test.go`. See
   [POSIX/GNU argument syntax](#posixgnu-argument-syntax-suite).

Anything that fails to compile or fails to pass is, by definition, a gap.

## POSIX wins, and the oracle gates on it

The two specs genuinely conflict in a few places — single-dash `-int` is the flag
`int` to the stdlib, but the cluster `-i -n -t` to POSIX. **That conflict is the
whole reason pflag exists, and POSIX wins.** So some vendored stdlib tests are
*expected* to fail, and the stdlib suite can never be all-green.

The success criterion is therefore not "everything passes" but:

> The suite builds and runs, and the set of failing tests matches the documented
> divergence catalogue **exactly** — no undocumented failure (a regression) and
> no documented test that quietly started passing (a divergence to retire).

That catalogue is [`divergences.json`](divergences.json) — the single source of
truth, listing each conflict, the POSIX Guideline / GNU rule that decides it, and
the exact tests it makes fail (categories: `posix-overrides-stdlib`,
`pflag-design-differs`, `pflag-omits-gnu-feature`). The oracle in
[`hack/oracle`](hack/oracle) reads `go test -json` and enforces the criterion
above; CI gates on it:

```sh
go test -tags conformance -json ./conformance | go run ./conformance/hack/oracle
```

The oracle tells you exactly what to do when it's red: fix a regression, add a
newly-discovered conflict to the catalogue, or retire one that no longer applies.
That is how the catalogue gets **calibrated** against the real implementation —
the current entries are a best-effort seed.

`go test ./...` (no tag) validates that `divergences.json` is well-formed
(`TestDivergenceManifest`) without needing v2 to compile.

### Lifecycle: building v2 with a green CI

The catalogue's top-level `status` and the `not-implemented-yet` category let the
job stay green while v2 is built out, tightening as it matures:

| Phase | `status` | catalogue | oracle |
| --- | --- | --- | --- |
| v2 empty (now) | `bootstrapping` | permanent divergences only | build failure tolerated → **green** |
| v2 compiles, partial | `enforcing` | add `not-implemented-yet` entries for tests that fail only because a feature is missing | build required; those failures tolerated → **green** |
| v2 complete | `enforcing` | `not-implemented-yet` burned down to empty | only permanent divergences remain |

`not-implemented-yet` is lenient: a listed test may fail, be skipped, or not run.
When one starts **passing**, the oracle stays green and just lists it as
*retirable* — delete the entry. (A *permanent* divergence that passes is the
opposite: a hard failure, because it never should.) So implementing a feature
never turns CI red; forgetting to categorise a new failure does.

# Standard-library `flag` suite

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
| `gnuposix_test.go` | hand-written | the POSIX/GNU argument-syntax suite |
| `divergences.json` | hand-written | the catalogue of intentional differences (single source of truth) |
| `divergences.go` / `divergences_test.go` | hand-written | embed + validate the catalogue (untagged, runs in normal `go test`) |
| `internal/divergence/` | hand-written | catalogue types + parsing, shared by the test and the oracle (no build tag, no v2 dependency) |
| `hack/oracle/` | hand-written | the CI gate: compares `go test -json` against the catalogue |
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

# POSIX/GNU argument syntax suite

`gnuposix_test.go` checks pflag's *other* promise: compatibility with the POSIX
Utility Syntax Guidelines and the GNU long-option extensions. Each test cites the
rule it covers — `Gn` for a numbered POSIX Guideline, or "GNU" for an extension:

| Area | Rules covered |
| --- | --- |
| Short options | preceded by `-` (G4), single alphanumeric name (G3), clustering (G5), separate or attached arg (G6), cluster ending in an arg-taking option (G5) |
| Long options (GNU) | `--name`, dashes in names, `--name=value`, `--name value` |
| Special tokens | `--` ends options (G10), lone `-` is an operand (G13) |
| Ordering | GNU interspersing (default) vs strict POSIX stop-at-first-operand (G9), order-independence and repetition (G11) |

### Why hand-written instead of vendored?

Unlike the stdlib `flag` suite, there is **no reusable Go-native POSIX
conformance corpus to vendor**. The authoritative sources are prose — the
[Open Group Utility Conventions ch. 12][posix] (Guidelines 1–14) and
[GNU Argument Syntax][gnu] — and the machine-runnable suites (glibc/gnulib
`tst-getopt*.c`) are C, coupled to the C `getopt`/optstring API, so they mostly
exercise C-isms pflag does not share. The tests are therefore authored directly
from the guidelines. They are shaped as *parse → normalized result*, so if a
shared cross-language corpus ever appears we can expose a tiny CLI to drive it
without restructuring.

[posix]: https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html
[gnu]: https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html

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
