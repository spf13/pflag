package pflag

import (
	"errors"
	"testing"
)

func TestNotExistError(t *testing.T) {
	err := &NotExistError{
		name:                "foo",
		specifiedShorthands: "bar",
	}

	if err.GetSpecifiedName() != "foo" {
		t.Errorf("Expected GetSpecifiedName to return %q, got %q", "foo", err.GetSpecifiedName())
	}
	if err.GetSpecifiedShortnames() != "bar" {
		t.Errorf("Expected GetSpecifiedShortnames to return %q, got %q", "bar", err.GetSpecifiedShortnames())
	}
}

func TestValueRequiredError(t *testing.T) {
	err := &ValueRequiredError{
		flag:                &Flag{},
		specifiedName:       "foo",
		specifiedShorthands: "bar",
	}

	if err.GetFlag() == nil {
		t.Error("Expected GetSpecifiedName to return its flag field, but got nil")
	}
	if err.GetSpecifiedName() != "foo" {
		t.Errorf("Expected GetSpecifiedName to return %q, got %q", "foo", err.GetSpecifiedName())
	}
	if err.GetSpecifiedShortnames() != "bar" {
		t.Errorf("Expected GetSpecifiedShortnames to return %q, got %q", "bar", err.GetSpecifiedShortnames())
	}
}

func TestLongFlagSingleDashError(t *testing.T) {
	err := &LongFlagSingleDashError{name: "name"}
	if err.GetName() != "name" {
		t.Errorf("expected name %q, got %q", "name", err.GetName())
	}
	if err.Error() != `bad flag syntax: -name; did you mean --name?` {
		t.Errorf("unexpected error: %q", err.Error())
	}
}

func TestParseLongFlagWithSingleDash(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	var name string
	f.StringVarP(&name, "name", "n", "", "name")
	err := f.Parse([]string{"-name", "wrong"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var longDash *LongFlagSingleDashError
	if !errors.As(err, &longDash) {
		t.Fatalf("expected LongFlagSingleDashError, got %T: %v", err, err)
	}
}

func TestInvalidValueError(t *testing.T) {
	expectedCause := errors.New("error")
	err := &InvalidValueError{
		flag:  &Flag{},
		value: "foo",
		cause: expectedCause,
	}

	if err.GetFlag() == nil {
		t.Error("Expected GetSpecifiedName to return its flag field, but got nil")
	}
	if err.GetValue() != "foo" {
		t.Errorf("Expected GetValue to return %q, got %q", "foo", err.GetValue())
	}
	if actual := err.Unwrap(); actual != expectedCause { //nolint:errorlint // not using errors.Is for compatibility with go1.12
		t.Errorf("Expected Unwrwap to return %q, got %q", expectedCause, actual)
	}
}

func TestInvalidSyntaxError(t *testing.T) {
	err := &InvalidSyntaxError{
		specifiedFlag: "--=",
	}

	if err.GetSpecifiedFlag() != "--=" {
		t.Errorf("Expected GetSpecifiedFlag to return %q, got %q", "--=", err.GetSpecifiedFlag())
	}
}
