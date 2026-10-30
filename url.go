package pflag

import "net/url"

type urlValue url.URL

func newURLValue(value url.URL, p *url.URL) *urlValue {
	*p = value
	return (*urlValue)(p)
}

func (u *urlValue) String() string {
	value := url.URL(*u)
	return value.String()
}

func (u *urlValue) Set(s string) error {
	value, err := url.Parse(s)
	if err != nil {
		return err
	}
	*u = urlValue(*value)
	return nil
}

func (u *urlValue) Type() string {
	return "url"
}

func urlConv(s string) (interface{}, error) {
	value, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	return *value, nil
}

// GetURL returns the URL value of the flag with the given name.
func (f *FlagSet) GetURL(name string) (url.URL, error) {
	value, err := f.getFlagType(name, "url", urlConv)
	if err != nil {
		return url.URL{}, err
	}
	return value.(url.URL), nil
}

// URLVar defines a url.URL flag with the specified name, default value, and usage string.
// The argument p points to the url.URL variable in which to store the value.
// Values are parsed with url.Parse, which also accepts relative and opaque URLs.
func (f *FlagSet) URLVar(p *url.URL, name string, value url.URL, usage string) {
	f.URLVarP(p, name, "", value, usage)
}

// URLVarP is like URLVar, but accepts a shorthand letter after a single dash.
func (f *FlagSet) URLVarP(p *url.URL, name, shorthand string, value url.URL, usage string) {
	f.VarP(newURLValue(value, p), name, shorthand, usage)
}

// URLVar defines a url.URL flag with the specified name, default value, and usage string.
// The argument p points to the url.URL variable in which to store the value.
func URLVar(p *url.URL, name string, value url.URL, usage string) {
	CommandLine.URLVar(p, name, value, usage)
}

// URLVarP is like URLVar, but accepts a shorthand letter after a single dash.
func URLVarP(p *url.URL, name, shorthand string, value url.URL, usage string) {
	CommandLine.URLVarP(p, name, shorthand, value, usage)
}

// URL defines a url.URL flag with the specified name, default value, and usage string.
// The return value points to the url.URL variable that stores the value.
func (f *FlagSet) URL(name string, value url.URL, usage string) *url.URL {
	return f.URLP(name, "", value, usage)
}

// URLP is like URL, but accepts a shorthand letter after a single dash.
func (f *FlagSet) URLP(name, shorthand string, value url.URL, usage string) *url.URL {
	p := new(url.URL)
	f.URLVarP(p, name, shorthand, value, usage)
	return p
}

// URL defines a url.URL flag with the specified name, default value, and usage string.
// The return value points to the url.URL variable that stores the value.
func URL(name string, value url.URL, usage string) *url.URL {
	return CommandLine.URL(name, value, usage)
}

// URLP is like URL, but accepts a shorthand letter after a single dash.
func URLP(name, shorthand string, value url.URL, usage string) *url.URL {
	return CommandLine.URLP(name, shorthand, value, usage)
}
