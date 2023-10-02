package pflag

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintUsage(t *testing.T) {
	t.Run("with print", func(t *testing.T) {
		f := NewFlagSet("test", ExitOnError)

		f.Bool("long-form", false, "Some description")
		f.Bool("long-form2", false, "Some description\n  with multiline")
		f.BoolP("long-name", "s", false, "Some description")
		f.BoolP("long-name2", "t", false, "Some description with\n  multiline")

		const expectedOutput = `      --long-form    Some description
      --long-form2   Some description
                       with multiline
  -s, --long-name    Some description
  -t, --long-name2   Some description with
                       multiline
`

		require.Equal(t, expectedOutput, printFlagDefaults(f))
	})

	t.Run("with wrapped columns", func(t *testing.T) {
		const cols = 80

		f := NewFlagSet("test", ExitOnError)

		f.Bool("long-form", false, "Some description")
		f.Bool("long-form2", false, "Some description\n  with multiline")
		f.BoolP("long-name", "s", false, "Some description")
		f.BoolP("long-name2", "t", false, "Some description with\n  multiline")
		f.StringP("some-very-long-arg", "l", "test", "Some very long description having break the limit")
		f.StringP("other-very-long-arg", "o", "long-default-value", "Some very long description having break the limit")
		f.String("some-very-long-arg2", "very long default value", "Some very long description\nwith line break\nmultiple")

		const expectedOutput = `      --long-form                    Some description
      --long-form2                   Some description
                                       with multiline
  -s, --long-name                    Some description
  -t, --long-name2                   Some description with
                                       multiline
  -o, --other-very-long-arg string   Some very long description having
                                     break the limit (default
                                     "long-default-value")
  -l, --some-very-long-arg string    Some very long description having
                                     break the limit (default "test")
      --some-very-long-arg2 string   Some very long description
                                     with line break
                                     multiple (default "very long default
                                     value")
`

		require.Equal(t, expectedOutput, f.FlagUsagesWrapped(cols))
	})
}
