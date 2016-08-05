package pflag

import (
	"fmt"
	"os"
)

// IsSet indicates whether the specified flag is set in the given FlagSet
func (fs *FlagSet) IsSet(name string) bool {
	normalName := fs.normalizeFlagName(name)

	return fs.actual[normalName] != nil
}

func IsSet(name string) bool {
	return CommandLine.IsSet(name)
}

// Merge is a helper function that merges n FlagSets into a single dest FlagSet
// In case of name collision between the flagsets it will apply
// the destination FlagSet's errorHandling behavior.
func Merge(dest *FlagSet, flagsets ...*FlagSet) error {

	if dest.formal == nil {
		dest.formal = make(map[NormalizedName]*Flag)
	}

	for _, fset := range flagsets {
		if fset.formal == nil {
			continue
		}
		for k, f := range fset.formal {
			c := f.Shorthand[0]

			_, sok := dest.shorthands[c]
			_, ok := dest.formal[k]

			if sok || ok {
				var err error
				if fset.name == "" {
					err = fmt.Errorf("flag redefined: %s", k)
				} else {
					err = fmt.Errorf("%s flag redefined: %s,short: %s", fset.name, k, f.Shorthand)
				}
				fmt.Fprintln(fset.out(), err.Error())
				// Happens only if flags are declared with identical names
				switch dest.errorHandling {
				case ContinueOnError:
					return err
				case ExitOnError:
					os.Exit(2)
				case PanicOnError:
					panic(err)
				}
			}

			dest.formal[k] = fset.formal[k]
			dest.shorthands[c] = fset.formal[k]
		}
	}
	return nil
}

type nArgRequirementType int

// Indicator used to pass to BadArgs function
const (
	Exact nArgRequirementType = iota
	Max
	Min
)

type nArgRequirement struct {
	Type nArgRequirementType
	N    int
}

// Require adds a requirement about the number of arguments for the FlagSet.
// The first parameter can be Exact, Max, or Min to respectively specify the exact,
// the maximum, or the minimal number of arguments required.
// The actual check is done in FlagSet.CheckArgs().
func (fs *FlagSet) Require(nArgRequirementType nArgRequirementType, nArg int) {
	fs.nArgRequirements = append(fs.nArgRequirements, nArgRequirement{nArgRequirementType, nArg})
}

// CheckArgs uses the requirements set by FlagSet.Require() to validate
// the number of arguments. If the requirements are not met,
// an error message string is returned.
func (fs *FlagSet) CheckArgs() (message string) {
	for _, req := range fs.nArgRequirements {
		var arguments string
		if req.N == 1 {
			arguments = "1 argument"
		} else {
			arguments = fmt.Sprintf("%d arguments", req.N)
		}

		str := func(kind string) string {
			return fmt.Sprintf("%q requires %s%s", fs.name, kind, arguments)
		}

		switch req.Type {
		case Exact:
			if fs.NArg() != req.N {
				return str("")
			}
		case Max:
			if fs.NArg() > req.N {
				return str("a maximum of ")
			}
		case Min:
			if fs.NArg() < req.N {
				return str("a minimum of ")
			}
		}
	}
	return ""
}

// ParseFlags calls fs.Parse(args) and prints a relevant error message if there are
// incorrect number of arguments. It returns error only if error handling is
// set to ContinueOnError and parsing fails. If error handling is set to
// ExitOnError, it's safe to ignore the return value.
func (fs *FlagSet) ParseFlags(args []string) error {

	if err := fs.Parse(args); err != nil {
		return err
	}

	if str := fs.CheckArgs(); str != "" {
		fs.ReportError(str)
		os.Exit(1)

	}
	return nil
}

// ReportError is a utility method that prints a user-friendly message
// containing the error that occurred during parsing and a suggestion to get help
func (fs *FlagSet) ReportError(str string) {
	if os.Args[0] == fs.name {
		str += ".\nSee '" + os.Args[0] + " --help'"
	} else {
		str += ".\nSee '" + os.Args[0] + " " + fs.name + " --help'"
	}
	fmt.Fprintf(fs.out(), "%s: %s.\n", os.Args[0], str)
}
