package pflag

import (
	"fmt"
	"net"
	"strings"
)

// IPNet adapts net.IPNet for use as a flag.
type ipNetValue net.IPNet

func (ipnet ipNetValue) String() string {
	n := net.IPNet(ipnet)
	return n.String()
}

func (ipnet *ipNetValue) Set(value string) error {
	_, n, err := net.ParseCIDR(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	*ipnet = ipNetValue(*n)
	return nil
}

func (*ipNetValue) Type() string {
	return "ipNet"
}

func newIPNetValue(val net.IPNet, p *net.IPNet) *ipNetValue {
	*p = val
	return (*ipNetValue)(p)
}

func ipNetConv(sval string) (interface{}, error) {
	_, n, err := net.ParseCIDR(strings.TrimSpace(sval))
	if err == nil {
		return *n, nil
	}
	return nil, fmt.Errorf("invalid string being converted to IPNet: %s", sval)
}

// GetIPNet return the net.IPNet value of a flag with the given name
func (f *FlagSet) GetIPNet(name string) (net.IPNet, error) {
	val, err := f.getFlagType(name, "ipNet", ipNetConv)
	if err != nil {
		return net.IPNet{}, err
	}
	return val.(net.IPNet), nil
}

// IPNetVar defines an net.IPNet flag with specified name, default value, and usage string.
// The argument p points to an net.IPNet variable in which to store the value of the flag.
func (f *FlagSet) IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string) {
	f.IPNetVarP(p, name, "", value, usage)
}

// IPNetVarN is like IPNetVarP, but adds the name as shorthand (non-posix).
func (f *FlagSet) IPNetVarN(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	f.VarN(newIPNetValue(value, p), name, shorthand, usage)
}

// IPNetVarP is like IPNetVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPNetVarP(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	f.VarP(newIPNetValue(value, p), name, shorthand, usage)
}

// IPNetVarS is like IPNetVar, but accepts a shorthand letter that can be used after a single , alone, alone.
func (f *FlagSet) IPNetVarS(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	f.VarS(newIPNetValue(value, p), name, shorthand, usage)
}

// IPNetVar defines an net.IPNet flag with specified name, default value, and usage string.
// The argument p points to an net.IPNet variable in which to store the value of the flag.
func IPNetVar(p *net.IPNet, name string, value net.IPNet, usage string) {
	CommandLine.IPNetVar(p, name, value, usage)
}

// IPNetVarN is like IPNetVarP, but adds the name as shorthand (non-posix).
func IPNetVarN(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	CommandLine.IPNetVarN(p, name, shorthand, value, usage)
}

// IPNetVarP is like IPNetVar, but accepts a shorthand letter that can be used after a single dash.
func IPNetVarP(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	CommandLine.IPNetVarP(p, name, shorthand, value, usage)
}

// IPNetVarS is like IPNetVar, but accepts a shorthand letter that can be used after a single dash, alone.
func IPNetVarS(p *net.IPNet, name, shorthand string, value net.IPNet, usage string) {
	CommandLine.IPNetVarS(p, name, shorthand, value, usage)
}

// IPNet defines an net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of an net.IPNet variable that stores the value of the flag.
func (f *FlagSet) IPNet(name string, value net.IPNet, usage string) *net.IPNet {
	return f.IPNetP(name, "", value, usage)
}

// IPNetN is like IPNetP, but adds the name as shorthand (non-posix).
func (f *FlagSet) IPNetN(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	p := new(net.IPNet)
	f.IPNetVarN(p, name, shorthand, value, usage)
	return p
}

// IPNetP is like IPNet, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	p := new(net.IPNet)
	f.IPNetVarP(p, name, shorthand, value, usage)
	return p
}

// IPNetS is like IPNet, but accepts a shorthand letter that can be used after a single dash, alone.
func (f *FlagSet) IPNetS(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	p := new(net.IPNet)
	f.IPNetVarS(p, name, shorthand, value, usage)
	return p
}

// IPNet defines an net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of an net.IPNet variable that stores the value of the flag.
func IPNet(name string, value net.IPNet, usage string) *net.IPNet {
	return CommandLine.IPNet(name, value, usage)
}

// IPNetN is like IPNetP, but adds the name as shorthand (non-posix).
func IPNetN(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	return CommandLine.IPNetN(name, shorthand, value, usage)
}

// IPNetP is like IPNet, but accepts a shorthand letter that can be used after a single dash.
func IPNetP(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	return CommandLine.IPNetP(name, shorthand, value, usage)
}

// IPNetS is like IPNet, but accepts a shorthand letter that can be used after a single dash, alone.
func IPNetS(name, shorthand string, value net.IPNet, usage string) *net.IPNet {
	return CommandLine.IPNetS(name, shorthand, value, usage)
}
