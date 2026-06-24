// +build go1.15

package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- complex128Slice Value
type complex128SliceValue struct {
	value   *[]complex128
	changed bool
}

func newComplex128SliceValue(val []complex128, p *[]complex128) *complex128SliceValue {
	isv := new(complex128SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *complex128SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]complex128, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseComplex(d, 128)
		if err != nil {
			return err
		}

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *complex128SliceValue) Type() string {
	return "complex128Slice"
}

func (s *complex128SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%f", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *complex128SliceValue) fromString(val string) (complex128, error) {
	return strconv.ParseComplex(val, 128)
}

func (s *complex128SliceValue) toString(val complex128) string {
	return fmt.Sprintf("%f", val)
}

func (s *complex128SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *complex128SliceValue) Replace(val []string) error {
	out := make([]complex128, len(val))
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

func (s *complex128SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func complex128SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []complex128{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]complex128, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseComplex(d, 128)
		if err != nil {
			return nil, err
		}

	}
	return out, nil
}

// GetComplex128Slice return the []complex128 value of a flag with the given name
func (f *FlagSet) GetComplex128Slice(name string) ([]complex128, error) {
	val, err := f.getFlagType(name, "complex128Slice", complex128SliceConv)
	if err != nil {
		return []complex128{}, err
	}
	return val.([]complex128), nil
}

// Complex128SliceVar defines a complex128Slice flag with specified name, default value, and usage string.
// The argument p points to a []complex128 variable in which to store the value of the flag.
func (f *FlagSet) Complex128SliceVar(p *[]complex128, name string, value []complex128, usage string) {
	f.VarP(newComplex128SliceValue(value, p), name, "", usage)
}

// Complex128SliceVarP is like Complex128SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Complex128SliceVarP(p *[]complex128, name, shorthand string, value []complex128, usage string) {
	f.VarP(newComplex128SliceValue(value, p), name, shorthand, usage)
}

// Complex128SliceVar defines a complex128[] flag with specified name, default value, and usage string.
// The argument p points to a complex128[] variable in which to store the value of the flag.
func Complex128SliceVar(p *[]complex128, name string, value []complex128, usage string) {
	CommandLine.VarP(newComplex128SliceValue(value, p), name, "", usage)
}

// Complex128SliceVarP is like Complex128SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Complex128SliceVarP(p *[]complex128, name, shorthand string, value []complex128, usage string) {
	CommandLine.VarP(newComplex128SliceValue(value, p), name, shorthand, usage)
}

// Complex128Slice defines a []complex128 flag with specified name, default value, and usage string.
// The return value is the address of a []complex128 variable that stores the value of the flag.
func (f *FlagSet) Complex128Slice(name string, value []complex128, usage string) *[]complex128 {
	p := []complex128{}
	f.Complex128SliceVarP(&p, name, "", value, usage)
	return &p
}

// Complex128SliceP is like Complex128Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Complex128SliceP(name, shorthand string, value []complex128, usage string) *[]complex128 {
	p := []complex128{}
	f.Complex128SliceVarP(&p, name, shorthand, value, usage)
	return &p
}

// Complex128Slice defines a []complex128 flag with specified name, default value, and usage string.
// The return value is the address of a []complex128 variable that stores the value of the flag.
func Complex128Slice(name string, value []complex128, usage string) *[]complex128 {
	return CommandLine.Complex128SliceP(name, "", value, usage)
}

// Complex128SliceP is like Complex128Slice, but accepts a shorthand letter that can be used after a single dash.
func Complex128SliceP(name, shorthand string, value []complex128, usage string) *[]complex128 {
	return CommandLine.Complex128SliceP(name, shorthand, value, usage)
}
