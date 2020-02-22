package pflag

import (
	"github.com/docker/go-units"
)

const byteSizeFlagType = "byte-size"

// byteSizeValue used to pass byte sizes to a go-flags CLI
type byteSizeValue uint64

func newByteSizeValue(val uint64, p *uint64) *byteSizeValue {
	*p = val
	return (*byteSizeValue)(p)
}

// MarshalFlag implements go-flags Marshaller interface
func (b byteSizeValue) MarshalFlag() (string, error) {
	return units.HumanSize(float64(b)), nil
}

// UnmarshalFlag implements go-flags Unmarshaller interface
func (b *byteSizeValue) UnmarshalFlag(value string) error {
	sz, err := units.FromHumanSize(value)
	if err != nil {
		return err
	}
	*b = byteSizeValue(uint64(sz))
	return nil
}

// String method for a bytesize (pflag value and stringer interface)
func (b byteSizeValue) String() string {
	return units.HumanSize(float64(b))
}

// Set the value of this bytesize (pflag value interfaces)
func (b *byteSizeValue) Set(value string) error {
	return b.UnmarshalFlag(value)
}

// Type returns the type of the pflag value (pflag value interface)
func (b *byteSizeValue) Type() string {
	return byteSizeFlagType
}

func byteSizeConv(sval string) (interface{}, error) {
	var b byteSizeValue
	err := b.UnmarshalFlag(sval)
	return uint64(b), err
}

// GetByteSize return the ByteSize value of a flag with the given name
func (f *FlagSet) GetByteSize(name string) (uint64, error) {
	val, err := f.getFlagType(name, byteSizeFlagType, byteSizeConv)
	if err != nil {
		return 0, err
	}
	return val.(uint64), nil
}

// ByteSizeVar defines an uint64 flag with specified name, default value, and usage string.
// The argument p pouint64s to an uint64 variable in which to store the value of the flag.
func (f *FlagSet) ByteSizeVar(p *uint64, name string, value uint64, usage string) {
	f.VarP(newByteSizeValue(value, p), name, "", usage)
}

// ByteSizeVarP is like ByteSizeVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) ByteSizeVarP(p *uint64, name, shorthand string, value uint64, usage string) {
	f.VarP(newByteSizeValue(value, p), name, shorthand, usage)
}

// ByteSizeVar defines an uint64 flag with specified name, default value, and usage string.
// The argument p pouint64s to an uint64 variable in which to store the value of the flag.
func ByteSizeVar(p *uint64, name string, value uint64, usage string) {
	CommandLine.VarP(newByteSizeValue(value, p), name, "", usage)
}

// ByteSizeVarP is like ByteSizeVar, but accepts a shorthand letter that can be used after a single dash.
func ByteSizeVarP(p *uint64, name, shorthand string, value uint64, usage string) {
	CommandLine.VarP(newByteSizeValue(value, p), name, shorthand, usage)
}

// ByteSize defines an uint64 flag with specified name, default value, and usage string.
// The return value is the address of an uint64 variable that stores the value of the flag.
func (f *FlagSet) ByteSize(name string, value uint64, usage string) *uint64 {
	p := new(uint64)
	f.ByteSizeVarP(p, name, "", value, usage)
	return p
}

// ByteSizeP is like ByteSize, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) ByteSizeP(name, shorthand string, value uint64, usage string) *uint64 {
	p := new(uint64)
	f.ByteSizeVarP(p, name, shorthand, value, usage)
	return p
}

// ByteSize defines an uint64 flag with specified name, default value, and usage string.
// The return value is the address of an uint64 variable that stores the value of the flag.
func ByteSize(name string, value uint64, usage string) *uint64 {
	return CommandLine.ByteSizeP(name, "", value, usage)
}

// ByteSizeP is like ByteSize, but accepts a shorthand letter that can be used after a single dash.
func ByteSizeP(name, shorthand string, value uint64, usage string) *uint64 {
	return CommandLine.ByteSizeP(name, shorthand, value, usage)
}
