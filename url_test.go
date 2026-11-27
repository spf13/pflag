package pflag

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestURL(t *testing.T) {
	for _, input := range []string{
		"https://example.com/path?q=hello%20world#section",
		"../relative/path?key=value",
		"mailto:person@example.com",
		"https://user:pass@example.com:8080/a%2Fb",
		"",
	} {
		t.Run(input, func(t *testing.T) {
			flags := NewFlagSet("test", ContinueOnError)
			value := flags.URLP("endpoint", "e", url.URL{}, "server URL")
			if err := flags.Parse([]string{"-e", input}); err != nil {
				t.Fatal(err)
			}
			want, err := url.Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(value, want) {
				t.Fatalf("got %#v, want %#v", value, want)
			}
			got, err := flags.GetURL("endpoint")
			if err != nil || !reflect.DeepEqual(got, *want) {
				t.Fatalf("GetURL() = %#v, %v; want %#v", got, err, *want)
			}
		})
	}
}

func TestURLInvalidValuePreservesPreviousValue(t *testing.T) {
	for _, input := range []string{"https://example.com/%zz", "https://[invalid", "\n"} {
		flags := NewFlagSet("test", ContinueOnError)
		value := flags.URL("endpoint", url.URL{Scheme: "https", Host: "example.com"}, "")
		if err := flags.Set("endpoint", input); err == nil {
			t.Fatalf("Set(%q) succeeded", input)
		}
		if value.String() != "https://example.com" {
			t.Fatalf("invalid input changed value to %q", value)
		}
	}
}

func TestURLDefaultsAndRegistration(t *testing.T) {
	want := url.URL{Scheme: "https", Host: "example.com", Path: "/api"}
	flags := NewFlagSet("test", ContinueOnError)
	var value, shorthand url.URL
	flags.URLVar(&value, "endpoint", want, "server URL")
	flags.URLVarP(&shorthand, "alternate", "a", want, "other URL")
	if value != want || shorthand != want {
		t.Fatalf("defaults = %#v, %#v; want %#v", value, shorthand, want)
	}
	got, err := flags.GetURL("endpoint")
	if err != nil || got != want {
		t.Fatalf("default GetURL() = %#v, %v; want %#v", got, err, want)
	}
	if !strings.Contains(flags.FlagUsages(), "https://example.com/api") {
		t.Fatal("help does not include default URL")
	}
	if err := flags.Parse([]string{"--endpoint=../one", "-a", "../two"}); err != nil {
		t.Fatal(err)
	}
	if value.Path != "../one" || shorthand.Path != "../two" {
		t.Fatalf("parsed paths = %q, %q", value.Path, shorthand.Path)
	}
}

func TestGetURLErrors(t *testing.T) {
	flags := NewFlagSet("test", ContinueOnError)
	flags.String("text", "https://example.com", "")
	for _, name := range []string{"missing", "text"} {
		if _, err := flags.GetURL(name); err == nil {
			t.Fatalf("GetURL(%q) succeeded", name)
		}
	}
}

func TestURLHelpDefaults(t *testing.T) {
	for _, path := range []string{"", "0", "false", "relative"} {
		flags := NewFlagSet("test", ContinueOnError)
		flags.URL("endpoint", url.URL{Path: path}, "server URL")
		help := flags.FlagUsages()
		if path == "" {
			if strings.Contains(help, "(default") {
				t.Fatalf("empty URL has a displayed default: %q", help)
			}
		} else if !strings.Contains(help, "(default "+path+")") {
			t.Fatalf("help does not include URL default %q: %q", path, help)
		}
	}
}

func TestURLPackageHelpers(t *testing.T) {
	previous := CommandLine
	CommandLine = NewFlagSet("test", ContinueOnError)
	defer func() { CommandLine = previous }()
	var bound, shortBound url.URL
	URLVar(&bound, "bound", url.URL{}, "")
	URLVarP(&shortBound, "short-bound", "b", url.URL{}, "")
	value := URL("value", url.URL{}, "")
	shortValue := URLP("short-value", "v", url.URL{}, "")
	if err := CommandLine.Parse([]string{
		"--bound=/one", "-b", "/two", "--value=/three", "-v", "/four",
	}); err != nil {
		t.Fatal(err)
	}
	if bound.Path != "/one" || shortBound.Path != "/two" || value.Path != "/three" || shortValue.Path != "/four" {
		t.Fatalf("parsed paths = %q, %q, %q, %q", bound.Path, shortBound.Path, value.Path, shortValue.Path)
	}
}
