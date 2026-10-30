// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestAllSliceTypes_Comprehensive(t *testing.T) {
	f := NewFlagSet("test_slices", ContinueOnError)
	f.SetOutput(io.Discard)

	var (
		bools  []bool
		i32s   []int32
		i64s   []int64
		u32s   []uint
		f32s   []float32
		f64s   []float64
		durs   []time.Duration
		ips    []net.IP
		ipnets []net.IPNet
	)

	// BoolSlice
	f.BoolSliceVar(&bools, "bools", []bool{false}, "usage")
	f.BoolSliceVarP(&bools, "bools-p", "a", []bool{false}, "usage")
	pBools := f.BoolSlice("bools-ptr", []bool{true}, "usage")
	pBoolsp := f.BoolSliceP("bools-ptr-p", "b", []bool{true}, "usage")

	// Int32Slice & Int64Slice
	f.Int32SliceVar(&i32s, "i32s", nil, "usage")
	f.Int32SliceVarP(&i32s, "i32s-p", "c", nil, "usage")
	pI32 := f.Int32Slice("i32-ptr", []int32{1}, "usage")
	pI32p := f.Int32SliceP("i32-ptr-p", "d", []int32{1}, "usage")

	f.Int64SliceVar(&i64s, "i64s", nil, "usage")
	f.Int64SliceVarP(&i64s, "i64s-p", "e", nil, "usage")
	pI64 := f.Int64Slice("i64-ptr", []int64{2}, "usage")
	pI64p := f.Int64SliceP("i64-ptr-p", "g", []int64{2}, "usage")

	// UintSlice
	f.UintSliceVar(&u32s, "uints", nil, "usage")
	f.UintSliceVarP(&u32s, "uints-p", "h", nil, "usage")
	pU := f.UintSlice("uint-ptr", []uint{3}, "usage")
	pUp := f.UintSliceP("uint-ptr-p", "i", []uint{3}, "usage")

	// Float32Slice & Float64Slice
	f.Float32SliceVar(&f32s, "f32s", nil, "usage")
	f.Float32SliceVarP(&f32s, "f32s-p", "j", nil, "usage")
	pF32 := f.Float32Slice("f32-ptr", []float32{1.5}, "usage")
	pF32p := f.Float32SliceP("f32-ptr-p", "k", []float32{1.5}, "usage")

	f.Float64SliceVar(&f64s, "f64s", nil, "usage")
	f.Float64SliceVarP(&f64s, "f64s-p", "l", nil, "usage")
	pF64 := f.Float64Slice("f64-ptr", []float64{2.5}, "usage")
	pF64p := f.Float64SliceP("f64-ptr-p", "m", []float64{2.5}, "usage")

	// DurationSlice
	f.DurationSliceVar(&durs, "durs", nil, "usage")
	f.DurationSliceVarP(&durs, "durs-p", "n", nil, "usage")
	pDurs := f.DurationSlice("durs-ptr", []time.Duration{time.Second}, "usage")
	pDursp := f.DurationSliceP("durs-ptr-p", "o", []time.Duration{time.Second}, "usage")

	// IPSlice & IPNetSlice
	f.IPSliceVar(&ips, "ips", nil, "usage")
	f.IPSliceVarP(&ips, "ips-p", "q", nil, "usage")
	pIPs := f.IPSlice("ips-ptr", []net.IP{net.ParseIP("127.0.0.1")}, "usage")
	pIPsp := f.IPSliceP("ips-ptr-p", "r", []net.IP{net.ParseIP("127.0.0.1")}, "usage")

	_, defaultNet, _ := net.ParseCIDR("192.168.1.0/24")
	f.IPNetSliceVar(&ipnets, "nets", nil, "usage")
	f.IPNetSliceVarP(&ipnets, "nets-p", "s", nil, "usage")
	pNets := f.IPNetSlice("nets-ptr", []net.IPNet{*defaultNet}, "usage")
	pNetsp := f.IPNetSliceP("nets-ptr-p", "u", []net.IPNet{*defaultNet}, "usage")

	args := []string{
		"--bools=true,false,true",
		"-b", "false,false",
		"--i32s=10,20,30",
		"-d", "40,50",
		"--i64s=100,200",
		"-g", "300",
		"--uints=1,2,3",
		"-i", "4,5",
		"--f32s=1.2,3.4",
		"-k", "5.6",
		"--f64s=10.1,20.2",
		"-m", "30.3",
		"--durs=1s,2m,3h",
		"-o", "4h",
		"--ips=1.1.1.1,8.8.8.8",
		"-r", "10.0.0.1",
		"--nets=10.0.0.0/8,172.16.0.0/12",
		"-u", "192.168.0.0/16",
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	// Verify values and getters
	bVal, err := f.GetBoolSlice("bools")
	if err != nil || len(bVal) != 3 || !bVal[0] || bVal[1] || !bVal[2] {
		t.Errorf("GetBoolSlice failed: %v, %v", bVal, err)
	}
	if len(*pBoolsp) != 2 || (*pBoolsp)[0] || (*pBoolsp)[1] {
		t.Errorf("pBoolsp mismatch: %v", *pBoolsp)
	}

	i32Val, err := f.GetInt32Slice("i32s")
	if err != nil || len(i32Val) != 3 || i32Val[0] != 10 {
		t.Errorf("GetInt32Slice failed: %v, %v", i32Val, err)
	}

	i64Val, err := f.GetInt64Slice("i64s")
	if err != nil || len(i64Val) != 2 || i64Val[1] != 200 {
		t.Errorf("GetInt64Slice failed: %v, %v", i64Val, err)
	}

	uVal, err := f.GetUintSlice("uints")
	if err != nil || len(uVal) != 3 || uVal[2] != 3 {
		t.Errorf("GetUintSlice failed: %v, %v", uVal, err)
	}

	f32Val, err := f.GetFloat32Slice("f32s")
	if err != nil || len(f32Val) != 2 {
		t.Errorf("GetFloat32Slice failed: %v, %v", f32Val, err)
	}

	f64Val, err := f.GetFloat64Slice("f64s")
	if err != nil || len(f64Val) != 2 {
		t.Errorf("GetFloat64Slice failed: %v, %v", f64Val, err)
	}

	durVal, err := f.GetDurationSlice("durs")
	if err != nil || len(durVal) != 3 || durVal[2] != 3*time.Hour {
		t.Errorf("GetDurationSlice failed: %v, %v", durVal, err)
	}

	ipVal, err := f.GetIPSlice("ips")
	if err != nil || len(ipVal) != 2 || ipVal[0].String() != "1.1.1.1" {
		t.Errorf("GetIPSlice failed: %v, %v", ipVal, err)
	}

	netVal, err := f.GetIPNetSlice("nets")
	if err != nil || len(netVal) != 2 || netVal[0].String() != "10.0.0.0/8" {
		t.Errorf("GetIPNetSlice failed: %v, %v", netVal, err)
	}

	// Verify pointers
	if len(*pBools) != 1 || len(*pI32) != 1 || len(*pI32p) != 2 || len(*pI64) != 1 || len(*pI64p) != 1 || len(*pU) != 1 || len(*pUp) != 2 || len(*pF32) != 1 || len(*pF32p) != 1 || len(*pF64) != 1 || len(*pF64p) != 1 || len(*pDurs) != 1 || len(*pDursp) != 1 || len(*pIPs) != 1 || len(*pIPsp) != 1 || len(*pNets) != 1 || len(*pNetsp) != 1 {
		t.Errorf("pointer slice lengths mismatch")
	}

	// Test Append, Replace, and GetSlice on SliceValue
	testSliceMethods(t)

	// Error tests on invalid slice items
	sliceErrTests := []struct {
		name string
		val  string
	}{
		{"bools", "not_a_bool"},
		{"i32s", "not_an_int"},
		{"i64s", "not_an_int"},
		{"uints", "-1"},
		{"f32s", "not_a_float"},
		{"f64s", "not_a_float"},
		{"durs", "not_a_duration"},
		{"ips", "not.an.ip.address.really"},
		{"nets", "not_a_cidr"},
	}

	for _, st := range sliceErrTests {
		if err := f.Set(st.name, st.val); err == nil {
			t.Errorf("expected error setting %s to %s", st.name, st.val)
		}
	}
}

