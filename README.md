# carapace-pflag

[![PkgGoDev](https://pkg.go.dev/badge/github.com/rsteube/carapace-pflag)](https://pkg.go.dev/github.com/rsteube/carapace-pflag)
[![GoReportCard](https://goreportcard.com/badge/github.com/rsteube/carapace-pflag)](https://goreportcard.com/report/github.com/rsteube/carapace-pflag)
[![Coverage Status](https://coveralls.io/repos/github/rsteube/carapace-pflag/badge.svg?branch=master)](https://coveralls.io/github/rsteube/carapace-pflag?branch=master)

Fork of [spf13/pflag](https://github.com/spf13/pflag) aimed to provide support for non-posix variants in [carapace](https://github.com/rsteube/carapace) standalone mode (e.g. [carapace-bin](https://github.com/rsteube/carapace-bin)).

## Customizations

### Shorthand-Only

Support shorthand-only flags (e.g. `-h` without `--help`) using `S` suffix for flag functions:

```go
pflag.BoolS("help", "h", false, "show help") // -h
```

### Long Shorthand

Support shorthand flags that are more than one character long (e.g. `pkill -<sig>`).
> This implicitly disables posix shorthand chaining (e.g. `ls -lha`):

```go
pflag.BoolS("STOP", "STOP", false, "Stop process, unblockable") // -STOP
pflag.BoolN("help", "h", false, "show help") // -h, -help
```
