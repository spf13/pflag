package pflag

// -- func Value
type funcValue func(string) error

func (f funcValue) Set(s string) error {
	return f(s)
}

func (funcValue) Type() string {
	return "value"
}

func (funcValue) String() string { return "" }

// Func defines a flag with the specified name and usage string that calls the
// specified function with its command-line argument like Value.Set.  Unlike
// other functions that define flags, Func does not return any value.
func (f *FlagSet) Func(name, usage string, fn func(string) error) {
	f.VarP(funcValue(fn), name, "", usage)
}

// FuncP is like Func, but accepts a shorthand letter that can be used after a single dash.
func (f *FlagSet) FuncP(name, shorthand, usage string, fn func(string) error) {
	f.VarP(funcValue(fn), name, shorthand, usage)
}

// Func defines a flag with the specified name and usage string that calls the
// specified function with its command-line argument like Value.Set.  Unlike
// other functions that define flags, Func does not return any value.
func Func(name, usage string, fn func(string) error) {
	CommandLine.FuncP(name, "", usage, fn)
}

// FuncP is like Func, but accepts a shorthand letter that can be used after a single dash.
func FuncP(name, shorthand, usage string, fn func(string) error) {
	CommandLine.FuncP(name, shorthand, usage, fn)
}
