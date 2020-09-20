// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- int8Slice Value
type int8SliceValue struct {
	value   *[]int8
	changed bool
}

func newInt8SliceValue(val []int8, p *[]int8) *int8SliceValue {
	isv := new(int8SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int8SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]int8, len(ss))
	for i, d := range ss {
		var err error
		var temp64 int64
		temp64, err = strconv.ParseInt(d, 0, 8)
		if err != nil {
			return err
		}
		out[i] = int8(temp64)

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *int8SliceValue) Type() string {
	return "int8Slice"
}

func (s *int8SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *int8SliceValue) fromString(val string) (int8, error) {
	t64, err := strconv.ParseInt(val, 0, 8)
	if err != nil {
		return 0, err
	}
	return int8(t64), nil
}

func (s *int8SliceValue) toString(val int8) string {
	return fmt.Sprintf("%d", val)
}

func (s *int8SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int8SliceValue) Replace(val []string) error {
	out := make([]int8, len(val))
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

func (s *int8SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func int8SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []int8{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]int8, len(ss))
	for i, d := range ss {
		var err error
		var temp64 int64
		temp64, err = strconv.ParseInt(d, 0, 8)
		if err != nil {
			return nil, err
		}
		out[i] = int8(temp64)

	}
	return out, nil
}

// GetInt8Slice return the []int8 value of a flag with the given name
func (f *FlagSet) GetInt8Slice(name string) ([]int8, error) {
	val, err := f.getFlagType(name, "int8Slice", int8SliceConv)
	if err != nil {
		return []int8{}, err
	}
	return val.([]int8), nil
}

// Int8SliceVar defines a int8Slice flag with specified name, default value, and usage string.
// The argument p points to a []int8 variable in which to store the value of the flag.
func (f *FlagSet) Int8SliceVar(p *[]int8, name string, value []int8, usage string) {
	f.VarP(newInt8SliceValue(value, p), name, "", usage)
}

// Int8SliceVarP is like Int8SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int8SliceVarP(p *[]int8, name, shorthand string, value []int8, usage string) {
	f.VarP(newInt8SliceValue(value, p), name, shorthand, usage)
}

// Int8SliceVar defines a int8[] flag with specified name, default value, and usage string.
// The argument p points to a int8[] variable in which to store the value of the flag.
func Int8SliceVar(p *[]int8, name string, value []int8, usage string) {
	CommandLine.VarP(newInt8SliceValue(value, p), name, "", usage)
}

// Int8SliceVarP is like Int8SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Int8SliceVarP(p *[]int8, name, shorthand string, value []int8, usage string) {
	CommandLine.VarP(newInt8SliceValue(value, p), name, shorthand, usage)
}

// Int8Slice defines a []int8 flag with specified name, default value, and usage string.
// The return value is the address of a []int8 variable that stores the value of the flag.
func (f *FlagSet) Int8Slice(name string, value []int8, usage string) *[]int8 {
	p := []int8{}
	f.Int8SliceVarP(&p, name, "", value, usage)
	return &p
}

// Int8SliceP is like Int8Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int8SliceP(name, shorthand string, value []int8, usage string) *[]int8 {
	p := []int8{}
	f.Int8SliceVarP(&p, name, shorthand, value, usage)
	return &p
}

// Int8Slice defines a []int8 flag with specified name, default value, and usage string.
// The return value is the address of a []int8 variable that stores the value of the flag.
func Int8Slice(name string, value []int8, usage string) *[]int8 {
	return CommandLine.Int8SliceP(name, "", value, usage)
}

// Int8SliceP is like Int8Slice, but accepts a shorthand letter that can be used after a single dash.
func Int8SliceP(name, shorthand string, value []int8, usage string) *[]int8 {
	return CommandLine.Int8SliceP(name, shorthand, value, usage)
}
