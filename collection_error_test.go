package pflag

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var collectionErrorCases = []struct {
	name     string
	register func(*FlagSet)
	prefix   string
	valid    string
	invalid  string
	hint     string
	overflow string
}{
	{"boolSlice", func(f *FlagSet) { f.BoolSlice("value", nil, "") }, "true,", "true", "maybe", "must be true or false", ""},
	{"intSlice", func(f *FlagSet) { f.IntSlice("value", nil, "") }, "1,", "1", "x", "must be an integer", "9223372036854775808"},
	{"int32Slice", func(f *FlagSet) { f.Int32Slice("value", nil, "") }, "1,", "1", "x", "must be an integer", "2147483648"},
	{"int64Slice", func(f *FlagSet) { f.Int64Slice("value", nil, "") }, "1,", "1", "x", "must be an integer", "9223372036854775808"},
	{"uintSlice", func(f *FlagSet) { f.UintSlice("value", nil, "") }, "1,", "1", "-1", "must be a non-negative integer", "18446744073709551616"},
	{"float32Slice", func(f *FlagSet) { f.Float32Slice("value", nil, "") }, "1,", "1", "x", "must be a number", "1e39"},
	{"float64Slice", func(f *FlagSet) { f.Float64Slice("value", nil, "") }, "1,", "1", "x", "must be a number", "1e309"},
	{"durationSlice", func(f *FlagSet) { f.DurationSlice("value", nil, "") }, "1s,", "1s", "soon", `must be a duration like "30s" or "5m"`, ""},
	{"stringToInt", func(f *FlagSet) { f.StringToInt("value", nil, "") }, "ok=1,bad=", "ok=1", "x", "must be an integer", "9223372036854775808"},
	{"stringToInt64", func(f *FlagSet) { f.StringToInt64("value", nil, "") }, "ok=1,bad=", "ok=1", "x", "must be an integer", "9223372036854775808"},
}

func TestCollectionArgumentErrors(t *testing.T) {
	for _, tc := range collectionErrorCases {
		t.Run(tc.name, func(t *testing.T) {
			inputs := []struct{ value, hint string }{{tc.invalid, tc.hint}}
			if tc.overflow != "" {
				inputs = append(inputs, struct{ value, hint string }{tc.overflow, "is out of range"})
			}
			for _, input := range inputs {
				for _, previouslySet := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/previouslySet=%t", input.value, previouslySet), func(t *testing.T) {
						f := NewFlagSet("test", ContinueOnError)
						f.SetOutput(ioutil.Discard)
						tc.register(f)
						if previouslySet {
							if err := f.Parse([]string{"--value=" + tc.valid}); err != nil {
								t.Fatal(err)
							}
						}
						flag := f.Lookup("value")
						before := flag.Value.String()
						raw := tc.prefix + input.value
						err := f.Parse([]string{"--value=" + raw})
						if err == nil {
							t.Fatalf("Parse(%q) succeeded", raw)
						}
						want := fmt.Sprintf("invalid argument %q for %q flag: %q %s", raw, "--value", input.value, input.hint)
						if err.Error() != want {
							t.Errorf("Parse(%q) error = %q, want %q", raw, err, want)
						}
						if got := flag.Value.String(); got != before {
							t.Errorf("failed Parse changed value from %q to %q", before, got)
						}
						if flag.Changed != previouslySet {
							t.Errorf("failed Parse changed Changed from %t to %t", previouslySet, flag.Changed)
						}
					})
				}
			}
		})
	}
}
