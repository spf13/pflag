package pflag

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPSlice(t *testing.T) {
	t.Parallel()

	newFlag := func(ipsp *[]net.IP) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.IPSliceVar(ipsp, "ips", []net.IP{}, "Command separated list!")
		return f
	}

	t.Run("with empty slice", func(t *testing.T) {
		ips := make([]net.IP, 0)
		f := newFlag(&ips)
		require.NoError(t, f.Parse([]string{}))

		getIPS, err := f.GetIPSlice("ips")
		require.NoErrorf(t, err,
			"got an error from GetIPSlice(): %v", err,
		)
		require.Empty(t, getIPS)
	})

	t.Run("with values", func(t *testing.T) {
		vals := []string{"192.168.1.1", "10.0.0.1", "0:0:0:0:0:0:0:2"}
		ips := make([]net.IP, 0, len(vals))
		f := newFlag(&ips)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--ips=%s", strings.Join(vals, ",")),
		}))

		for i, v := range ips {
			ip := net.ParseIP(vals[i])
			require.NotNilf(t, ip,
				"invalid string being converted to IP address: %s", vals[i],
			)
			require.Truef(t, ip.Equal(v),
				"expected ips[%d] to be %s but got: %s from GetIPSlice", i, vals[i], v,
			)
		}
	})

	newFlagWithDefault := func(ipsp *[]net.IP) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.IPSliceVar(ipsp, "ips",
			[]net.IP{
				net.ParseIP("192.168.1.1"),
				net.ParseIP("0:0:0:0:0:0:0:1"),
			},
			"Command separated list!")
		return f
	}

	t.Run("with defaults (1)", func(t *testing.T) {
		vals := []string{"192.168.1.1", "0:0:0:0:0:0:0:1"}
		ips := make([]net.IP, 0, len(vals))
		f := newFlagWithDefault(&ips)

		require.NoError(t, f.Parse([]string{}))

		for i, v := range ips {
			ip := net.ParseIP(vals[i])
			require.NotNilf(t, ip,
				"invalid string being converted to IP address: %s", vals[i],
			)
			require.Truef(t, ip.Equal(v),
				"expected ips[%d] to be %s but got: %s", i, vals[i], v,
			)
		}

		getIPS, eri := f.GetIPSlice("ips")
		require.NoErrorf(t, eri,
			"got an error from GetIPSlice: %v", eri,
		)

		for i, v := range getIPS {
			ip := net.ParseIP(vals[i])
			require.NotNilf(t, ip,
				"invalid string being converted to IP address: %s", vals[i],
			)
			require.Truef(t, ip.Equal(v),
				"expected ips[%d] to be %s but got: %s", i, vals[i], v,
			)
		}
	})

	t.Run("with defaults (2)", func(t *testing.T) {
		vals := []string{"192.168.1.1", "0:0:0:0:0:0:0:1"}
		ips := make([]net.IP, 0, len(vals))
		f := newFlagWithDefault(&ips)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf("--ips=%s", strings.Join(vals, ",")),
		}))

		for i, v := range ips {
			ip := net.ParseIP(vals[i])
			require.NotNilf(t, ip,
				"invalid string being converted to IP address: %s", vals[i],
			)
			require.Truef(t, ip.Equal(v),
				"expected ips[%d] to be %s but got: %s", i, vals[i], v,
			)
		}

		getIPS, err := f.GetIPSlice("ips")
		require.NoErrorf(t, err,
			"got an error from GetIPSlice: %v", err,
		)

		for i, v := range getIPS {
			ip := net.ParseIP(vals[i])
			require.NotNilf(t, ip,
				"invalid string being converted to IP address: %s", vals[i],
			)
			require.Truef(t, ip.Equal(v),
				"expected ips[%d] to be %s but got: %s", i, vals[i], v,
			)
		}
	})

	t.Run("called twice", func(t *testing.T) {
		const argfmt = "--ips=%s"
		in := []string{"192.168.1.2,0:0:0:0:0:0:0:1", "10.0.0.1"}
		ips := make([]net.IP, 0, len(in))
		f := newFlag(&ips)
		expected := []net.IP{net.ParseIP("192.168.1.2"), net.ParseIP("0:0:0:0:0:0:0:1"), net.ParseIP("10.0.0.1")}

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		require.Equal(t, expected, ips)
	})

	t.Run("as slice value", func(t *testing.T) {
		const argfmt = "--ips=%s"
		in := []string{"192.168.1.1", "0:0:0:0:0:0:0:1"}
		ips := make([]net.IP, 0, len(in))
		f := newFlag(&ips)

		require.NoError(t, f.Parse([]string{
			fmt.Sprintf(argfmt, in[0]),
			fmt.Sprintf(argfmt, in[1]),
		}))

		f.VisitAll(func(f *Flag) {
			if val, ok := f.Value.(SliceValue); ok {
				require.NoError(t, val.Replace([]string{"192.168.1.2"}))
			}
		})

		require.Equalf(t, []net.IP{net.ParseIP("192.168.1.2")}, ips,
			"expected ss to be overwritten with '192.168.1.2', but got: %v", ips,
		)
	})

	t.Run("bad quoting", func(t *testing.T) {
		tests := []struct {
			Want    []net.IP
			FlagArg []string
		}{
			{
				Want: []net.IP{
					net.ParseIP("a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568"),
					net.ParseIP("203.107.49.208"),
					net.ParseIP("14.57.204.90"),
				},
				FlagArg: []string{
					"a4ab:61d:f03e:5d7d:fad7:d4c2:a1a5:568",
					"203.107.49.208",
					"14.57.204.90",
				},
			},
			{
				Want: []net.IP{
					net.ParseIP("204.228.73.195"),
					net.ParseIP("86.141.15.94"),
				},
				FlagArg: []string{
					"204.228.73.195",
					"86.141.15.94",
				},
			},
			{
				Want: []net.IP{
					net.ParseIP("c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f"),
					net.ParseIP("4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472"),
				},
				FlagArg: []string{
					"c70c:db36:3001:890f:c6ea:3f9b:7a39:cc3f",
					"4d17:1d6e:e699:bd7a:88c5:5e7e:ac6a:4472",
				},
			},
			{
				Want: []net.IP{
					net.ParseIP("5170:f971:cfac:7be3:512a:af37:952c:bc33"),
					net.ParseIP("93.21.145.140"),
					net.ParseIP("2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca"),
				},
				FlagArg: []string{
					" 5170:f971:cfac:7be3:512a:af37:952c:bc33  , 93.21.145.140     ",
					"2cac:61d3:c5ff:6caf:73e0:1b1a:c336:c1ca",
				},
			},
			{
				Want: []net.IP{
					net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
					net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
					net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
					net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"),
				},
				FlagArg: []string{
					`"2e5e:66b2:6441:848:5b74:76ea:574c:3a7b,        2e5e:66b2:6441:848:5b74:76ea:574c:3a7b,2e5e:66b2:6441:848:5b74:76ea:574c:3a7b     "`,
					" 2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"},
			},
		}

		for i, test := range tests {
			var ips []net.IP
			f := newFlag(&ips)

			if err := f.Parse([]string{fmt.Sprintf("--ips=%s", strings.Join(test.FlagArg, ","))}); err != nil {
				t.Fatalf("flag parsing failed with error: %s\nparsing:\t%#v\nwant:\t\t%s",
					err, test.FlagArg, test.Want[i])
			}

			for j, b := range ips {
				if !b.Equal(test.Want[j]) {
					t.Fatalf("bad value parsed for test %d on net.IP %d:\nwant:\t%s\ngot:\t%s", i, j, test.Want[j], b)
				}
			}
		}
	})
}
