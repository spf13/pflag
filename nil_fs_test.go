package pflag

import "testing"

func TestNilFlagSetLookup(t *testing.T) {
	var f *FlagSet
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panicked: %v", r)
		}
	}()
	if f.Lookup("x") != nil {
		t.Fatal("want nil")
	}
	if f.ShorthandLookup("x") != nil {
		t.Fatal("want nil")
	}
}
