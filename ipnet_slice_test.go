package pflag

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPNetSlice(t *testing.T) {
	t.Parallel()

	newFlag := func(ipsp *[]net.IPNet) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.IPNetSliceVar(ipsp, "cidrs", []net.IPNet{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		cidrs := make([]net.IPNet, 0)
		f := newFlag(&cidrs)
		require.NoError(t, f.Parse([]string{}))

		getIPNet, err := f.GetIPNetSlice("cidrs")
		require.NoErrorf(t, err,
			"got an error from GetIPNetSlice(): %v", err,
		)
		require.Empty(t, getIPNet)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"192.168.1.1/24", "10.0.0.1/16", "fd00:0:0:0:0:0:0:2/64"}
		ips := make([]net.IPNet, 0, len(vals))
		f := newFlag(&ips)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--cidrs=%s", strings.Join(vals, ",")),
		}))

		for i, v := range ips {
			_, cidr, _ := net.ParseCIDR(vals[i])
			require.NotNilf(t, cidr,
				"invalid string being converted to CIDR: %s", vals[i],
			)
			require.Truef(t, equalCIDR(*cidr, v),
				"expected ips[%d] to be %s but got: %s from GetIPSlice", i, vals[i], v,
			)
		}
	})

	newFlagWithDefault := func(ipsp *[]net.IPNet) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.IPNetSliceVar(ipsp, "cidrs",
			[]net.IPNet{
				getCIDR(net.ParseCIDR("192.168.1.1/16")),
				getCIDR(net.ParseCIDR("fd00::/64")),
			},
			"Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"192.168.1.1/16", "fd00::/64"}
		cidrs := make([]net.IPNet, 0, len(vals))
		f := newFlagWithDefault(&cidrs)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range cidrs {
			_, cidr, _ := net.ParseCIDR(vals[i])
			require.NotNilf(t, cidr,
				"invalid string being converted to CIDR: %s", vals[i],
			)
			require.Truef(t, equalCIDR(*cidr, v),
				"expected cidrs[%d] to be %s but got: %s", i, vals[i], v,
			)
		}

		getIPNet, err := f.GetIPNetSlice("cidrs")
		require.NoErrorf(t, err,
			"got an error from GetIPNetSlice: %v", err,
		)

		for i, v := range getIPNet {
			_, cidr, _ := net.ParseCIDR(vals[i])
			require.NotNilf(t, cidr,
				"invalid string being converted to CIDR: %s", vals[i],
			)
			require.Truef(t, equalCIDR(*cidr, v),
				"expected cidrs[%d] to be %s but got: %s", i, vals[i], v,
			)
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"192.168.1.1/16", "fd00::/64"}
		cidrs := make([]net.IPNet, 0, len(vals))
		f := newFlagWithDefault(&cidrs)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--cidrs=%s", strings.Join(vals, ",")),
		}))

		for i, v := range cidrs {
			_, cidr, _ := net.ParseCIDR(vals[i])
			require.NotNilf(t, cidr,
				"invalid string being converted to CIDR: %s", vals[i],
			)
			require.Truef(t, equalCIDR(*cidr, v),
				"expected cidrs[%d] to be %s but got: %s", i, vals[i], v,
			)
		}

		getIPNet, err := f.GetIPNetSlice("cidrs")
		require.NoErrorf(t, err,
			"got an error from GetIPNetSlice: %v", err,
		)

		for i, v := range getIPNet {
			_, cidr, _ := net.ParseCIDR(vals[i])
			require.NotNilf(t, cidr,
				"invalid string being converted to CIDR: %s", vals[i],
			)
			require.Truef(t, equalCIDR(*cidr, v),
				"expected cidrs[%d] to be %s but got: %s", i, vals[i], v,
			)
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--cidrs=%s"
		in := []string{"192.168.1.2/16,fd00::/64", "10.0.0.1/24"}
		cidrs := make([]net.IPNet, 0, len(in))
		f := newFlag(&cidrs)
		expected := []net.IPNet{
			getCIDR(net.ParseCIDR("192.168.1.2/16")),
			getCIDR(net.ParseCIDR("fd00::/64")),
			getCIDR(net.ParseCIDR("10.0.0.1/24")),
		}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, cidrs)
	})

	t.Run("bad quoting", func(t *testing.T) {
		tests := []struct {
			Want    []net.IPNet
			FlagArg []string
		}{
			{
				Want: []net.IPNet{
					getCIDR(net.ParseCIDR("a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568/128")),
					getCIDR(net.ParseCIDR("203.107.49.208/32")),
					getCIDR(net.ParseCIDR("14.57.204.90/32")),
				},
				FlagArg: []string{
					"a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568/128",
					"203.107.49.208/32",
					"14.57.204.90/32",
				},
			},
			{
				Want: []net.IPNet{
					getCIDR(net.ParseCIDR("204.228.73.195/32")),
					getCIDR(net.ParseCIDR("86.141.15.94/32")),
				},
				FlagArg: []string{
					"204.228.73.195/32",
					"86.141.15.94/32",
				},
			},
			{
				Want: []net.IPNet{
					getCIDR(net.ParseCIDR("c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f/128")),
					getCIDR(net.ParseCIDR("4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472/128")),
				},
				FlagArg: []string{
					"c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f/128",
					"4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472/128",
				},
			},
			{
				Want: []net.IPNet{
					getCIDR(net.ParseCIDR("5170:f971:cfac:7be3:512a:af37:952c:bc33/128")),
					getCIDR(net.ParseCIDR("93.21.145.140/32")),
					getCIDR(net.ParseCIDR("2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca/128")),
				},
				FlagArg: []string{
					" 5170:f971:cfac:7be3:512a:af37:952c:bc33/128  , 93.21.145.140/32     ",
					"2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca/128",
				},
			},
			{
				Want: []net.IPNet{
					getCIDR(net.ParseCIDR("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128")),
					getCIDR(net.ParseCIDR("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128")),
					getCIDR(net.ParseCIDR("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128")),
					getCIDR(net.ParseCIDR("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128")),
				},
				FlagArg: []string{
					`"2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128,        2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128,2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128     "`,
					" 2e5e:66b2:6441:848:5b74:76ea:574c:3a7b/128"},
			},
		}

		for i, test := range tests {
			cidrs := make([]net.IPNet, 0, len(test.Want))
			f := newFlag(&cidrs)

			require.NoErrorf(t, f.Parse([]string{fmt.Sprintf("--cidrs=%s", strings.Join(test.FlagArg, ","))}),
				"flag parsing failed with error:\nparsing:\t%#v\nwant:\t\t%s",
				test.FlagArg, test.Want,
			)

			for j, b := range cidrs {
				require.Truef(t, equalCIDR(b, test.Want[j]),
					"bad value parsed for test %d on net.IP %d:\nwant:\t%s\ngot:\t%s", i, j, test.Want[j], b,
				)
			}
		}
	})
}

func getCIDR(_ net.IP, cidr *net.IPNet, _ error) net.IPNet {
	return *cidr
}

func equalCIDR(c1 net.IPNet, c2 net.IPNet) bool {
	return c1.String() == c2.String()
}
