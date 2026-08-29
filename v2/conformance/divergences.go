package conformance

import (
	_ "embed"

	"github.com/spf13/pflag/v2/conformance/internal/divergence"
)

//go:embed divergences.json
var manifestJSON []byte

// Manifest returns the parsed divergence catalogue (divergences.json): the
// documented, intentional differences between pflag v2 and the stdlib flag
// package. The conformance oracle (hack/oracle) uses it to decide which test
// failures are expected.
func Manifest() (divergence.Manifest, error) {
	return divergence.Parse(manifestJSON)
}
