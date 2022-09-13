# carapace-pflag

[![PkgGoDev](https://pkg.go.dev/badge/github.com/rsteube/carapace-pflag)](https://pkg.go.dev/github.com/rsteube/carapace-pflag)
[![GoReportCard](https://goreportcard.com/badge/github.com/rsteube/carapace-pflag)](https://goreportcard.com/report/github.com/rsteube/carapace-pflag)
[![Coverage Status](https://coveralls.io/repos/github/rsteube/carapace-pflag/badge.svg?branch=master)](https://coveralls.io/github/rsteube/carapace-pflag?branch=master)


## Customizations

### Shorthand-Only

Support shorthand-only flags (e.g. `-h` without `--help`) using `S` suffix for flag functions:

```go
pflag.BoolS("help", "h", false, "show help")
```
