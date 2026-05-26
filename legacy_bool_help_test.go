package pflag

import (
	"bytes"
	"strings"
	"testing"
)

func TestLegacyBoolHelp(t *testing.T) {
	buf := bytes.Buffer{}
	f := NewFlagSet("test", ContinueOnError)
	f.LegacyBoolHelp = true
	f.Bool("enabled", false, "Enable the feature")
	f.Bool("verbose", true, "Verbose logging")
	f.SetOutput(&buf)
	f.PrintDefaults()

	got := buf.String()
	if strings.Contains(got, "[=true|false]") {
		t.Fatalf("legacy bool help should omit [=true|false], got:\n%s", got)
	}
	if strings.Contains(got, "(default true)") || strings.Contains(got, "(default false)") {
		t.Fatalf("legacy bool help should omit bool defaults, got:\n%s", got)
	}
	if !strings.Contains(got, "Enable the feature") || !strings.Contains(got, "Verbose logging") {
		t.Fatalf("expected usage text in output, got:\n%s", got)
	}
}

func TestLegacyBoolHelpBoolFunc(t *testing.T) {
	buf := bytes.Buffer{}
	f := NewFlagSet("test", ContinueOnError)
	f.LegacyBoolHelp = true
	f.BoolFunc("callback", "Run callback", func(string) error { return nil })
	f.SetOutput(&buf)
	f.PrintDefaults()

	got := buf.String()
	if strings.Contains(got, "(default ") {
		t.Fatalf("legacy bool help should omit boolfunc defaults, got:\n%s", got)
	}
	if !strings.Contains(got, "Run callback") {
		t.Fatalf("expected usage text in output, got:\n%s", got)
	}
}
