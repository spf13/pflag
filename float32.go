package pflag

import "strconv"

// -- float32 Value
type float32Value float32

func newFloat32Value(val float32, p *float32) *float32Value {
	*p = val
	return (*float32Value)(p)
}

func (f *float32Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	*f = float32Value(v)
	return err
}

func (f *float32Value) Type() string {
	return "float32"
}

func (f *float32Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

func float32Conv(sval string) (interface{}, error) {
	v, err := strconv.ParseFloat(sval, 32)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

// GetFloat32 return the float32 value of a flag with the given name
func (f *FlagSet) GetFloat32(name string) (float32, error) {
	val, err := f.getFlagType(name, "float32", float32Conv)
	if err != nil {
		return 0, err
	}
	return val.(float32), nil
}

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func (f *FlagSet) Float32Var(p *float32, name string, value float32, usage string) {
	f.Float32VarP(p, name, "", value, usage)
}

// Float32VarN is like Float32VarP, but adds the name as shorthand (non-posix).
func (f *FlagSet) Float32VarN(p *float32, name, shorthand string, value float32, usage string) {
	f.VarN(newFloat32Value(value, p), name, shorthand, usage)
}

// Float32VarP is like Float32Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32VarP(p *float32, name, shorthand string, value float32, usage string) {
	f.VarP(newFloat32Value(value, p), name, shorthand, usage)
}

// Float32VarS is like Float32Var, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) Float32VarS(p *float32, name, shorthand string, value float32, usage string) {
	f.VarS(newFloat32Value(value, p), name, shorthand, usage)
}

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func Float32Var(p *float32, name string, value float32, usage string) {
	CommandLine.Float32Var(p, name, value, usage)
}

// Float32VarN is like Float32VarP, but adds the name as shorthand (non-posix).
func Float32VarN(p *float32, name, shorthand string, value float32, usage string) {
	CommandLine.Float32VarN(p, name, shorthand, value, usage)
}

// Float32VarP is like Float32Var, but accepts a shorthand letter that can be used after a single dash.
func Float32VarP(p *float32, name, shorthand string, value float32, usage string) {
	CommandLine.Float32VarP(p, name, shorthand, value, usage)
}

// Float32VarS is like Float32Var, but accepts a shorthand letter that can be used after a single dash, alone.
func Float32VarS(p *float32, name, shorthand string, value float32, usage string) {
	CommandLine.Float32VarS(p, name, shorthand, value, usage)
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func (f *FlagSet) Float32(name string, value float32, usage string) *float32 {
	return f.Float32P(name, "", value, usage)
}

// Float32N is like Float32P, but adds the name as shorthand (non-posix).
func (f *FlagSet) Float32N(name, shorthand string, value float32, usage string) *float32 {
	p := new(float32)
	f.Float32VarN(p, name, shorthand, value, usage)
	return p
}

// Float32P is like Float32, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Float32P(name, shorthand string, value float32, usage string) *float32 {
	p := new(float32)
	f.Float32VarP(p, name, shorthand, value, usage)
	return p
}

// Float32S is like Float32, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) Float32S(name, shorthand string, value float32, usage string) *float32 {
	p := new(float32)
	f.Float32VarS(p, name, shorthand, value, usage)
	return p
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func Float32(name string, value float32, usage string) *float32 {
	return CommandLine.Float32(name, value, usage)
}

// Float32N is like Float32P, but adds the name as shorthand (non-posix).
func Float32N(name, shorthand string, value float32, usage string) *float32 {
	return CommandLine.Float32N(name, shorthand, value, usage)
}

// Float32P is like Float32, but accepts a shorthand letter that can be used after a single dash.
func Float32P(name, shorthand string, value float32, usage string) *float32 {
	return CommandLine.Float32P(name, shorthand, value, usage)
}

// Float32S is like Float32, but accepts a shorthand letter that can be used after a single dash, alone.
func Float32S(name, shorthand string, value float32, usage string) *float32 {
	return CommandLine.Float32S(name, shorthand, value, usage)
}
