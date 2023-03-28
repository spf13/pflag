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
func (f *FlagSet) IntVar(p *int, name string, value int, usage string, validation ...func(value int) error) {
	if len(validation) > 0 {
		validationFunc := interface{}(validation[0])
		f.VarP(newIntValue(value, p), name, "", usage, validationFunc)
		return
	}
	f.VarP(newIntValue(value, p), name, "", usage)
}

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntVarP(p *int, name, shorthand string, value int, usage string, validation ...func(value int) error) {
	if len(validation) > 0 {
		validationFunc := interface{}(validation[0])
		f.VarP(newIntValue(value, p), name, shorthand, usage, validationFunc)
		return
	}
	f.VarP(newIntValue(value, p), name, shorthand, usage)
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string, validation ...func(value int) error) {
	if len(validation) > 0 {
		validationFunc := interface{}(validation[0])
		CommandLine.VarP(newIntValue(value, p), name, "", usage, validationFunc)
		return
	}
	CommandLine.VarP(newIntValue(value, p), name, "", usage)
}

// IntVarP is like IntVar, but accepts a shorthand letter that can be used after a single dash.
func IntVarP(p *int, name, shorthand string, value int, usage string, validation ...func(value int) error) {
	if len(validation) > 0 {
		validationFunc := interface{}(validation[0])
		CommandLine.VarP(newIntValue(value, p), name, shorthand, usage, validationFunc)
		return
	}
	CommandLine.VarP(newIntValue(value, p), name, shorthand, usage)
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (f *FlagSet) Int(name string, value int, usage string, validation ...func(value int) error) *int {
	p := new(int)
	f.IntVarP(p, name, "", value, usage, validation...)
	return p
}

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IntP(name, shorthand string, value int, usage string, validation ...func(value int) error) *int {
	p := new(int)
	f.IntVarP(p, name, shorthand, value, usage, validation...)
	return p
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string, validation ...func(value int) error) *int {
	return CommandLine.IntP(name, "", value, usage, validation...)
}

// IntP is like Int, but accepts a shorthand letter that can be used after a single dash.
func IntP(name, shorthand string, value int, usage string, validation ...func(value int) error) *int {
	return CommandLine.IntP(name, shorthand, value, usage, validation...)
}
