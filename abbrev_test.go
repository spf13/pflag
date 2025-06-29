package pflag

import (
	"fmt"
	"runtime"
	"testing"
)

func TestSimple(t *testing.T) {
	assert := newAsserter(t)

	words := []string{"hello", "help", "sync", "uint", "uint16", "uint64"}
	ret := map[string]string{
		"hello":  "hello",
		"hell":   "hello",
		"help":   "help",
		"sync":   "sync",
		"syn":    "sync",
		"sy":     "sync",
		"s":      "sync",
		"uint":   "uint",
		"uint16": "uint16",
		"uint1":  "uint16",
		"uint64": "uint64",
		"uint6":  "uint64",
	}

	ab := abbrev(words)

	for k, v := range ab {
		x, ok := ret[k]
		assert(ok, "unexpected abbrev %s", k)
		assert(x == v, "abbrev %s: exp %s, saw %s", k, x, v)
	}

	for k, v := range ret {
		x, ok := ab[k]
		assert(ok, "unknown abbrev %s", k)
		assert(x == v, "abbrev %s: exp %s, saw %s", k, x, v)
	}

}

// make an assert() function for use in environment 't' and return it
func newAsserter(t *testing.T) func(cond bool, msg string, args ...interface{}) {
	return func(cond bool, msg string, args ...interface{}) {
		if cond {
			return
		}

		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}

		s := fmt.Sprintf(msg, args...)
		t.Fatalf("%s: %d: Assertion failed: %s\n", file, line, s)
	}
}
