// +build go1.15

package pflag

import "strconv"

// -- complex128 Value
type complex128Value complex128

func newComplex128Value(val complex128, p *complex128) *complex128Value {
	*p = val
	return (*complex128Value)(p)
}

func (c *complex128Value) Set(s string) error {
	v, err := strconv.ParseComplex(s, 128)
	*c = complex128Value(v)
	return err
}

func (c *complex128Value) Type() string {
	return "complex128"
}

func (c *complex128Value) String() string { return strconv.FormatComplex(complex128(*c), 'g', -1, 128) }

func complex128Conv(sval string) (interface{}, error) {
	return strconv.ParseComplex(sval, 128)
}

// GetComplex128 return the complex128 value of a flag with the given name
func (f *FlagSet) GetComplex128(name string) (complex128, error) {
	val, err := f.getFlagType(name, "complex128", complex128Conv)
	if err != nil {
		return 0, err
	}
	return val.(complex128), nil
}

// Complex128Var defines a complex128 flag with specified name, default value, and usage string.
// The argument p points to a complex128 variable in which to store the value of the flag.
func (f *FlagSet) Complex128Var(p *complex128, name string, value complex128, usage string) {
	f.VarP(newComplex128Value(value, p), name, "", usage)
}

// Complex128VarP is like Complex128Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Complex128VarP(p *complex128, name, shorthand string, value complex128, usage string) {
	f.VarP(newComplex128Value(value, p), name, shorthand, usage)
}

// Complex128Var defines a complex128 flag with specified name, default value, and usage string.
// The argument p points to a complex128 variable in which to store the value of the flag.
func Complex128Var(p *complex128, name string, value complex128, usage string) {
	CommandLine.VarP(newComplex128Value(value, p), name, "", usage)
}

// Complex128VarP is like Complex128Var, but accepts a shorthand letter that can be used after a single dash.
func Complex128VarP(p *complex128, name, shorthand string, value complex128, usage string) {
	CommandLine.VarP(newComplex128Value(value, p), name, shorthand, usage)
}

// Complex128 defines a complex128 flag with specified name, default value, and usage string.
// The return value is the address of a complex128 variable that stores the value of the flag.
func (f *FlagSet) Complex128(name string, value complex128, usage string) *complex128 {
	p := new(complex128)
	f.Complex128VarP(p, name, "", value, usage)
	return p
}

// Complex128P is like Complex128, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Complex128P(name, shorthand string, value complex128, usage string) *complex128 {
	p := new(complex128)
	f.Complex128VarP(p, name, shorthand, value, usage)
	return p
}

// Complex128 defines a complex128 flag with specified name, default value, and usage string.
// The return value is the address of a complex128 variable that stores the value of the flag.
func Complex128(name string, value complex128, usage string) *complex128 {
	return CommandLine.Complex128P(name, "", value, usage)
}

// Complex128P is like Complex128, but accepts a shorthand letter that can be used after a single dash.
func Complex128P(name, shorthand string, value complex128, usage string) *complex128 {
	return CommandLine.Complex128P(name, shorthand, value, usage)
}
