package pflag

import "strings"

// splitCommaSeparatedSliceValue preserves an explicit empty flag value as an
// empty slice instead of strings.Split's single empty-string element.
func splitCommaSeparatedSliceValue(val string) []string {
	if val == "" {
		return []string{}
	}
	return strings.Split(val, ",")
}
