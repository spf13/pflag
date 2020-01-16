package pflag

import (
	"reflect"
	"strconv"
)

// -- int64 Value
type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*i = int64Value(v)
	return err
}

func (i *int64Value) Type() string {
	return "int64"
}

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func int64Conv(sval string) (interface{}, error) {
	return strconv.ParseInt(sval, 0, 64)
}

// GetInt64 return the int64 value of a flag with the given name
func (f *FlagSet) GetInt64(name string) (int64, error) {
	val, err := f.getFlagType(name, "int64", int64Conv)
	if err != nil {
		return 0, err
	}
	return val.(int64), nil
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	f.VarP(newInt64Value(value, p), name, "", usage)
}

// Int64VarP is like Int64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64VarP(p *int64, name, shorthand string, value int64, usage string) {
	f.VarP(newInt64Value(value, p), name, shorthand, usage)
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func Int64Var(p *int64, name string, value int64, usage string) {
	CommandLine.VarP(newInt64Value(value, p), name, "", usage)
}

// Int64VarP is like Int64Var, but accepts a shorthand letter that can be used after a single dash.
func Int64VarP(p *int64, name, shorthand string, value int64, usage string) {
	CommandLine.VarP(newInt64Value(value, p), name, shorthand, usage)
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
	p := new(int64)
	f.Int64VarP(p, name, "", value, usage)
	return p
}

// Int64P is like Int64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int64P(name, shorthand string, value int64, usage string) *int64 {
	p := new(int64)
	f.Int64VarP(p, name, shorthand, value, usage)
	return p
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func Int64(name string, value int64, usage string) *int64 {
	return CommandLine.Int64P(name, "", value, usage)
}

// Int64P is like Int64, but accepts a shorthand letter that can be used after a single dash.
func Int64P(name, shorthand string, value int64, usage string) *int64 {
	return CommandLine.Int64P(name, shorthand, value, usage)
}

// Set int64 flag to flagSet
func setInt64Flag(flagSet *FlagSet, name, shorthand, value, usage string) error {
	defVal, err := int64Conv(value)
	if err != nil {
		return err
	}
	flagSet.Int64P(name, shorthand, defVal.(int64), usage)
	return nil
}

// Set int64 value from flagSet
func setInt64Value(flagSet *FlagSet, name string, fieldV reflect.Value) error {
	val, err := flagSet.GetInt64(name)
	if err != nil {
		return err
	}
	fieldV.Set(reflect.ValueOf(val))
	return nil
}
