package pflag

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

type testOptions struct {
	String      string   `flag:"string" short:"a" default:"abc" desc:"this is string"`
	StringSlice []string `flag:"stringSlice" default:"[a,b]" desc:"this is string slice"`

	Bool      bool   `flag:"bool" short:"b" default:"true" desc:"this is bool"`
	BoolSlice []bool `flag:"boolSlice" default:"[true,false]" desc:"this is bool slice"`

	Int          int       `flag:"int" short:"c" default:"1" desc:"this is int"`
	IntSlice     []int     `flag:"intSlice" default:"[1,2]" desc:"this is int slice"`
	Int8         int8      `flag:"int8" default:"1" desc:"this is int8"`
	Int16        int16     `flag:"int16" default:"1" desc:"this is int16"`
	Int32        int32     `flag:"int32" default:"1" desc:"this is int32"`
	Int32Slice   []int32   `flag:"int32Slice" default:"[1,2]" desc:"this is int32 slice"`
	Int64        int64     `flag:"int64" default:"1" desc:"this is int64"`
	Int64Slice   []int64   `flag:"int64Slice"  default:"[1,2]" desc:"this is int64 slice"`
	Uint         uint      `flag:"uint" default:"1" desc:"this is uint"`
	UintSlice    []uint    `flag:"uintSlice" default:"[1,2]" desc:"this is uint slice"`
	Uint8        uint8     `flag:"uint8" default:"1" desc:"this is uint8"`
	Uint16       uint16    `flag:"uint16" default:"1" desc:"this is uint16"`
	Uint32       uint32    `flag:"uint32" default:"1" desc:"this is uint32"`
	Uint64       uint64    `flag:"uint64" default:"1" desc:"this is uint64"`
	Float32      float32   `flag:"float32" default:"1.1" desc:"this is float32"`
	Float32Slice []float32 `flag:"float32Slice" default:"1.1,2.1" desc:"this is float32 slice"`
	Float64      float64   `flag:"float64" default:"1.1" desc:"this is float64"`
	Float64Slice []float64 `flag:"float64Slice" default:"1.1,2.1" desc:"this is float64 slice"`

	Duration      time.Duration   `flag:"duration" short:"d" default:"1s" desc:"this is duration"`
	DurationSlice []time.Duration `flag:"durationSlice" default:"1s,2m" desc:"this is duration slice"`

	IP      net.IP     `flag:"ip" short:"i" default:"127.0.0.1" desc:"this is ip"`
	IPSlice []net.IP   `flag:"ipSlice" default:"127.0.0.1,127.0.0.2" desc:"this is ip slice"`
	IPMask  net.IPMask `flag:"ipMask" default:"255.255.255.255" desc:"this is ipMask"`
	IPNet   net.IPNet  `flag:"ipNet" default:"192.0.2.1/24" desc:"this is ipNet"`

	Env string `flag:"env" short:"e" default:"${PWD}" desc:"this is env"`

	Others otherInfo `flag:"others"`
}

type otherInfo struct {
	// can't set shorhand
	String string `flag:"string" default:"string" desc:"this is others.string"`
}

func TestSetValues(t *testing.T) {
	f := NewFlagSet("", ContinueOnError)
	if err := f.AddFlags(testOptions{}); err != nil {
		t.Error(err)
	}

	args := []string{"--bool", "true", "--string", "abcd"}
	if err := f.Parse(args); err != nil {
		t.Error(err)
	}

	var opts = testOptions{}
	if err := f.SetValues(&opts); err != nil {
		t.Error(err)
	}

	if opts.Bool != true {
		t.Error("flag(bool) value is incorrent")
	}
	if opts.String != "abcd" {
		t.Error("flag(string) value is incorrent")
	}

}

