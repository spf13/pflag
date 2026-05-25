package pflag

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func setUpIP(ip *net.IP) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.IPVar(ip, "address", net.ParseIP("0.0.0.0"), "IP Address")
	return f
}

func TestIP(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		{"0.0.0.0", true, "0.0.0.0"},
		{" 0.0.0.0 ", true, "0.0.0.0"},
		{"1.2.3.4", true, "1.2.3.4"},
		{"127.0.0.1", true, "127.0.0.1"},
		{"255.255.255.255", true, "255.255.255.255"},
		{"", true, "0.0.0.0"},
		{"0", false, ""},
		{"localhost", false, ""},
		{"0.0.0", false, ""},
		{"0.0.0.", false, ""},
		{"0.0.0.0.", false, ""},
		{"0.0.0.256", false, ""},
		{"0 . 0 . 0 . 0", false, ""},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		var addr net.IP
		f := setUpIP(&addr)

		tc := &testCases[i]

		arg := fmt.Sprintf("--address=%s", tc.input)
		err := f.Parse([]string{arg})
		if err != nil && tc.success == true {
			t.Errorf("expected success, got %q", err)
			continue
		} else if err == nil && tc.success == false {
			t.Errorf("expected failure")
			continue
		} else if tc.success {
			ip, err := f.GetIP("address")
			if err != nil {
				t.Errorf("Got error trying to fetch the IP flag: %v", err)
			}
			if ip.String() != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, ip.String())
			}
		}
	}
}

// TestIPNilDefault covers #351: declaring an IP flag with a nil default and
// then reading it without setting it must return (nil, nil), not an error.
// Before the ipConv guard, the unset flag's String() returned "<nil>" and
// GetIP tried to parse it as an address.
func TestIPNilDefault(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.IP("ip", nil, "")

	ip, err := f.GetIP("ip")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != nil {
		t.Errorf("expected nil IP, got %v", ip)
	}

	var bound net.IP
	f2 := NewFlagSet("test2", ContinueOnError)
	f2.IPVar(&bound, "ip", nil, "")
	ip, err = f2.GetIP("ip")
	if err != nil {
		t.Fatalf("unexpected error (IPVar variant): %v", err)
	}
	if ip != nil {
		t.Errorf("expected nil IP (IPVar variant), got %v", ip)
	}
}
