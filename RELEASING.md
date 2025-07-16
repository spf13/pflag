# Releasing a new version

Releasing a new version is very simple, but documented here to enable any maintainer with enough access to help out.
We release new versions whenever there are changes that seem valuable to make available for consumption.

To create a new relase:

1. Go to the [releases page](https://github.com/spf13/pflag/releases) on Github, and select
[Draft a new release](https://github.com/spf13/pflag/releases/new).

2. Click the tag selector, and create a new tag named `vX.Y.Z`, where `X.Y.Z` is the semantic version number
   of the new version. (Use your judgement when deciding whether to increment patch or minor version; increment
   major version if there are deliberate breaking changes to exported APIs or documented behaviors).

3. Click `Generate release notes` to generate a writeup that includes all merged changes since the last release,
   and all first-time contributors. Add a section at the top with some info about the relase if there is any
   relevant info to share.

4. Either save your release as a draft (remember to re-generate release notes when you come back to publish it!)
   or publish it directly. Github will automatically create the appropriate git tag when the relase is published.
