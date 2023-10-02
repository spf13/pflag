package pflag

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIP(t *testing.T) {
	newFlag := func(ip *net.IP) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.IPVar(ip, "address", net.ParseIP("0.0.0.0"), "IP Address")
		return f
	}

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
		{"", false, ""},
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

		ip, err := f.GetIP("address")
		require.NoErrorf(t, err,
			"got error trying to fetch the IP flag: %v", err,
		)
		require.Equal(t, tc.expected, ip.String())
	}
}
