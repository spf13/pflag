package pflag

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// BytesHex adapts []byte for use as a flag. Value of flag is HEX encoded
type bytesHexValue []byte

// String implements pflag.Value.String.
func (bytesHex bytesHexValue) String() string {
	return fmt.Sprintf("%X", []byte(bytesHex))
}

// Set implements pflag.Value.Set.
func (bytesHex *bytesHexValue) Set(value string) error {
	bin, err := hex.DecodeString(strings.TrimSpace(value))

	if err != nil {
		return err
	}

	*bytesHex = bin

	return nil
}

// Type implements pflag.Value.Type.
func (*bytesHexValue) Type() string {
	return "bytesHex"
}

func newBytesHexValue(val []byte, p *[]byte) *bytesHexValue {
	*p = val
	return (*bytesHexValue)(p)
}

func bytesHexConv(sval string) (interface{}, error) {

	bin, err := hex.DecodeString(sval)

	if err == nil {
		return bin, nil
	}

	return nil, fmt.Errorf("invalid string being converted to Bytes: %s %s", sval, err)
}

// GetBytesHex return the []byte value of a flag with the given name
func (f *FlagSet) GetBytesHex(name string) ([]byte, error) {
	val, err := f.getFlagType(name, "bytesHex", bytesHexConv)

	if err != nil {
		return []byte{}, err
	}

	return val.([]byte), nil
}

// BytesHexVar defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func (f *FlagSet) BytesHexVar(p *[]byte, name string, value []byte, usage string) {
	f.BytesHexVarP(p, name, "", value, usage)
}

// BytesHexVarN is like BytesHexVarP, but adds the name as shorthand (non-posix).
func (f *FlagSet) BytesHexVarN(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarN(newBytesHexValue(value, p), name, shorthand, usage)
}

// BytesHexVarP is like BytesHexVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesHexVarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarP(newBytesHexValue(value, p), name, shorthand, usage)
}

// BytesHexVarS is like BytesHexVarP, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) BytesHexVarS(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarS(newBytesHexValue(value, p), name, shorthand, usage)
}

// BytesHexVar defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func BytesHexVar(p *[]byte, name string, value []byte, usage string) {
	CommandLine.BytesHexVar(p, name, value, usage)
}

// BytesHexVarN is like BytesHexVarP, but adds the name as shorthand (non-posix).
func BytesHexVarN(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.BytesHexVarN(p, name, shorthand, value, usage)
}

// BytesHexVarP is like BytesHexVar, but accepts a shorthand letter that can be used after a single dash.
func BytesHexVarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.BytesHexVarP(p, name, shorthand, value, usage)
}

// BytesHexVarS is like BytesHexVarP, but accepts a shorthand letter that can be used after a single dash, alone.
func BytesHexVarS(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.BytesHexVarS(p, name, shorthand, value, usage)
}

// BytesHex defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func (f *FlagSet) BytesHex(name string, value []byte, usage string) *[]byte {
	return f.BytesHexP(name, "", value, usage)
}

// BytesHexN is like BytesHexP, but adds the name as shorthand (non-posix).
func (f *FlagSet) BytesHexN(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesHexVarN(p, name, shorthand, value, usage)
	return p
}

// BytesHexP is like BytesHex, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesHexP(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesHexVarP(p, name, shorthand, value, usage)
	return p
}

// BytesHexS is like BytesHexP, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) BytesHexS(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesHexVarS(p, name, shorthand, value, usage)
	return p
}

// BytesHex defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func BytesHex(name string, value []byte, usage string) *[]byte {
	return CommandLine.BytesHexP(name, "", value, usage)
}

// BytesHexN is like BytesHexP, but adds the name as shorthand (non-posix).
func BytesHexN(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesHexN(name, shorthand, value, usage)
}

// BytesHexP is like BytesHex, but accepts a shorthand letter that can be used after a single dash.
func BytesHexP(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesHexP(name, shorthand, value, usage)
}

// BytesHexS is like BytesHexP, but accepts a shorthand letter that can be used after a single dash, alone.
func BytesHexS(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesHexS(name, shorthand, value, usage)
}

// BytesBase64 adapts []byte for use as a flag. Value of flag is Base64 encoded
type bytesBase64Value []byte

// String implements pflag.Value.String.
func (bytesBase64 bytesBase64Value) String() string {
	return base64.StdEncoding.EncodeToString([]byte(bytesBase64))
}

