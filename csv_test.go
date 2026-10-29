package pflag

import (
	"net"
	"reflect"
	"testing"
)

func TestSliceFlagsRejectExtraCSVRecords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		register func(*FlagSet)
	}{
		{"strings", "first\nsecond", func(f *FlagSet) { f.StringSlice("value", []string{"default"}, "") }},
		{"bools", "true\nfalse", func(f *FlagSet) { f.BoolSlice("value", []bool{false}, "") }},
		{"ips", "192.0.2.1\ninvalid", func(f *FlagSet) { f.IPSlice("value", []net.IP{net.ParseIP("192.0.2.2")}, "") }},
		{"networks", "192.0.2.0/24\ninvalid", func(f *FlagSet) { f.IPNetSlice("value", nil, "") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			tt.register(f)
			flag := f.Lookup("value")
			before := flag.Value.String()
			if err := f.Parse([]string{"--value", tt.input}); err == nil {
				t.Error("expected an error for multiple CSV records")
			}
			if got := flag.Value.String(); got != before {
				t.Errorf("failed parsing changed the value: got %q, want %q", got, before)
			}
			if flag.Changed {
				t.Error("failed parsing marked the flag as changed")
			}
		})
	}
}

func TestStringSliceQuotedNewline(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	values := f.StringSlice("value", nil, "")
	if err := f.Parse([]string{"--value", "\"first\nsecond\",third"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"first\nsecond", "third"}
	if !reflect.DeepEqual(*values, want) {
		t.Fatalf("got %q, want %q", *values, want)
	}
	got, err := f.GetStringSlice("value")
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("GetStringSlice() = %q, %v; want %q, nil", got, err, want)
	}
}
