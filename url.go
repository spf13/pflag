package pflag

import (
	"fmt"
	"net/url"
)

// -- url.Value
type urlValue struct {
	value   *url.URL
	changed bool
}

func newURLValue(val *url.URL, p *url.URL) *urlValue {
	uv := new(urlValue)
	uv.value = p
	*uv.value = *val
	return uv
}

func (u *urlValue) Set(s string) error {
	v, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	*u.value = *v
	u.changed = true
	return nil
}

func (u *urlValue) Type() string {
	return "url"
}

func (u *urlValue) String() string {
	if u.value == nil {
		return ""
	}
	return u.value.String()
}

// URLVar defines a url.URL flag with specified name, default value, and usage string.
func (f *FlagSet) URLVar(p *url.URL, name string, value url.URL, usage string) {
	f.VarP(newURLValue(&value, p), name, "", usage)
}

// URLVarP is like URLVar, but accepts a shorthand letter.
func (f *FlagSet) URLVarP(p *url.URL, name, shorthand string, value url.URL, usage string) {
	f.VarP(newURLValue(&value, p), name, shorthand, usage)
}

// URL defines a url.URL flag with specified name, default value, and usage string.
func (f *FlagSet) URL(name string, value url.URL, usage string) *url.URL {
	p := new(url.URL)
	f.URLVarP(p, name, "", value, usage)
	return p
}

// URLP is like URL, but accepts a shorthand letter.
func (f *FlagSet) URLP(name, shorthand string, value url.URL, usage string) *url.URL {
	p := new(url.URL)
	f.URLVarP(p, name, shorthand, value, usage)
	return p
}

// URLVar defines a url.URL flag with specified name, default value, and usage string.
func URLVar(p *url.URL, name string, value url.URL, usage string) {
	CommandLine.VarP(newURLValue(&value, p), name, "", usage)
}

// URLVarP is like URLVar, but accepts a shorthand letter.
func URLVarP(p *url.URL, name, shorthand string, value url.URL, usage string) {
	CommandLine.VarP(newURLValue(&value, p), name, shorthand, usage)
}

// URL defines a url.URL flag with specified name, default value, and usage string.
func URL(name string, value url.URL, usage string) *url.URL {
	p := new(url.URL)
	CommandLine.URLVarP(p, name, "", value, usage)
	return p
}

// URLP is like URL, but accepts a shorthand letter.
func URLP(name, shorthand string, value url.URL, usage string) *url.URL {
	p := new(url.URL)
	CommandLine.URLVarP(p, name, shorthand, value, usage)
	return p
}
