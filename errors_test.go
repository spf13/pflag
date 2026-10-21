package pflag

import (
	"errors"
	"testing"
)

func TestNotExistError(t *testing.T) {
	err := &NotExistError{
		name:                "foo",
		specifiedShorthands: "bar",
		messageType:         flagNotExistMessage,
	}

	if err.GetSpecifiedName() != "foo" {
		t.Errorf("Expected GetSpecifiedName to return %q, got %q", "foo", err.GetSpecifiedName())
	}
	if err.GetSpecifiedShortnames() != "bar" {
		t.Errorf("Expected GetSpecifiedShortnames to return %q, got %q", "bar", err.GetSpecifiedShortnames())
	}
	if got := err.Error(); got != `flag "foo" does not exist` {
		t.Errorf("Expected Error() to return %q, got %q", `flag "foo" does not exist`, got)
	}

	errNoSuch := &NotExistError{name: "f", messageType: flagNoSuchFlagMessage}
	if got := errNoSuch.Error(); got != "no such flag -f" {
		t.Errorf("Expected Error() to return %q, got %q", "no such flag -f", got)
	}

	errUnknown := &NotExistError{name: "foo", messageType: flagUnknownFlagMessage}
	if got := errUnknown.Error(); got != "unknown flag: --foo" {
		t.Errorf("Expected Error() to return %q, got %q", "unknown flag: --foo", got)
	}

	errUnknownShort := &NotExistError{name: "f", specifiedShorthands: "bar", messageType: flagUnknownShorthandFlagMessage}
	if got := errUnknownShort.Error(); got != `unknown shorthand flag: 'f' in -bar` {
		t.Errorf("Expected Error() to return %q, got %q", `unknown shorthand flag: 'f' in -bar`, got)
	}
}

func TestValueRequiredError(t *testing.T) {
	err := &ValueRequiredError{
		flag:                &Flag{},
		specifiedName:       "foo",
		specifiedShorthands: "bar",
	}

	if err.GetFlag() == nil {
		t.Error("Expected GetFlag to return its flag field, but got nil")
	}
	if err.GetSpecifiedName() != "foo" {
		t.Errorf("Expected GetSpecifiedName to return %q, got %q", "foo", err.GetSpecifiedName())
	}
	if err.GetSpecifiedShortnames() != "bar" {
		t.Errorf("Expected GetSpecifiedShortnames to return %q, got %q", "bar", err.GetSpecifiedShortnames())
	}
	if got := err.Error(); got != `flag needs an argument: 'f' in -bar` {
		t.Errorf("Expected Error() to return %q, got %q", `flag needs an argument: 'f' in -bar`, got)
	}

	errLong := &ValueRequiredError{
		flag:          &Flag{},
		specifiedName: "foo",
	}
	if got := errLong.Error(); got != "flag needs an argument: --foo" {
		t.Errorf("Expected Error() to return %q, got %q", "flag needs an argument: --foo", got)
	}
}

func TestInvalidValueError(t *testing.T) {
	expectedCause := errors.New("error")
	err := &InvalidValueError{
		flag: &Flag{
			Name:      "foo",
			Shorthand: "f",
		},
		value: "bar",
		cause: expectedCause,
	}

	if err.GetFlag() == nil {
		t.Error("Expected GetFlag to return its flag field, but got nil")
	}
	if err.GetValue() != "bar" {
		t.Errorf("Expected GetValue to return %q, got %q", "bar", err.GetValue())
	}
	if actual := err.Unwrap(); actual != expectedCause { //nolint:errorlint // not using errors.Is for compatibility with go1.12
		t.Errorf("Expected Unwrap to return %q, got %q", expectedCause, actual)
	}
	if got := err.Error(); got != `invalid argument "bar" for "-f, --foo" flag: error` {
		t.Errorf("Expected Error() to return %q, got %q", `invalid argument "bar" for "-f, --foo" flag: error`, got)
	}

	errNoShort := &InvalidValueError{
		flag:  &Flag{Name: "foo"},
		value: "bar",
		cause: expectedCause,
	}
	if got := errNoShort.Error(); got != `invalid argument "bar" for "--foo" flag: error` {
		t.Errorf("Expected Error() to return %q, got %q", `invalid argument "bar" for "--foo" flag: error`, got)
	}
}

func TestInvalidSyntaxError(t *testing.T) {
	err := &InvalidSyntaxError{
		specifiedFlag: "--=",
	}

	if err.GetSpecifiedFlag() != "--=" {
		t.Errorf("Expected GetSpecifiedFlag to return %q, got %q", "--=", err.GetSpecifiedFlag())
	}
	if got := err.Error(); got != "bad flag syntax: --=" {
		t.Errorf("Expected Error() to return %q, got %q", "bad flag syntax: --=", got)
	}
}
