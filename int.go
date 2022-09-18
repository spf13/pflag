package pflag

import "strconv"

// -- int Value
type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = intValue(v)
	return err
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

func intConv(sval string) (interface{}, error) {
	return strconv.Atoi(sval)
}

// GetInt return the int value of a flag with the given name
func (f *FlagSet) GetInt(name string) (int, error) {
	val, err := f.getFlagType(name, "int", intConv)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.IntVarP(p, name, "", value, usage)
}

// IntVarN is like IntVarP, but adds the name as shorthand (non-posix).
func (f *FlagSet) IntVarN(p *int, name, shorthand string, value int, usage string) {
	f.VarN(newIntValue(value, p), name, shorthand, usage)
}

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntVarP(p *int, name, shorthand string, value int, usage string) {
	f.VarP(newIntValue(value, p), name, shorthand, usage)
}

// IntVarS is like IntVar, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) IntVarS(p *int, name, shorthand string, value int, usage string) {
	f.VarS(newIntValue(value, p), name, shorthand, usage)
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string) {
	CommandLine.IntVar(p, name, value, usage)
}

// IntVarN is like IntVarP, but adds the name as shorthand (non-posix).
func IntVarN(p *int, name, shorthand string, value int, usage string) {
	CommandLine.IntVarN(p, name, shorthand, value, usage)
}

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func IntVarP(p *int, name, shorthand string, value int, usage string) {
	CommandLine.IntVarP(p, name, shorthand, value, usage)
}

// IntVarS is like IntVar, but accepts a shorthand letter that can be used after a single dash, alone.
func IntVarS(p *int, name, shorthand string, value int, usage string) {
	CommandLine.IntVarS(p, name, shorthand, value, usage)
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (f *FlagSet) Int(name string, value int, usage string) *int {
	return f.IntP(name, "", value, usage)
}

// IntN is like IntP, but adds the name as shorthand (non-posix).
func (f *FlagSet) IntN(name, shorthand string, value int, usage string) *int {
	p := new(int)
	f.IntVarN(p, name, shorthand, value, usage)
	return p
}

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntP(name, shorthand string, value int, usage string) *int {
	p := new(int)
	f.IntVarP(p, name, shorthand, value, usage)
	return p
}

// IntS is like Int, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) IntS(name, shorthand string, value int, usage string) *int {
	p := new(int)
	f.IntVarS(p, name, shorthand, value, usage)
	return p
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string) *int {
	return CommandLine.Int(name, value, usage)
}

// IntN is like IntP, but adds the name as shorthand (non-posix).
func IntN(name, shorthand string, value int, usage string) *int {
	return CommandLine.IntN(name, shorthand, value, usage)
}

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func IntP(name, shorthand string, value int, usage string) *int {
	return CommandLine.IntP(name, shorthand, value, usage)
}

// IntS is like Int, but accepts a shorthand letter that can be used after a single dash, alone.
func IntS(name, shorthand string, value int, usage string) *int {
	return CommandLine.IntS(name, shorthand, value, usage)
}
