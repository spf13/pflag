package pflag

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPNet(t *testing.T) {
	newFlag := func(ip *net.IPNet) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		_, def, _ := net.ParseCIDR("0.0.0.0/0")
		f.IPNetVar(ip, "address", *def, "IP Address")
		return f
	}

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
		f := newFlag(&addr)
		tc := &testCases[i]

		err := f.Parse([]string{
			fmt.Sprintf("--address=%s", tc.input),
		})
		if !tc.success {
			require.Errorf(t, err, "expected failure")

			continue
		}

		require.NoErrorf(t, err, "expected success, got %q", err)

		ip, err := f.GetIPNet("address")
		require.NoErrorf(t, err,
			"got error trying to fetch the IPnet flag: %v", err,
		)
		require.Equal(t, tc.expected, ip.String())
	}
}