func TestAddFlags(t *testing.T) {
	f := NewFlagSet("", ContinueOnError)
	err := f.AddFlags(testOptions{})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		Name      string
		Shorthand string
		Def       string
		Usage     string
	}{
		{
			Name:      "string",
			Shorthand: "a",
			Def:       "abc",
			Usage:     "this is string",
		},
		{
			Name:  "stringSlice",
			Def:   "[a,b]",
			Usage: "this is string slice",
		},
		{
			Name:      "bool",
			Shorthand: "b",
			Def:       "true",
			Usage:     "this is bool",
		},
		{
			Name:  "boolSlice",
			Def:   "[true,false]",
			Usage: "this is bool slice",
		},
		{
			Name:      "int",
			Shorthand: "c",
			Def:       "1",
			Usage:     "this is int",
		},
		{
			Name:  "intSlice",
			Def:   "[1,2]",
			Usage: "this is int slice",
		},
		{
			Name:  "int8",
			Def:   "1",
			Usage: "this is int8",
		},
		{
			Name:  "int16",
			Def:   "1",
			Usage: "this is int16",
		},
		{
			Name:  "int32",
			Def:   "1",
			Usage: "this is int32",
		},
		{
			Name:  "int32Slice",
			Def:   "[1,2]",
			Usage: "this is int32 slice",
		},
		{
			Name:  "int64",
			Def:   "1",
			Usage: "this is int64",
		},
		{
			Name:  "int64Slice",
			Def:   "[1,2]",
			Usage: "this is int64 slice",
		},
		{
			Name:  "uint",
			Def:   "1",
			Usage: "this is uint",
		},
		{
			Name:  "uintSlice",
			Def:   "[1,2]",
			Usage: "this is uint slice",
		},
		{
			Name:  "uint8",
			Def:   "1",
			Usage: "this is uint8",
		},
		{
			Name:  "uint16",
			Def:   "1",
			Usage: "this is uint16",
		},
		{
			Name:  "uint32",
			Def:   "1",
			Usage: "this is uint32",
		},
		{
			Name:  "uint64",
			Def:   "1",
			Usage: "this is uint64",
		},
		{
			Name:  "float32",
			Def:   "1.1",
			Usage: "this is float32",
		},
		{
			Name:  "float32Slice",
			Def:   "[1.100000,2.100000]",
			Usage: "this is float32 slice",
		},
		{
			Name:  "float64",
			Def:   "1.1",
			Usage: "this is float64",
		},
		{
			Name:  "float64Slice",
			Def:   "[1.100000,2.100000]",
			Usage: "this is float64 slice",
		},
		{
			Name:      "duration",
			Shorthand: "d",
			Def:       "1s",
			Usage:     "this is duration",
		},
		{
			Name:  "durationSlice",
			Def:   "[1s,2m0s]",
			Usage: "this is duration slice",
		},
		{
			Name:      "ip",
			Shorthand: "i",
			Def:       "127.0.0.1",
			Usage:     "this is ip",
		},
		{
			Name:  "ipSlice",
			Def:   "[127.0.0.1,127.0.0.2]",
			Usage: "this is ip slice",
		},
		{
			Name:  "ipMask",
			Def:   "ffffffff",
			Usage: "this is ipMask",
		},
		{
			Name:  "ipNet",
			Def:   "192.0.2.0/24",
			Usage: "this is ipNet",
		},
		{
			Name:      "env",
			Shorthand: "e",
			Def:       os.Getenv("PWD"),
			Usage:     "this is env",
		},
		{
			Name:  "others.string",
			Def:   "string",
			Usage: "this is others.string",
		},
	}

	for _, c := range cases {
		normalName := f.normalizeFlagName(c.Name)
		itemFlag := f.formal[normalName]
		err := checkFlag(itemFlag, c.Name, c.Shorthand, c.Def, c.Usage)
		if err != nil {
			t.Error(err)
		}
	}
}

func checkFlag(flag *Flag, name, short, def, desc string) error {
	if flag == nil {
		return fmt.Errorf("flag is nil")
	}
	if flag.Name != name {
		return fmt.Errorf("flag(%s) not equal flag.Name(%s) \n", name, flag.Name)
	}
	if flag.Shorthand != short {
		return fmt.Errorf("short(%s) not equal in flag.Shorthand(%s) \n", short, flag.Shorthand)
	}
	if flag.DefValue != def {
		return fmt.Errorf("def(%s) not equal in flag.DefValue(%s) \n", def, flag.DefValue)
	}
	if flag.Usage != desc {
		return fmt.Errorf("desc(%s) not equal in flag.Usage(%s) \n", desc, flag.Usage)
	}
	return nil
}
