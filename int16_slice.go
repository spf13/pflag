package pflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- int16Slice Value
type int16SliceValue struct {
	value   *[]int16
	changed bool
}

func newInt16SliceValue(val []int16, p *[]int16) *int16SliceValue {
	isv := new(int16SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int16SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]int16, len(ss))
	for i, d := range ss {
		var err error
		var temp64 int64
		temp64, err = strconv.ParseInt(d, 0, 16)
		if err != nil {
			return err
		}
		out[i] = int16(temp64)

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *int16SliceValue) Type() string {
	return "int16Slice"
}

func (s *int16SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *int16SliceValue) fromString(val string) (int16, error) {
	t64, err := strconv.ParseInt(val, 0, 16)
	if err != nil {
		return 0, err
	}
	return int16(t64), nil
}

func (s *int16SliceValue) toString(val int16) string {
	return fmt.Sprintf("%d", val)
}

func (s *int16SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int16SliceValue) Replace(val []string) error {
	out := make([]int16, len(val))
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

func (s *int16SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func int16SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	// Empty string would cause a slice with one (empty) entry
	if len(val) == 0 {
		return []int16{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]int16, len(ss))
	for i, d := range ss {
		var err error
		var temp64 int64
		temp64, err = strconv.ParseInt(d, 0, 16)
		if err != nil {
			return nil, err
		}
		out[i] = int16(temp64)

	}
	return out, nil
}

// GetInt16Slice return the []int16 value of a flag with the given name
func (f *FlagSet) GetInt16Slice(name string) ([]int16, error) {
	val, err := f.getFlagType(name, "int16Slice", int16SliceConv)
	if err != nil {
		return []int16{}, err
	}
	return val.([]int16), nil
}

// Int16SliceVar defines a int16Slice flag with specified name, default value, and usage string.
// The argument p points to a []int16 variable in which to store the value of the flag.
func (f *FlagSet) Int16SliceVar(p *[]int16, name string, value []int16, usage string) {
	f.VarP(newInt16SliceValue(value, p), name, "", usage)
}

// Int16SliceVarP is like Int16SliceVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int16SliceVarP(p *[]int16, name, shorthand string, value []int16, usage string) {
	f.VarP(newInt16SliceValue(value, p), name, shorthand, usage)
}

// Int16SliceVar defines a int16[] flag with specified name, default value, and usage string.
// The argument p points to a int16[] variable in which to store the value of the flag.
func Int16SliceVar(p *[]int16, name string, value []int16, usage string) {
	CommandLine.VarP(newInt16SliceValue(value, p), name, "", usage)
}

// Int16SliceVarP is like Int16SliceVar, but accepts a shorthand letter that can be used after a single dash.
func Int16SliceVarP(p *[]int16, name, shorthand string, value []int16, usage string) {
	CommandLine.VarP(newInt16SliceValue(value, p), name, shorthand, usage)
}

// Int16Slice defines a []int16 flag with specified name, default value, and usage string.
// The return value is the address of a []int16 variable that stores the value of the flag.
func (f *FlagSet) Int16Slice(name string, value []int16, usage string) *[]int16 {
	p := []int16{}
	f.Int16SliceVarP(&p, name, "", value, usage)
	return &p
}

// Int16SliceP is like Int16Slice, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) Int16SliceP(name, shorthand string, value []int16, usage string) *[]int16 {
	p := []int16{}
	f.Int16SliceVarP(&p, name, shorthand, value, usage)
	return &p
}

// Int16Slice defines a []int16 flag with specified name, default value, and usage string.
// The return value is the address of a []int16 variable that stores the value of the flag.
func Int16Slice(name string, value []int16, usage string) *[]int16 {
	return CommandLine.Int16SliceP(name, "", value, usage)
}

// Int16SliceP is like Int16Slice, but accepts a shorthand letter that can be used after a single dash.
func Int16SliceP(name, shorthand string, value []int16, usage string) *[]int16 {
	return CommandLine.Int16SliceP(name, shorthand, value, usage)
}
