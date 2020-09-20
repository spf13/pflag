package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- uint8Slice Value
type uint8SliceValue struct {
	value   *[]uint8
	changed bool
}

func newUint8SliceValue(val []uint8, p *[]uint8) *uint8SliceValue {
	isv := new(uint8SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint8SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]uint8, len(ss))
	for i, d := range ss {
		var err error
		var temp64 uint64
		temp64, err = strconv.ParseUint(d, 0, 8)
		if err != nil {
			return err
		}
		out[i] = uint8(temp64)

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *uint8SliceValue) Type() string {
	return "uint8Slice"
}

func (s *uint8SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *uint8SliceValue) fromString(val string) (uint8, error) {
	t64, err := strconv.ParseUint(val, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(t64), nil
}

func (s *uint8SliceValue) toString(val uint8) string {
	return fmt.Sprintf("%d", val)
}

func (s *uint8SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *uint8SliceValue) Replace(val []string) error {
	out := make([]uint8, len(val))
	for i, d := range val {
		var err error
		out[i], err = s.fromString(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *uint8SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func uint8SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []uint8{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]uint8, len(ss))
	for i, d := range ss {
		var err error
		var temp64 uint64
		temp64, err = strconv.ParseUint(d, 0, 8)
		if err != nil {
			return nil, err
		}
		out[i] = uint8(temp64)

	}
	return out, nil
}

// GetUint8Slice return the []uint8 value of a flag with the given name
func (f *FlagSet) GetUint8Slice(name string) ([]uint8, error) {
	val, err := f.getFlagType(name, "uint8Slice", uint8SliceConv)
	if err != nil {
		return []uint8{}, err
	}
	return val.([]uint8), nil
}

// Uint8SliceVar defines a uint8Slice flag with specified name, default value, and usage string.
// The argument p pouints to a []uint8 variable in which to store the value of the flag.
func (f *FlagSet) Uint8SliceVar(p *[]uint8, name string, value []uint8, usage string) {
	f.VarP(newUint8SliceValue(value, p), name, "", usage)
}

// Uint8SliceVarP is like Uint8SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Uint8SliceVarP(p *[]uint8, name, shorthand string, value []uint8, usage string) {
	f.VarP(newUint8SliceValue(value, p), name, shorthand, usage)
}

// Uint8SliceVar defines a uint8[] flag with specified name, default value, and usage string.
// The argument p pouints to a uint8[] variable in which to store the value of the flag.
func Uint8SliceVar(p *[]uint8, name string, value []uint8, usage string) {
	CommandLine.VarP(newUint8SliceValue(value, p), name, "", usage)
}

// Uint8SliceVarP is like Uint8SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Uint8SliceVarP(p *[]uint8, name, shorthand string, value []uint8, usage string) {
	CommandLine.VarP(newUint8SliceValue(value, p), name, shorthand, usage)
}

// Uint8Slice defines a []uint8 flag with specified name, default value, and usage string.
// The return value is the address of a []uint8 variable that stores the value of the flag.
func (f *FlagSet) Uint8Slice(name string, value []uint8, usage string) *[]uint8 {
	p := []uint8{}
	f.Uint8SliceVarP(&p, name, "", value, usage)
	return &p
}

// Uint8SliceP is like Uint8Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Uint8SliceP(name, shorthand string, value []uint8, usage string) *[]uint8 {
	p := []uint8{}
	f.Uint8SliceVarP(&p, name, shorthand, value, usage)
	return &p
}

// Uint8Slice defines a []uint8 flag with specified name, default value, and usage string.
// The return value is the address of a []uint8 variable that stores the value of the flag.
func Uint8Slice(name string, value []uint8, usage string) *[]uint8 {
	return CommandLine.Uint8SliceP(name, "", value, usage)
}

// Uint8SliceP is like Uint8Slice, but accepts a shorthand letter that can be used after a single dash.
func Uint8SliceP(name, shorthand string, value []uint8, usage string) *[]uint8 {
	return CommandLine.Uint8SliceP(name, shorthand, value, usage)
}
