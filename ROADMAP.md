# `pflag` roadmap

## Beef up the team of maintainers

3 to 5 volunteers. No less, no more.

Would be nice to open a slack or discord channel to make debate more lively.
Otherwise, we may use github discussions.

## Agree on the [CONTRIBUTING guidelines](./CONTRIBUTING.md)

There are quite a few debatable things there: let's find a minimal consensus
(possibly defer decisions about some details to a later stage).

## Agree on a general direction

Proposed principles:

* promise of compatibility (i.e. no breaking `v2` ever)

> Example: solving issue #384 by adding support for a default value would break compatibility.
> If we want to solve this issue, we need a new API.

* prioritize use cases from `cobra` & `viper`, acknowledged as the main venues for importing `pflag` indirectly

> Example: solving issue #250 typically comes from a `cobra` user.

* prioritize code simplification & maintenance reduction over new features

> Examples: PRs #235, #248, #332, #369, #378, #380 are representative of that effort

* deflect stream of innovation, desire to share new ideas to some dedicated `contrib` repo

> Examples: issue #246 is typically such a new flag type, requiring a specific import.


## Issue triage

There are 88 outstanding issues. Let's triage them first.

For every one of them:
  * categorize: bug report, feature request/enhancement, code quality/testing, extensions ... (add new labels)
  * post a comment/reply to the OP

Proposed issue labels:
* bug
* duplicate
* invalid (not a bug)
* enhancement
* won't fix
* feature request
* extension
* code quality
* investigate/need repro

Questions and the threads of answers may be used to populate a FAQ.

## PR triage

There are 70 outstanding PRs. Let's triage them first.

Identify obsolete/won't do PRs, e.g. PR #206, probably #233 & #234 as well

Prioritize merges as follows:
  * bug fixes
  * code quality/linting/fixups & minor refact proposals (add new labels)
  * extra leeways/options (e.g. using the vararg `Option` pattern)

Same labels as for issues may be used.

At this point, isolate new features / enhancement proposals: they will be discussed later, once we get a clearer view.

New PRs may come in with more refactoring and/or improvements to testing (e.g. use testify, add fuzz test...).

## CI & tools

* [ ] Re-enact CI (github actions), with (i) tests, (ii) - later - linting, (iii) govulneck
* [ ] Simplify (no real need for the tools in `verify`)
* [ ] Enforce linting

## Repo climate

* update & revamp repo badges (code coverage, CI passed, `pkg.go.dev` etc)
* shorter README, defer details to other docs (may be published as wiki or github pages)

> I don't think we need to use github `Projects`: this is too small to justify this.
> Github `Discussions` could come in handy to collect thoughts from the community.

## Setup a dedicated repo for extensions

I suggest to create `github.com/spf13/pflag-contrib` to host such additions.

> I don't have a knack for punch lines and catchy headlines.
> There is most likely a better name to be found.

Identify PR candidates for extensions. We can manually clone or rewrite the most interesting ones.

Popular non-breaking extensions might be refitted into the main repository at some point in the future.

Examples:

* PR #247 (new type)
* PR #374 (flag validation)
* PR #114 (localization)
