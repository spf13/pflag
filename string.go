package pflag

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) String() string { return string(*s) }

func stringConv(sval string) (interface{}, error) {
	return sval, nil
}

// GetString return the string value of a flag with the given name
func (f *FlagSet) GetString(name string) (string, error) {
	val, err := f.getFlagType(name, "string", stringConv)
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (f *FlagSet) StringVar(p *string, name string, value string, usage string, opts ...FlagOption) {
	f.Var(newStringValue(value, p), name, usage, opts...)
}

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
//
// Deprecated: Use StringVar with the WithShorthand option instead
func (f *FlagSet) StringVarP(p *string, name, shorthand string, value string, usage string) {
	f.Var(newStringValue(value, p), name, usage, WithShorthand(shorthand))
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string, opts ...FlagOption) {
	CommandLine.Var(newStringValue(value, p), name, usage, opts...)
}

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
//
// Deprecated: Use StringVar with the WithShorthand option instead
func StringVarP(p *string, name, shorthand string, value string, usage string) {
	CommandLine.Var(newStringValue(value, p), name, usage, WithShorthand(shorthand))
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (f *FlagSet) String(name string, value string, usage string, opts ...FlagOption) *string {
	p := new(string)
	f.StringVar(p, name, value, usage, opts...)
	return p
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
//
// Deprecated: Use String with the WithShorthand option instead
func (f *FlagSet) StringP(name, shorthand string, value string, usage string) *string {
	p := new(string)
	f.StringVarP(p, name, shorthand, value, usage)
	return p
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string, opts ...FlagOption) *string {
	return CommandLine.String(name, value, usage, opts...)
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
//
// Deprecated: Use String with the WithShorthand option instead
func StringP(name, shorthand string, value string, usage string) *string {
	return CommandLine.StringP(name, shorthand, value, usage)
}