func testSliceMethods(t *testing.T) {
	// boolSliceValue
	var bs []bool
	bv := newBoolSliceValue([]bool{true}, &bs)
	if err := bv.Append("false"); err != nil {
		t.Fatalf("bv.Append failed: %v", err)
	}
	if len(bs) != 2 || bs[1] != false {
		t.Errorf("bv.Append unexpected bs: %v", bs)
	}
	if err := bv.Replace([]string{"true", "true"}); err != nil {
		t.Fatalf("bv.Replace failed: %v", err)
	}
	if slice := bv.GetSlice(); len(slice) != 2 || slice[0] != "true" {
		t.Errorf("bv.GetSlice unexpected: %v", slice)
	}
	if err := bv.Append("invalid"); err == nil {
		t.Errorf("expected error from bv.Append('invalid')")
	}
	if err := bv.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected error from bv.Replace(['invalid'])")
	}

	// durationSliceValue
	var ds []time.Duration
	dv := newDurationSliceValue([]time.Duration{time.Second}, &ds)
	_ = dv.Append("1m")
	if len(ds) != 2 || ds[1] != time.Minute {
		t.Errorf("dv.Append unexpected: %v", ds)
	}
	_ = dv.Replace([]string{"2s", "3s"})
	if slice := dv.GetSlice(); len(slice) != 2 || slice[0] != "2s" {
		t.Errorf("dv.GetSlice unexpected: %v", slice)
	}
	if err := dv.Append("invalid"); err == nil {
		t.Errorf("expected error from dv.Append('invalid')")
	}
	if err := dv.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected error from dv.Replace(['invalid'])")
	}

	// float32SliceValue & float64SliceValue
	var f32s []float32
	f32v := newFloat32SliceValue(nil, &f32s)
	_ = f32v.Append("1.23")
	_ = f32v.Replace([]string{"4.56"})
	_ = f32v.GetSlice()
	if err := f32v.Append("invalid"); err == nil {
		t.Errorf("expected err f32v.Append")
	}
	if err := f32v.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected err f32v.Replace")
	}

	var f64s []float64
	f64v := newFloat64SliceValue(nil, &f64s)
	_ = f64v.Append("7.89")
	_ = f64v.Replace([]string{"0.12"})
	_ = f64v.GetSlice()
	if err := f64v.Append("invalid"); err == nil {
		t.Errorf("expected err f64v.Append")
	}
	if err := f64v.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected err f64v.Replace")
	}

	// int32SliceValue & int64SliceValue & uintSliceValue
	var i32s []int32
	i32v := newInt32SliceValue(nil, &i32s)
	_ = i32v.Append("10")
	_ = i32v.Replace([]string{"20"})
	_ = i32v.GetSlice()
	if err := i32v.Append("invalid"); err == nil {
		t.Errorf("expected err i32v.Append")
	}
	if err := i32v.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected err i32v.Replace")
	}

	var i64s []int64
	i64v := newInt64SliceValue(nil, &i64s)
	_ = i64v.Append("100")
	_ = i64v.Replace([]string{"200"})
	_ = i64v.GetSlice()
	if err := i64v.Append("invalid"); err == nil {
		t.Errorf("expected err i64v.Append")
	}
	if err := i64v.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected err i64v.Replace")
	}

	var us []uint
	uv := newUintSliceValue(nil, &us)
	_ = uv.Append("300")
	_ = uv.Replace([]string{"400"})
	_ = uv.GetSlice()
	if err := uv.Append("invalid"); err == nil {
		t.Errorf("expected err uv.Append")
	}
	if err := uv.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected err uv.Replace")
	}

	// ipSliceValue
	var ips []net.IP
	ipv := newIPSliceValue(nil, &ips)
	_ = ipv.Append("1.1.1.1")
	_ = ipv.Replace([]string{"2.2.2.2"})
	_ = ipv.GetSlice()
	if err := ipv.Append("invalid"); err == nil {
		t.Errorf("expected err ipv.Append")
	}
	if err := ipv.Replace([]string{"invalid"}); err == nil {
		t.Errorf("expected err ipv.Replace")
	}

	// ipNetSliceValue
	var ipnets []net.IPNet
	netv := newIPNetSliceValue(nil, &ipnets)
	_ = netv.Set("10.0.0.0/8, 192.168.0.0/16")
	_ = netv.String()
	_ = netv.Type()
	if err := netv.Set("invalid_cidr"); err == nil {
		t.Errorf("expected err netv.Set(invalid)")
	}
}

