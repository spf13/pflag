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