// Set implements pflag.Value.Set.
func (bytesBase64 *bytesBase64Value) Set(value string) error {
	bin, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))

	if err != nil {
		return err
	}

	*bytesBase64 = bin

	return nil
}

// Type implements pflag.Value.Type.
func (*bytesBase64Value) Type() string {
	return "bytesBase64"
}

func newBytesBase64Value(val []byte, p *[]byte) *bytesBase64Value {
	*p = val
	return (*bytesBase64Value)(p)
}

func bytesBase64ValueConv(sval string) (interface{}, error) {

	bin, err := base64.StdEncoding.DecodeString(sval)
	if err == nil {
		return bin, nil
	}

	return nil, fmt.Errorf("invalid string being converted to Bytes: %s %s", sval, err)
}

// GetBytesBase64 return the []byte value of a flag with the given name
func (f *FlagSet) GetBytesBase64(name string) ([]byte, error) {
	val, err := f.getFlagType(name, "bytesBase64", bytesBase64ValueConv)

	if err != nil {
		return []byte{}, err
	}

	return val.([]byte), nil
}

// BytesBase64Var defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func (f *FlagSet) BytesBase64Var(p *[]byte, name string, value []byte, usage string) {
	f.BytesBase64VarP(p, name, "", value, usage)
}

// BytesBase64VarN is like BytesBase64VarP, but adds the name as shorthand (non-posix).
func (f *FlagSet) BytesBase64VarN(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarN(newBytesBase64Value(value, p), name, shorthand, usage)
}

// BytesBase64VarP is like BytesBase64Var, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesBase64VarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarP(newBytesBase64Value(value, p), name, shorthand, usage)
}

// BytesBase64VarS is like BytesBase64Var, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) BytesBase64VarS(p *[]byte, name, shorthand string, value []byte, usage string) {
	f.VarS(newBytesBase64Value(value, p), name, shorthand, usage)
}

// BytesBase64Var defines an []byte flag with specified name, default value, and usage string.
// The argument p points to an []byte variable in which to store the value of the flag.
func BytesBase64Var(p *[]byte, name string, value []byte, usage string) {
	CommandLine.BytesBase64Var(p, name, value, usage)
}

// BytesBase64VarN is like BytesBase64VarP, but adds the name as shorthand (non-posix).
func BytesBase64VarN(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.BytesBase64VarN(p, name, shorthand, value, usage)
}

// BytesBase64VarP is like BytesBase64Var, but accepts a shorthand letter that can be used after a single dash.
func BytesBase64VarP(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.BytesBase64VarP(p, name, shorthand, value, usage)
}

// BytesBase64VarS is like BytesBase64Var, but accepts a shorthand letter that can be used after a single dash, alone.
func BytesBase64VarS(p *[]byte, name, shorthand string, value []byte, usage string) {
	CommandLine.BytesBase64VarS(p, name, shorthand, value, usage)
}

// BytesBase64 defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func (f *FlagSet) BytesBase64(name string, value []byte, usage string) *[]byte {
	return f.BytesBase64P(name, "", value, usage)
}

// BytesBase64N is like BytesBase64P, but adds the name as shorthand (non-posix).
func (f *FlagSet) BytesBase64N(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesBase64VarN(p, name, shorthand, value, usage)
	return p
}

// BytesBase64P is like BytesBase64, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesBase64VarP(p, name, shorthand, value, usage)
	return p
}

// BytesBase64S is like BytesBase64, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) BytesBase64S(name, shorthand string, value []byte, usage string) *[]byte {
	p := new([]byte)
	f.BytesBase64VarS(p, name, shorthand, value, usage)
	return p
}

// BytesBase64 defines an []byte flag with specified name, default value, and usage string.
// The return value is the address of an []byte variable that stores the value of the flag.
func BytesBase64(name string, value []byte, usage string) *[]byte {
	return CommandLine.BytesBase64P(name, "", value, usage)
}

// BytesBase64N is like BytesBase64P, but adds the name as shorthand (non-posix).
func BytesBase64N(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesBase64N(name, shorthand, value, usage)
}

// BytesBase64P is like BytesBase64, but accepts a shorthand letter that can be used after a single dash.
func BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesBase64P(name, shorthand, value, usage)
}

// BytesBase64S is like BytesBase64, but accepts a shorthand letter that can be used after a single dash, alone.
func BytesBase64S(name, shorthand string, value []byte, usage string) *[]byte {
	return CommandLine.BytesBase64S(name, shorthand, value, usage)
}