func TestAllMapTypes_Comprehensive(t *testing.T) {
	f := NewFlagSet("test_maps", ContinueOnError)
	f.SetOutput(io.Discard)

	var (
		s2i  map[string]int
		s2i6 map[string]int64
		s2s  map[string]string
	)

	f.StringToIntVar(&s2i, "s2i", nil, "usage")
	f.StringToIntVarP(&s2i, "s2i-p", "a", nil, "usage")
	pS2i := f.StringToInt("s2i-ptr", map[string]int{"k": 1}, "usage")
	pS2ip := f.StringToIntP("s2i-ptr-p", "b", map[string]int{"k": 1}, "usage")

	f.StringToInt64Var(&s2i6, "s2i6", nil, "usage")
	f.StringToInt64VarP(&s2i6, "s2i6-p", "c", nil, "usage")
	pS2i6 := f.StringToInt64("s2i6-ptr", map[string]int64{"k6": 2}, "usage")
	pS2i6p := f.StringToInt64P("s2i6-ptr-p", "d", map[string]int64{"k6": 2}, "usage")

	f.StringToStringVar(&s2s, "s2s", nil, "usage")
	f.StringToStringVarP(&s2s, "s2s-p", "e", nil, "usage")
	pS2s := f.StringToString("s2s-ptr", map[string]string{"k": "v"}, "usage")
	pS2sp := f.StringToStringP("s2s-ptr-p", "g", map[string]string{"k": "v"}, "usage")

	args := []string{
		"--s2i=apple=5,banana=10",
		"-b", "cherry=15",
		"--s2i6=x=1000,y=2000",
		"-d", "z=3000",
		"--s2s=env=prod,region=us",
		"-g", "cluster=east",
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	// Verify values and getters
	m1, err := f.GetStringToInt("s2i")
	if err != nil || m1["apple"] != 5 || m1["banana"] != 10 {
		t.Errorf("GetStringToInt failed: %v, %v", m1, err)
	}
	if (*pS2ip)["cherry"] != 15 || (*pS2i)["k"] != 1 {
		t.Errorf("pS2i mismatch")
	}

	m2, err := f.GetStringToInt64("s2i6")
	if err != nil || m2["x"] != 1000 || m2["y"] != 2000 {
		t.Errorf("GetStringToInt64 failed: %v, %v", m2, err)
	}
	if (*pS2i6p)["z"] != 3000 || (*pS2i6)["k6"] != 2 {
		t.Errorf("pS2i6 mismatch")
	}

	m3, err := f.GetStringToString("s2s")
	if err != nil || m3["env"] != "prod" || m3["region"] != "us" {
		t.Errorf("GetStringToString failed: %v, %v", m3, err)
	}
	if (*pS2sp)["cluster"] != "east" || (*pS2s)["k"] != "v" {
		t.Errorf("pS2s mismatch")
	}

	// Error paths
	if err := f.Set("s2i", "no_equal_sign"); err == nil {
		t.Errorf("expected error setting s2i without equal sign")
	}
	if err := f.Set("s2i", "key=not_an_int"); err == nil {
		t.Errorf("expected error setting s2i with non-int value")
	}
	if err := f.Set("s2i6", "no_equal_sign"); err == nil {
		t.Errorf("expected error setting s2i6 without equal sign")
	}
	if err := f.Set("s2i6", "key=not_an_int"); err == nil {
		t.Errorf("expected error setting s2i6 with non-int value")
	}
	if err := f.Set("s2s", "no_equal_sign"); err == nil {
		t.Errorf("expected error setting s2s without equal sign")
	}
}

func TestIP_And_Func_Comprehensive(t *testing.T) {
	f := NewFlagSet("test_ip_func", ContinueOnError)
	f.SetOutput(io.Discard)

	// IP
	var ip net.IP
	defaultIP := net.ParseIP("127.0.0.1")
	f.IPVar(&ip, "ip", defaultIP, "usage")
	f.IPVarP(&ip, "ip-p", "i", defaultIP, "usage")
	pIP := f.IP("ip-ptr", defaultIP, "usage")
	pIPp := f.IPP("ip-ptr-p", "j", defaultIP, "usage")

	// IPNet
	var ipnet net.IPNet
	_, defaultCIDR, _ := net.ParseCIDR("10.0.0.0/8")
	f.IPNetVar(&ipnet, "cidr", *defaultCIDR, "usage")
	f.IPNetVarP(&ipnet, "cidr-p", "c", *defaultCIDR, "usage")
	pCIDR := f.IPNet("cidr-ptr", *defaultCIDR, "usage")
	pCIDRp := f.IPNetP("cidr-ptr-p", "d", *defaultCIDR, "usage")

	// BoolFunc & Func
	var boolFuncVal string
	f.BoolFunc("boolfunc", "usage", func(s string) error {
		boolFuncVal = s
		return nil
	})
	f.BoolFuncP("boolfunc-p", "b", "usage", func(s string) error {
		boolFuncVal = "p:" + s
		return nil
	})

	var funcVal string
	f.Func("customfunc", "usage", func(s string) error {
		funcVal = s
		return nil
	})
	f.FuncP("customfunc-p", "F", "usage", func(s string) error {
		funcVal = "p:" + s
		return nil
	})

	args := []string{
		"--ip=1.1.1.1",
		"-j", "8.8.8.8",
		"--cidr=172.16.0.0/12",
		"-d", "192.168.1.0/24",
		"-b",
		"-F", "hello_func",
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if ip.String() != "1.1.1.1" || pIPp.String() != "8.8.8.8" || pIP.String() != "127.0.0.1" {
		t.Errorf("IP mismatch: ip=%s, pIPp=%s", ip, pIPp)
	}
	getIP, err := f.GetIP("ip")
	if err != nil || getIP.String() != "1.1.1.1" {
		t.Errorf("GetIP failed: %v, %v", getIP, err)
	}

	if ipnet.String() != "172.16.0.0/12" || pCIDRp.String() != "192.168.1.0/24" || pCIDR.String() != "10.0.0.0/8" {
		t.Errorf("IPNet mismatch: ipnet=%s, pCIDRp=%s", ipnet.String(), pCIDRp.String())
	}
	getCIDR, err := f.GetIPNet("cidr")
	if err != nil || getCIDR.String() != "172.16.0.0/12" {
		t.Errorf("GetIPNet failed: %v, %v", getCIDR, err)
	}

	if boolFuncVal != "p:true" {
		t.Errorf("boolFuncVal = %q, want p:true", boolFuncVal)
	}
	if funcVal != "p:hello_func" {
		t.Errorf("funcVal = %q, want p:hello_func", funcVal)
	}

	// Error paths
	if err := f.Set("ip", "invalid.ip"); err == nil {
		t.Errorf("expected error setting invalid IP")
	}
	if err := f.Set("cidr", "invalid.cidr"); err == nil {
		t.Errorf("expected error setting invalid CIDR")
	}
}

func TestPackageLevel_SlicesAndMaps(t *testing.T) {
	oldCommandLine := CommandLine
	defer func() { CommandLine = oldCommandLine }()

	CommandLine = NewFlagSet("cmd_slices_maps", ContinueOnError)
	CommandLine.SetOutput(io.Discard)

	var (
		bs   []bool
		i32s []int32
		i64s []int64
		u    []uint
		f32s []float32
		f64s []float64
		ds   []time.Duration
		ips  []net.IP
		nets []net.IPNet
		s2i  map[string]int
		s2i6 map[string]int64
		s2s  map[string]string
		ip   net.IP
		cidr net.IPNet
	)

	BoolSliceVar(&bs, "p-bs", nil, "usage")
	BoolSliceVarP(&bs, "p-bs-p", "a", nil, "usage")
	_ = BoolSlice("pkg-bs", nil, "usage")
	_ = BoolSliceP("pkg-bs-p", "A", nil, "usage")

	Int32SliceVar(&i32s, "p-i32s", nil, "usage")
	Int32SliceVarP(&i32s, "p-i32s-p", "b", nil, "usage")
	_ = Int32Slice("pkg-i32s", nil, "usage")
	_ = Int32SliceP("pkg-i32s-p", "B", nil, "usage")

	Int64SliceVar(&i64s, "p-i64s", nil, "usage")
	Int64SliceVarP(&i64s, "p-i64s-p", "c", nil, "usage")
	_ = Int64Slice("pkg-i64s", nil, "usage")
	_ = Int64SliceP("pkg-i64s-p", "C", nil, "usage")

	UintSliceVar(&u, "p-u", nil, "usage")
	UintSliceVarP(&u, "p-u-p", "d", nil, "usage")
	_ = UintSlice("pkg-u", nil, "usage")
	_ = UintSliceP("pkg-u-p", "D", nil, "usage")

	Float32SliceVar(&f32s, "p-f32s", nil, "usage")
	Float32SliceVarP(&f32s, "p-f32s-p", "e", nil, "usage")
	_ = Float32Slice("pkg-f32s", nil, "usage")
	_ = Float32SliceP("pkg-f32s-p", "E", nil, "usage")

	Float64SliceVar(&f64s, "p-f64s", nil, "usage")
	Float64SliceVarP(&f64s, "p-f64s-p", "g", nil, "usage")
	_ = Float64Slice("pkg-f64s", nil, "usage")
	_ = Float64SliceP("pkg-f64s-p", "G", nil, "usage")

	DurationSliceVar(&ds, "p-ds", nil, "usage")
	DurationSliceVarP(&ds, "p-ds-p", "h", nil, "usage")
	_ = DurationSlice("pkg-ds", nil, "usage")
	_ = DurationSliceP("pkg-ds-p", "H", nil, "usage")

	IPSliceVar(&ips, "p-ips", nil, "usage")
	IPSliceVarP(&ips, "p-ips-p", "i", nil, "usage")
	_ = IPSlice("pkg-ips", nil, "usage")
	_ = IPSliceP("pkg-ips-p", "I", nil, "usage")

	_, defaultCIDR, _ := net.ParseCIDR("10.0.0.0/8")
	IPNetSliceVar(&nets, "p-nets", nil, "usage")
	IPNetSliceVarP(&nets, "p-nets-p", "j", nil, "usage")
	_ = IPNetSlice("pkg-nets", nil, "usage")
	_ = IPNetSliceP("pkg-nets-p", "J", nil, "usage")

	StringToIntVar(&s2i, "p-s2i", nil, "usage")
	StringToIntVarP(&s2i, "p-s2i-p", "k", nil, "usage")
	_ = StringToInt("pkg-s2i", nil, "usage")
	_ = StringToIntP("pkg-s2i-p", "K", nil, "usage")

	StringToInt64Var(&s2i6, "p-s2i6", nil, "usage")
	StringToInt64VarP(&s2i6, "p-s2i6-p", "l", nil, "usage")
	_ = StringToInt64("pkg-s2i6", nil, "usage")
	_ = StringToInt64P("pkg-s2i6-p", "L", nil, "usage")

	StringToStringVar(&s2s, "p-s2s", nil, "usage")
	StringToStringVarP(&s2s, "p-s2s-p", "m", nil, "usage")
	_ = StringToString("pkg-s2s", nil, "usage")
	_ = StringToStringP("pkg-s2s-p", "M", nil, "usage")

	IPVar(&ip, "p-ip", net.ParseIP("127.0.0.1"), "usage")
	IPVarP(&ip, "p-ip-p", "n", net.ParseIP("127.0.0.1"), "usage")
	_ = IP("pkg-ip", net.ParseIP("127.0.0.1"), "usage")
	_ = IPP("pkg-ip-p", "N", net.ParseIP("127.0.0.1"), "usage")

	IPNetVar(&cidr, "p-cidr", *defaultCIDR, "usage")
	IPNetVarP(&cidr, "p-cidr-p", "o", *defaultCIDR, "usage")
	_ = IPNet("pkg-cidr", *defaultCIDR, "usage")
	_ = IPNetP("pkg-cidr-p", "O", *defaultCIDR, "usage")

	BoolFunc("pkg-boolfunc", "usage", func(string) error { return nil })
	BoolFuncP("pkg-boolfunc-p", "p", "usage", func(string) error { return nil })
	Func("pkg-func", "usage", func(string) error { return nil })
	FuncP("pkg-func-p", "q", "usage", func(string) error { return nil })

	err := CommandLine.Parse([]string{
		"-a", "true",
		"-b", "100",
		"-c", "200",
		"-k", "k1=1",
		"-l", "k2=2",
		"-m", "k3=v3",
	})
	if err != nil {
		t.Fatalf("CommandLine.Parse failed: %v", err)
	}

	if !reflect.DeepEqual(bs, []bool{true}) || !reflect.DeepEqual(i32s, []int32{100}) || !reflect.DeepEqual(i64s, []int64{200}) {
		t.Errorf("unexpected package slice values: bs=%v, i32s=%v, i64s=%v", bs, i32s, i64s)
	}
	if s2i["k1"] != 1 || s2i6["k2"] != 2 || s2s["k3"] != "v3" {
		t.Errorf("unexpected package map values: s2i=%v, s2i6=%v, s2s=%v", s2i, s2i6, s2s)
	}
}

func TestRemainingTypes_Comprehensive(t *testing.T) {
	f := NewFlagSet("test_remaining", ContinueOnError)
	f.SetOutput(io.Discard)

	var (
		is  []int
		sa  []string
		ss  []string
		tip net.IP
	)

	f.IntSliceVar(&is, "is", nil, "usage")
	f.IntSliceVarP(&is, "is-p", "i", nil, "usage")
	pIs := f.IntSlice("is-ptr", []int{1}, "usage")
	pIsp := f.IntSliceP("is-ptr-p", "I", []int{1}, "usage")

	f.StringArrayVar(&sa, "sa", nil, "usage")
	f.StringArrayVarP(&sa, "sa-p", "a", nil, "usage")
	pSa := f.StringArray("sa-ptr", []string{"a"}, "usage")
	pSap := f.StringArrayP("sa-ptr-p", "A", []string{"a"}, "usage")

	f.StringSliceVar(&ss, "ss", nil, "usage")
	f.StringSliceVarP(&ss, "ss-p", "s", nil, "usage")
	pSs := f.StringSlice("ss-ptr", []string{"s"}, "usage")
	pSsp := f.StringSliceP("ss-ptr-p", "S", []string{"s"}, "usage")

	defIP := net.ParseIP("10.0.0.1")
	f.TextVar(&tip, "text-ip", &defIP, "usage")
	f.TextVarP(&tip, "text-ip-p", "t", &defIP, "usage")

	args := []string{
		"--is=1,2,3",
		"-I", "4,5",
		"--sa=hello", "--sa=world",
		"-A", "foo",
		"--ss=cat,dog",
		"-S", "bird",
		"-t", "192.168.1.1",
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if len(is) != 3 || len(*pIsp) != 2 || len(*pIs) != 1 {
		t.Errorf("IntSlice mismatch")
	}
	if len(sa) != 2 || len(*pSap) != 1 || len(*pSa) != 1 {
		t.Errorf("StringArray mismatch")
	}
	if len(ss) != 2 || len(*pSsp) != 1 || len(*pSs) != 1 {
		t.Errorf("StringSlice mismatch")
	}
	if tip.String() != "192.168.1.1" {
		t.Errorf("tip = %s, want 192.168.1.1", tip.String())
	}

	var outIP net.IP
	if err := f.GetText("text-ip-p", &outIP); err != nil || outIP.String() != "192.168.1.1" {
		t.Errorf("GetText failed: %v, %v", outIP, err)
	}

	// Package-level helpers
	oldCommandLine := CommandLine
	defer func() { CommandLine = oldCommandLine }()

	CommandLine = NewFlagSet("cmd_rem", ContinueOnError)
	CommandLine.SetOutput(io.Discard)

	var (
		pkgIs  []int
		pkgSa  []string
		pkgSs  []string
		pkgTip net.IP
	)

	IntSliceVar(&pkgIs, "p-is", nil, "usage")
	IntSliceVarP(&pkgIs, "p-is-p", "i", nil, "usage")
	_ = IntSlice("pkg-is", nil, "usage")
	_ = IntSliceP("pkg-is-p", "I", nil, "usage")

	StringArrayVar(&pkgSa, "p-sa", nil, "usage")
	StringArrayVarP(&pkgSa, "p-sa-p", "a", nil, "usage")
	_ = StringArray("pkg-sa", nil, "usage")
	_ = StringArrayP("pkg-sa-p", "A", nil, "usage")

	StringSliceVar(&pkgSs, "p-ss", nil, "usage")
	StringSliceVarP(&pkgSs, "p-ss-p", "s", nil, "usage")
	_ = StringSlice("pkg-ss", nil, "usage")
	_ = StringSliceP("pkg-ss-p", "S", nil, "usage")

	TextVar(&pkgTip, "p-text", &defIP, "usage")
	TextVarP(&pkgTip, "p-text-p", "t", &defIP, "usage")

	err := CommandLine.Parse([]string{"-i", "10,20", "-a", "one", "-s", "two", "-t", "8.8.8.8"})
	if err != nil {
		t.Fatalf("CommandLine.Parse failed: %v", err)
	}
	if len(pkgIs) != 2 || len(pkgSa) != 1 || len(pkgSs) != 1 || pkgTip.String() != "8.8.8.8" {
		t.Errorf("package level remaining types mismatch")
	}
}
