// +build gofuzz

package pflag

import (
	"strings"
)

func Fuzz(data []byte) int {
	f := NewFlagSet("test", ContinueOnError)
	err := f.Parse(strings.Split(string(data), " "))
	if err != nil {
		return 0
	}
	return 1
}
