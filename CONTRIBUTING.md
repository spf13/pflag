# Contributing guidelines

## Principles and vision

### [Roadmap to `v1.1.x`](./ROADMAP.md)

### go version

This repository supports _at least_ the 2 latest stable versions of the language.

Example: when `go1.21` is out, we may require `go1.19`

Any new language feature that is introduced should be guarded (and ideally, emulated for backward compatibility)
by go build flags, e.g. `go:build !go1.21`, `go:build go1.21`.

We don't upgrade the required version of the compiler automatically:
changes to the go version requirement in `go.mod` will occur only if a new language feature is deemed really useful.

Therefore, at any given time, older versions of the compiler might be supported. This not a promise, just a likely outcome.

Example: if we introduce generics, `go1.18` would be a requirement.

### The `pflag` compatibility promise

**Programs using older versions of `pflag` should be able to compile in the future.**

> In the spirit of the `go` language, we believe that the `golang` eco-system would benefit from a similar promise by this repository.

The required version of the compiler may however evolve.

We may use deprecation notices to advise users to move towards newer APIs.

Older APIs will continue being honored, though.

### The `pflag` promise of being a drop-in replacement of `flag`

`pflag` is a drop-in replacement for `go`'s flag package, implementing
POSIX/GNU-style --flags. That's it.

`pflag` will follow the evolution of the standard library `flag`, and continue to maintain the promise of being
a _drop-in replacement_ for `flag` (and then some of course, but at least that).

### Versioning

This package abides by `semver`.

Hence:
* `v1.0.6` : patch
* `v1.1.0` : new feature

As stated above, it is unlikely we ever issue a `v2.0.0` (breaking change).

## Pull requests

All contributions are much appreciated and welcome. There is no small contributions: additional documentation, examples, 
tests, etc are all valuable additions, not just new features.

### Linting

This repository uses `github.com/golangci/golangci-lint/cmd/golangci-lint` as meta-linter, enforced by CI.
The set of rules may evolve over time. `revive` and `govet` are must-have.

> TODO: agree on the initial .golangci.yml setup, with a no-nonsense approach to linting. Code quality without nitpicking... :)
> As a started, here is a sample configuration that I use on a daily basis for personal projects (no strong agenda on this):
> [golangci-lint.yml](https://github.com/fredbi/gflag/blob/master/.golangci.yml).

### Testing

Pull requests should include significant unit tests.

This repository uses `github.com/stretchr/testify` framework for unit tests.

> TODO: to be discussed, what we mean by "significant", e.g. >70% test coverage or so.

### Documentation

Keeping a compelling `godoc` documentation is part of the "drop-in replacement" promise.

Any newly exposed API should come with a clear, readable description of what it does.

Testable example are alway welcome (although not required).

Conversely, internals only need comments if a particular caveat or expectation requires wording besides just code.

### Dependencies

So far, `pflag` has managed to build without _any_ external dependency.

We may indulge into _a few_ dependencies, especially test dependencies, or things like `golang.org/x/...`.

We should not need much more than that, unless we are talking about extensions (e.g. new flag types). See below.

### Commit etiquette

Please keep your commit title and description concise and to the point.

Feel free to `git rebase -i`, `git commit --amend` etc, so this repository keeps a nice and readable git history,
without merge commits or fixup rescindments.

### Contributor Licence Agreement

Contributions are subject to a CLA DCO being signed off.

> TODO: is this something that is still needed? Anyhow, check if the CLA terms are up to date.
> The Linux Foundation has introduced [EasyCLA](https://docs.linuxfoundation.org/lfx/easycla).
> I believe that `cla.assistant.io` is a bit old.

### Signed commits

Commits shall be [signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits).

> TODO: this requirement is debatable.
> On the one hand, it is certainly something many junior developers are perplex about.
> On the other hand, supporting so many major pieces of infrastructure is a big responsibility.

### Linear history

Please rebase your PRs to the current master before a merging can be done. All PR's will be squashed.

## Contributed extensions

New types, great features built on top of `pflag` etc are welcome.

Considering adding those to the already rather rich collection of supported flag types is going always to be a long process.

To address that, a new repository `spf13/pflag-contrib` may host such new ideas and (possibly breaking) innovations.

That repository works under basically the same set of rules, but less stringent.

* backward-compatibility for a contributed feature remains desirable, but not required
* contributors have a freer hand regarding dependencies
* documentation & testing are needed, possibly with a lighter touch
* introduction of new go features is possible - just guard your extension with `go: build gox.yz` 
  to avoid breaking other contributions

> I am proposing a new git repo as I found rather impractical to maintain a nested `go.mod` in the main repo.
> Further, we most likely want to version contributed features independantly.
