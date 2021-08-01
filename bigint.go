package pflag

import (
	"fmt"
	"math/big"
)

// -- big.Int Value
type bigIntValue big.Int

func newBigIntValue(val big.Int, p *big.Int) *bigIntValue {
	*p = val
	return (*bigIntValue)(p)
}

func (i *bigIntValue) Set(s string) error {
	v, ok := big.NewInt(0).SetString(s, 10)

	*i = bigIntValue(*v)

	if !ok {
		return fmt.Errorf("invalid integer %q", s)
	}

	return nil
}

func (i *bigIntValue) Get() interface{} {
	return big.Int(*i)
}

func (i *bigIntValue) Type() string {
	return "bigInt"
}

func (i *bigIntValue) String() string { return (*big.Int)(i).String() }

// GetBigInt return the big.Int value of a flag with the given name
func (f *FlagSet) GetBigInt(name string) (big.Int, error) {
	val, err := f.getFlagType(name, "bigInt")
	if err != nil {
		return *big.NewInt(0), err
	}
	return val.(big.Int), nil
}

// BigIntVar defines an big.Int flag with specified name, default value, and usage string.
// The argument p points to an big.Int variable in which to store the value of the flag.
func (f *FlagSet) BigIntVar(p *big.Int, name string, value big.Int, usage string) {
	f.VarP(newBigIntValue(value, p), name, "", usage)
}

// BigIntVarP is like BigIntVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BigIntVarP(p *big.Int, name, shorthand string, value big.Int, usage string) {
	f.VarP(newBigIntValue(value, p), name, shorthand, usage)
}

// BigIntVar defines an big.Int flag with specified name, default value, and usage string.
// The argument p points to an big.Int variable in which to store the value of the flag.
func BigIntVar(p *big.Int, name string, value big.Int, usage string) {
	CommandLine.VarP(newBigIntValue(value, p), name, "", usage)
}

// BigIntVarP is like BigIntVar, but accepts a shorthand letter that can be used after a single dash.
func BigIntVarP(p *big.Int, name, shorthand string, value big.Int, usage string) {
	CommandLine.VarP(newBigIntValue(value, p), name, shorthand, usage)
}

// BigInt defines an big.Int flag with specified name, default value, and usage string.
// The return value is the address of an big.Int variable that stores the value of the flag.
func (f *FlagSet) BigInt(name string, value big.Int, usage string) *big.Int {
	p := new(big.Int)
	f.BigIntVarP(p, name, "", value, usage)
	return p
}

// BigIntP is like BigInt, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BigIntP(name, shorthand string, value big.Int, usage string) *big.Int {
	p := new(big.Int)
	f.BigIntVarP(p, name, shorthand, value, usage)
	return p
}

// BigInt defines an big.Int flag with specified name, default value, and usage string.
// The return value is the address of an big.Int variable that stores the value of the flag.
func BigInt(name string, value big.Int, usage string) *big.Int {
	return CommandLine.BigIntP(name, "", value, usage)
}

// BigIntP is like BigInt, but accepts a shorthand letter that can be used after a single dash.
func BigIntP(name, shorthand string, value big.Int, usage string) *big.Int {
	return CommandLine.BigIntP(name, shorthand, value, usage)
}
