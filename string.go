package pflag

// -- string Value
type StringValue string

func newStringValue(val string, p *string) *StringValue {
	*p = val
	return (*StringValue)(p)
}

func (s *StringValue) Set(val string) error {
	*s = StringValue(val)
	return nil
}
func (s *StringValue) Type() string {
	return "string"
}

func (s *StringValue) String() string { return string(*s) }

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
func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.VarP(newStringValue(value, p), name, "", usage)
}

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringVarP(p *string, name, shorthand string, value string, usage string) {
	f.VarP(newStringValue(value, p), name, shorthand, usage)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string) {
	CommandLine.VarP(newStringValue(value, p), name, "", usage)
}

// StringVarP is like StringVar, but accepts a shorthand letter that can be used after a single dash.
func StringVarP(p *string, name, shorthand string, value string, usage string) {
	CommandLine.VarP(newStringValue(value, p), name, shorthand, usage)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (f *FlagSet) String(name string, value string, usage string) *string {
	p := new(string)
	f.StringVarP(p, name, "", value, usage)
	return p
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) StringP(name, shorthand string, value string, usage string) *string {
	p := new(string)
	f.StringVarP(p, name, shorthand, value, usage)
	return p
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name string, value string, usage string) *string {
	return CommandLine.StringP(name, "", value, usage)
}

// StringP is like String, but accepts a shorthand letter that can be used after a single dash.
func StringP(name, shorthand string, value string, usage string) *string {
	return CommandLine.StringP(name, shorthand, value, usage)
}
