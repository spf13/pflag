package pflag

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func setUpIPNet(ip *net.IPNet) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	_, def, _ := net.ParseCIDR("0.0.0.0/0")
	f.IPNetVar(ip, "address", *def, "IP Address")
	return f
}

func TestIPNet(t *testing.T) {
	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		{"0.0.0.0/0", true, "0.0.0.0/0"},
		{" 0.0.0.0/0 ", true, "0.0.0.0/0"},
		{"1.2.3.4/8", true, "1.0.0.0/8"},
		{"127.0.0.1/16", true, "127.0.0.0/16"},
		{"255.255.255.255/19", true, "255.255.224.0/19"},
		{"255.255.255.255/32", true, "255.255.255.255/32"},
		{"", false, ""},
		{"/0", false, ""},
		{"0", false, ""},
		{"0/0", false, ""},
		{"localhost/0", false, ""},
		{"0.0.0/4", false, ""},
		{"0.0.0./8", false, ""},
		{"0.0.0.0./12", false, ""},
		{"0.0.0.256/16", false, ""},
		{"0.0.0.0 /20", false, ""},
		{"0.0.0.0/ 24", false, ""},
		{"0 . 0 . 0 . 0 / 28", false, ""},
		{"0.0.0.0/33", false, ""},
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull
	for i := range testCases {
		var addr net.IPNet
		f := setUpIPNet(&addr)

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
			ip, err := f.GetIPNet("address")
			if err != nil {
				t.Errorf("Got error trying to fetch the IP flag: %v", err)
			}
			if ip.String() != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, ip.String())
			}
		}
	}
}

func TestGetIPNetDefaults(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   net.IPNet
		want    string
		wantErr bool
	}{
		{name: "zero", want: "<nil>"},
		{name: "IPv4", value: net.IPNet{IP: net.IPv4(192, 0, 2, 3), Mask: net.CIDRMask(24, 32)}, want: "192.0.2.0/24"},
		{name: "IPv6", value: net.IPNet{IP: net.ParseIP("2001:db8::1"), Mask: net.CIDRMask(32, 128)}, want: "2001:db8::/32"},
		{name: "IP only", value: net.IPNet{IP: net.IPv4(192, 0, 2, 3)}, wantErr: true},
		{name: "mask only", value: net.IPNet{Mask: net.CIDRMask(24, 32)}, wantErr: true},
		{name: "empty IP", value: net.IPNet{IP: net.IP{}}, wantErr: true},
		{name: "empty mask", value: net.IPNet{Mask: net.IPMask{}}, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := NewFlagSet("test", ContinueOnError)
			f.IPNet("network", tc.value, "optional network")
			got, err := f.GetIPNet("network")
			if (err != nil) != tc.wantErr {
				t.Fatalf("GetIPNet default %#v = %v, %v", tc.value, got, err)
			}
			if tc.wantErr {
				return
			}
			if got.String() != tc.want || (tc.want == "<nil>" && (got.IP != nil || got.Mask != nil)) {
				t.Fatalf("GetIPNet = %#v, want %s", got, tc.want)
			}
			if len(got.IP) > 0 {
				got.IP[0] ^= 0xff
				got.Mask[0] ^= 0xff
			}
			again, err := f.GetIPNet("network")
			if err != nil || again.String() != tc.want {
				t.Fatalf("mutating getter result changed stored value: %v, %v", again, err)
			}
		})
	}
}

type ipNetTextValue struct{ stringValue }

func (*ipNetTextValue) Type() string { return "ipNet" }

func TestGetIPNetCustomValue(t *testing.T) {
	for _, input := range []string{"<nil>", "192.0.2.3/24"} {
		f := NewFlagSet("test", ContinueOnError)
		f.Var(&ipNetTextValue{stringValue(input)}, "network", "custom network")
		got, err := f.GetIPNet("network")
		if input == "<nil>" {
			if err == nil {
				t.Fatal("custom <nil> representation unexpectedly accepted")
			}
		} else if err != nil || got.String() != "192.0.2.0/24" {
			t.Fatalf("custom CIDR conversion = %v, %v", got, err)
		}
	}
}
