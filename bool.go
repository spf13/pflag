package pflag

import "strconv"

// optional interface to indicate boolean flags that can be
// supplied without "=value" text
type boolFlag interface {
	Value
	IsBoolFlag() bool
}

// -- bool Value
type boolValue bool

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Type() string {
	return "bool"
}

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

func boolConv(sval string) (interface{}, error) {
	return strconv.ParseBool(sval)
}

// GetBool return the bool value of a flag with the given name
func (f *FlagSet) GetBool(name string) (bool, error) {
	val, err := f.getFlagType(name, "bool", boolConv)
	if err != nil {
		return false, err
	}
	return val.(bool), nil
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string, validation ...func(value interface{}) error) {
	f.BoolVarP(p, name, "", value, usage, validation...)
}

// BoolVarP is like BoolVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolVarP(p *bool, name, shorthand string, value bool, usage string, validation ...func(value interface{}) error) {
	flag := f.VarPF(newBoolValue(value, p), name, shorthand, usage, validation...)
	flag.NoOptDefVal = "true"
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func BoolVar(p *bool, name string, value bool, usage string, validation ...func(value interface{}) error) {
	BoolVarP(p, name, "", value, usage, validation...)
}

// BoolVarP is like BoolVar, but accepts a shorthand letter that can be used after a single dash.
func BoolVarP(p *bool, name, shorthand string, value bool, usage string, validation ...func(value interface{}) error) {
	flag := CommandLine.VarPF(newBoolValue(value, p), name, shorthand, usage, validation...)
	flag.NoOptDefVal = "true"
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func (f *FlagSet) Bool(name string, value bool, usage string, validation ...func(value interface{}) error) *bool {
	return f.BoolP(name, "", value, usage, validation...)
}

// BoolP is like Bool, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BoolP(name, shorthand string, value bool, usage string, validation ...func(value interface{}) error) *bool {
	p := new(bool)
	f.BoolVarP(p, name, shorthand, value, usage, validation...)
	return p
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func Bool(name string, value bool, usage string, validation ...func(value interface{}) error) *bool {
	return BoolP(name, "", value, usage, validation...)
}

// BoolP is like Bool, but accepts a shorthand letter that can be used after a single dash.
func BoolP(name, shorthand string, value bool, usage string, validation ...func(value interface{}) error) *bool {
	b := CommandLine.BoolP(name, shorthand, value, usage, validation...)
	return b
}
