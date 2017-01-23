package pflag

import (
	"fmt"
	"net"
	"strings"
	"testing"
)

func setUpIPSFlagSet(ipsp *[]net.IP) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.IPSliceVar(ipsp, "ips", []net.IP{}, "Command separated list!")
	return f
}

func setUpIPSFlagSetWithDefault(ipsp *[]net.IP) *FlagSet {
	f := NewFlagSet("test", ContinueOnError)
	f.IPSliceVar(ipsp, "ips",
		[]net.IP{
			net.ParseIP("192.168.1.1"),
			net.ParseIP("0:0:0:0:0:0:0:1"),
		},
		"Command separated list!")
	return f
}

func TestEmptyIP(t *testing.T) {
	var ips []net.IP
	f := setUpIPSFlagSet(&ips)
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}

	getIPS, err := f.GetIPSlice("ips")
	if err != nil {
		t.Fatal("got an error from GetIPSlice():", err)
	}
	if len(getIPS) != 0 {
		t.Fatalf("got ips %v with len=%d but expected length=0", getIPS, len(getIPS))
	}
}

func TestIPS(t *testing.T) {
	var ips []net.IP
	f := setUpIPSFlagSet(&ips)

	vals := []string{"192.168.1.1", "10.0.0.1", "0:0:0:0:0:0:0:2"}
	arg := fmt.Sprintf("--ips=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v) {
			t.Fatalf("expected ips[%d] to be %s but got: %s from GetIPSlice", i, vals[i], v)
		}
	}
}

func TestIPSDefault(t *testing.T) {
	var ips []net.IP
	f := setUpIPSFlagSetWithDefault(&ips)

	vals := []string{"192.168.1.1", "0:0:0:0:0:0:0:1"}
	err := f.Parse([]string{})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		}
	}

	getIPS, err := f.GetIPSlice("ips")
	if err != nil {
		t.Fatal("got an error from GetIPSlice")
	}
	for i, v := range getIPS {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		}
	}
}

func TestIPSWithDefault(t *testing.T) {
	var ips []net.IP
	f := setUpIPSFlagSetWithDefault(&ips)

	vals := []string{"192.168.1.1", "0:0:0:0:0:0:0:1"}
	arg := fmt.Sprintf("--ips=%s", strings.Join(vals, ","))
	err := f.Parse([]string{arg})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		}
	}

	getIPS, err := f.GetIPSlice("ips")
	if err != nil {
		t.Fatal("got an error from GetIPSlice")
	}
	for i, v := range getIPS {
		if ip := net.ParseIP(vals[i]); ip == nil {
			t.Fatalf("invalid string being converted to IP address: %s", vals[i])
		} else if !ip.Equal(v) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, vals[i], v)
		}
	}
}

func TestIPSCalledTwice(t *testing.T) {
	var ips []net.IP
	f := setUpIPSFlagSet(&ips)

	in := []string{"192.168.1.2,0:0:0:0:0:0:0:1", "10.0.0.1"}
	expected := []net.IP{net.ParseIP("192.168.1.2"), net.ParseIP("0:0:0:0:0:0:0:1"), net.ParseIP("10.0.0.1")}
	argfmt := "ips=%s"
	arg1 := fmt.Sprintf(argfmt, in[0])
	arg2 := fmt.Sprintf(argfmt, in[1])
	err := f.Parse([]string{arg1, arg2})
	if err != nil {
		t.Fatal("expected no error; got", err)
	}
	for i, v := range ips {
		if !expected[i].Equal(v) {
			t.Fatalf("expected ips[%d] to be %s but got: %s", i, expected[i], v)
		}
	}
}
