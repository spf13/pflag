package pflag

import (
	"fmt"
	"strconv"
)

// collectionValueError reports the invalid element without exposing the
// parser's implementation details, while retaining its error for callers.
type collectionValueError struct {
	value string
	hint  string
	cause error
}

func (e *collectionValueError) Error() string {
	hint := e.hint
	if numErr, ok := e.cause.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange { //nolint:errorlint // strconv returns *NumError directly; retain Go 1.12 compatibility.
		hint = "is out of range"
	}
	return fmt.Sprintf("%q %s", e.value, hint)
}

func (e *collectionValueError) Unwrap() error { return e.cause }
