package pflag

import (
	"net/url"
	"testing"
)

func TestURL(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	u := fs.URL("url", url.URL{Scheme: "https", Host: "example.com"}, "test url flag")
	args := []string{"--url", "https://go.dev/pkg/net/url/"}
	if err := fs.Parse(args); err != nil {
		t.Fatal(err)
	}
	if u.String() != "https://go.dev/pkg/net/url/" {
		t.Errorf("expected https://go.dev/pkg/net/url/, got %s", u.String())
	}
}

func TestURLDefault(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	u := fs.URL("url", url.URL{Scheme: "https", Host: "example.com"}, "test url flag")
	if err := fs.Parse([]string{}); err != nil {
		t.Fatal(err)
	}
	if u.String() != "https://example.com" {
		t.Errorf("expected https://example.com, got %s", u.String())
	}
}

func TestURLInvalid(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	fs.URL("url", url.URL{}, "test url flag")
	args := []string{"--url", "://invalid"}
	if err := fs.Parse(args); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestURLP(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	u := fs.URLP("url", "u", url.URL{}, "test url flag")
	args := []string{"-u", "https://example.com/path"}
	if err := fs.Parse(args); err != nil {
		t.Fatal(err)
	}
	if u.String() != "https://example.com/path" {
		t.Errorf("expected https://example.com/path, got %s", u.String())
	}
}

func TestURLChanged(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	fs.URL("url", url.URL{}, "test url flag")
	args := []string{"--url", "https://changed.com"}
	if err := fs.Parse(args); err != nil {
		t.Fatal(err)
	}
	if !fs.Lookup("url").Changed {
		t.Error("expected Changed to be true")
	}
}
