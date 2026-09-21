package pflag

import "github.com/google/uuid"

// optional interface to indicate uuidean flags that can be
// supplied without "=value" text
type uuidFlag interface {
	Value
	IsUUIDFlag() uuid.UUID
}

// -- uuid Value
type uuidValue uuid.UUID

func newUUIDValue(val uuid.UUID, p *uuid.UUID) *uuidValue {
	*p = val
	return (*uuidValue)(p)
}

func (u *uuidValue) Set(s string) error {
	v, err := uuid.Parse(s)
	*u = uuidValue(v)
	return err
}

func (u *uuidValue) Type() string {
	return "uuid"
}

func (u *uuidValue) String() string { return uuid.UUID(*u).String() }

func (u *uuidValue) IsUUIDFlag() bool { return true }

func uuidConv(sval string) (interface{}, error) {
	return uuid.Parse(sval)
}

// GetUUID return the uuid value of a flag with the given name
func (f *FlagSet) GetUUID(name string) (uuid.UUID, error) {
	val, err := f.getFlagType(name, "uuid", uuidConv)
	if err != nil {
		return uuid.New(), err
	}
	return val.(uuid.UUID), nil
}

// UUIDVar defines a uuid flag with specified name, default value, and usage string.
// The argument p points to a uuid variable in which to store the value of the flag.
func (f *FlagSet) UUIDVar(p *uuid.UUID, name string, value uuid.UUID, usage string) {
	f.UUIDVarP(p, name, "", value, usage)
}

// UUIDVarP is like uuidVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UUIDVarP(p *uuid.UUID, name, shorthand string, value uuid.UUID, usage string) {
	flag := f.VarPF(newUUIDValue(value, p), name, shorthand, usage)
	flag.NoOptDefVal = "true"
}

// UUIDVar defines a uuid flag with specified name, default value, and usage string.
// The argument p points to a uuid variable in which to store the value of the flag.
func UUIDVar(p *uuid.UUID, name string, value uuid.UUID, usage string) {
	UUIDVarP(p, name, "", value, usage)
}

// UUIDVarP is like uuidVar, but accepts a shorthand letter that can be used after a single dash.
func UUIDVarP(p *uuid.UUID, name, shorthand string, value uuid.UUID, usage string) {
	flag := CommandLine.VarPF(newUUIDValue(value, p), name, shorthand, usage)
	flag.NoOptDefVal = "true"
}

// UUID defines a uuid flag with specified name, default value, and usage string.
// The return value is the address of a uuid variable that stores the value of the flag.
func (f *FlagSet) UUID(name string, value uuid.UUID, usage string) *uuid.UUID {
	return f.UUIDP(name, "", value, usage)
}

// UUIDP is like uuid, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) UUIDP(name, shorthand string, value uuid.UUID, usage string) *uuid.UUID {
	p := uuid.New()
	f.UUIDVarP(&p, name, shorthand, value, usage)
	return &p
}

// UUID defines a uuid flag with specified name, default value, and usage string.
// The return value is the address of a uuid variable that stores the value of the flag.
func UUID(name string, value uuid.UUID, usage string) *uuid.UUID {
	return UUIDP(name, "", value, usage)
}

// UUIDP is like UUID, but accepts a shorthand letter that can be used after a single dash.
func UUIDP(name, shorthand string, value uuid.UUID, usage string) *uuid.UUID {
	b := CommandLine.UUIDP(name, shorthand, value, usage)
	return b
}
