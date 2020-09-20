// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- uint32Slice Value
type uint32SliceValue struct {
	value   *[]uint32
	changed bool
}

func newUint32SliceValue(val []uint32, p *[]uint32) *uint32SliceValue {
	isv := new(uint32SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint32SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]uint32, len(ss))
	for i, d := range ss {
		var err error
		var temp64 uint64
		temp64, err = strconv.ParseUint(d, 0, 32)
		if err != nil {
			return err
		}
		out[i] = uint32(temp64)

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *uint32SliceValue) Type() string {
	return "uint32Slice"
}

func (s *uint32SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *uint32SliceValue) fromString(val string) (uint32, error) {
	t64, err := strconv.ParseUint(val, 0, 32)
	if err != nil {
		return 0, err
	}
	return uint32(t64), nil
}

func (s *uint32SliceValue) toString(val uint32) string {
	return fmt.Sprintf("%d", val)
}

func (s *uint32SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *uint32SliceValue) Replace(val []string) error {
	out := make([]uint32, len(val))
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

func (s *uint32SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func uint32SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []uint32{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]uint32, len(ss))
	for i, d := range ss {
		var err error
		var temp64 uint64
		temp64, err = strconv.ParseUint(d, 0, 32)
		if err != nil {
			return nil, err
		}
		out[i] = uint32(temp64)

	}
	return out, nil
}

// GetUint32Slice return the []uint32 value of a flag with the given name
func (f *FlagSet) GetUint32Slice(name string) ([]uint32, error) {
	val, err := f.getFlagType(name, "uint32Slice", uint32SliceConv)
	if err != nil {
		return []uint32{}, err
	}
	return val.([]uint32), nil
}

// Uint32SliceVar defines a uint32Slice flag with specified name, default value, and usage string.
// The argument p pouints to a []uint32 variable in which to store the value of the flag.
func (f *FlagSet) Uint32SliceVar(p *[]uint32, name string, value []uint32, usage string) {
	f.VarP(newUint32SliceValue(value, p), name, "", usage)
}

// Uint32SliceVarP is like Uint32SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Uint32SliceVarP(p *[]uint32, name, shorthand string, value []uint32, usage string) {
	f.VarP(newUint32SliceValue(value, p), name, shorthand, usage)
}

// Uint32SliceVar defines a uint32[] flag with specified name, default value, and usage string.
// The argument p pouints to a uint32[] variable in which to store the value of the flag.
func Uint32SliceVar(p *[]uint32, name string, value []uint32, usage string) {
	CommandLine.VarP(newUint32SliceValue(value, p), name, "", usage)
}

// Uint32SliceVarP is like Uint32SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Uint32SliceVarP(p *[]uint32, name, shorthand string, value []uint32, usage string) {
	CommandLine.VarP(newUint32SliceValue(value, p), name, shorthand, usage)
}

// Uint32Slice defines a []uint32 flag with specified name, default value, and usage string.
// The return value is the address of a []uint32 variable that stores the value of the flag.
func (f *FlagSet) Uint32Slice(name string, value []uint32, usage string) *[]uint32 {
	p := []uint32{}
	f.Uint32SliceVarP(&p, name, "", value, usage)
	return &p
}

// Uint32SliceP is like Uint32Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Uint32SliceP(name, shorthand string, value []uint32, usage string) *[]uint32 {
	p := []uint32{}
	f.Uint32SliceVarP(&p, name, shorthand, value, usage)
	return &p
}

// Uint32Slice defines a []uint32 flag with specified name, default value, and usage string.
// The return value is the address of a []uint32 variable that stores the value of the flag.
func Uint32Slice(name string, value []uint32, usage string) *[]uint32 {
	return CommandLine.Uint32SliceP(name, "", value, usage)
}

// Uint32SliceP is like Uint32Slice, but accepts a shorthand letter that can be used after a single dash.
func Uint32SliceP(name, shorthand string, value []uint32, usage string) *[]uint32 {
	return CommandLine.Uint32SliceP(name, shorthand, value, usage)
}
