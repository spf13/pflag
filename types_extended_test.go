// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pflag

import (
	"encoding/hex"
	"io"
	"net"
	"testing"
	"time"
)

func TestAllNumericTypes_FlagSet(t *testing.T) {
	f := NewFlagSet("test_num", ContinueOnError)
	f.SetOutput(io.Discard)

	var (
		i8  int8
		i16 int16
		i32 int32
		i64 int64
		iv  int
		u8  uint8
		u16 uint16
		u32 uint32
		u64 uint64
		uv  uint
		f32 float32
		f64 float64
		dur time.Duration
	)

	// Var & VarP
	f.Int8Var(&i8, "int8", 8, "int8 usage")
	f.Int8VarP(&i8, "int8-p", "a", 8, "int8 usage")
	f.Int16Var(&i16, "int16", 16, "int16 usage")
	f.Int16VarP(&i16, "int16-p", "b", 16, "int16 usage")
	f.Int32Var(&i32, "int32", 32, "int32 usage")
	f.Int32VarP(&i32, "int32-p", "c", 32, "int32 usage")
	f.Int64Var(&i64, "int64", 64, "int64 usage")
	f.Int64VarP(&i64, "int64-p", "d", 64, "int64 usage")
	f.IntVar(&iv, "int", 1, "int usage")
	f.IntVarP(&iv, "int-p", "e", 1, "int usage")

	f.Uint8Var(&u8, "uint8", 8, "uint8 usage")
	f.Uint8VarP(&u8, "uint8-p", "g", 8, "uint8 usage")
	f.Uint16Var(&u16, "uint16", 16, "uint16 usage")
	f.Uint16VarP(&u16, "uint16-p", "h", 16, "uint16 usage")
	f.Uint32Var(&u32, "uint32", 32, "uint32 usage")
	f.Uint32VarP(&u32, "uint32-p", "i", 32, "uint32 usage")
	f.Uint64Var(&u64, "uint64", 64, "uint64 usage")
	f.Uint64VarP(&u64, "uint64-p", "j", 64, "uint64 usage")
	f.UintVar(&uv, "uint", 1, "uint usage")
	f.UintVarP(&uv, "uint-p", "k", 1, "uint usage")

	f.Float32Var(&f32, "float32", 3.2, "float32 usage")
	f.Float32VarP(&f32, "float32-p", "l", 3.2, "float32 usage")
	f.Float64Var(&f64, "float64", 6.4, "float64 usage")
	f.Float64VarP(&f64, "float64-p", "m", 6.4, "float64 usage")

	f.DurationVar(&dur, "duration", time.Second, "duration usage")
	f.DurationVarP(&dur, "duration-p", "n", time.Second, "duration usage")

	// Pointers
	pi8 := f.Int8("int8-ptr", 1, "usage")
	pi8p := f.Int8P("int8-ptr-p", "o", 1, "usage")
	pi16 := f.Int16("int16-ptr", 2, "usage")
	pi16p := f.Int16P("int16-ptr-p", "q", 2, "usage")
	pi32 := f.Int32("int32-ptr", 3, "usage")
	pi32p := f.Int32P("int32-ptr-p", "r", 3, "usage")
	pi64 := f.Int64("int64-ptr", 4, "usage")
	pi64p := f.Int64P("int64-ptr-p", "s", 4, "usage")
	piv := f.Int("int-ptr", 5, "usage")
	pivp := f.IntP("int-ptr-p", "v", 5, "usage")

	pu8 := f.Uint8("uint8-ptr", 6, "usage")
	pu8p := f.Uint8P("uint8-ptr-p", "w", 6, "usage")
	pu16 := f.Uint16("uint16-ptr", 7, "usage")
	pu16p := f.Uint16P("uint16-ptr-p", "x", 7, "usage")
	pu32 := f.Uint32("uint32-ptr", 8, "usage")
	pu32p := f.Uint32P("uint32-ptr-p", "y", 8, "usage")
	pu64 := f.Uint64("uint64-ptr", 9, "usage")
	pu64p := f.Uint64P("uint64-ptr-p", "z", 9, "usage")
	puv := f.Uint("uint-ptr", 10, "usage")
	puvp := f.UintP("uint-ptr-p", "A", 10, "usage")

	pf32 := f.Float32("float32-ptr", 1.1, "usage")
	pf32p := f.Float32P("float32-ptr-p", "B", 1.1, "usage")
	pf64 := f.Float64("float64-ptr", 2.2, "usage")
	pf64p := f.Float64P("float64-ptr-p", "C", 2.2, "usage")
	pdur := f.Duration("duration-ptr", time.Minute, "usage")
	pdurp := f.DurationP("duration-ptr-p", "D", time.Minute, "usage")

	args := []string{
		"--int8=12", "-b", "24",
		"--int32=36", "-d", "48",
		"--int=50",
		"--uint8=80", "-h", "90",
		"--uint32=100", "-j", "110",
		"--uint=120",
		"--float32=12.5", "-m", "25.5",
		"--duration=2h45m",
		"-o", "7", "-q", "8", "-r", "9", "-s", "10", "-v", "11",
		"-w", "12", "-x", "13", "-y", "14", "-z", "15", "-A", "16",
		"-B", "3.3", "-C", "4.4", "-D", "10s",
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	// Verify getters
	if v, err := f.GetInt8("int8"); err != nil || v != 12 {
		t.Errorf("GetInt8('int8') = %v, %v; want 12, nil", v, err)
	}
	if v, err := f.GetInt16("int16-p"); err != nil || v != 24 {
		t.Errorf("GetInt16('int16-p') = %v, %v; want 24, nil", v, err)
	}
	if v, err := f.GetInt32("int32"); err != nil || v != 36 {
		t.Errorf("GetInt32('int32') = %v, %v; want 36, nil", v, err)
	}
	if v, err := f.GetInt64("int64-p"); err != nil || v != 48 {
		t.Errorf("GetInt64('int64-p') = %v, %v; want 48, nil", v, err)
	}
	if v, err := f.GetInt("int"); err != nil || v != 50 {
		t.Errorf("GetInt('int') = %v, %v; want 50, nil", v, err)
	}

	if v, err := f.GetUint8("uint8"); err != nil || v != 80 {
		t.Errorf("GetUint8('uint8') = %v, %v; want 80, nil", v, err)
	}
	if v, err := f.GetUint16("uint16-p"); err != nil || v != 90 {
		t.Errorf("GetUint16('uint16-p') = %v, %v; want 90, nil", v, err)
	}
	if v, err := f.GetUint32("uint32"); err != nil || v != 100 {
		t.Errorf("GetUint32('uint32') = %v, %v; want 100, nil", v, err)
	}
	if v, err := f.GetUint64("uint64-p"); err != nil || v != 110 {
		t.Errorf("GetUint64('uint64-p') = %v, %v; want 110, nil", v, err)
	}
	if v, err := f.GetUint("uint"); err != nil || v != 120 {
		t.Errorf("GetUint('uint') = %v, %v; want 120, nil", v, err)
	}

	if v, err := f.GetFloat32("float32"); err != nil || v != 12.5 {
		t.Errorf("GetFloat32('float32') = %v, %v; want 12.5, nil", v, err)
	}
	if v, err := f.GetFloat64("float64-p"); err != nil || v != 25.5 {
		t.Errorf("GetFloat64('float64-p') = %v, %v; want 25.5, nil", v, err)
	}
	if v, err := f.GetDuration("duration"); err != nil || v != 2*time.Hour+45*time.Minute {
		t.Errorf("GetDuration('duration') = %v, %v; want 2h45m, nil", v, err)
	}

	// Verify pointer dereferences
	if *pi8 != 1 || *pi8p != 7 || *pi16 != 2 || *pi16p != 8 || *pi32 != 3 || *pi32p != 9 || *pi64 != 4 || *pi64p != 10 || *piv != 5 || *pivp != 11 {
		t.Errorf("unexpected signed int pointer values")
	}
	if *pu8 != 6 || *pu8p != 12 || *pu16 != 7 || *pu16p != 13 || *pu32 != 8 || *pu32p != 14 || *pu64 != 9 || *pu64p != 15 || *puv != 10 || *puvp != 16 {
		t.Errorf("unexpected unsigned int pointer values")
	}
	if *pf32 != 1.1 || *pf32p != 3.3 || *pf64 != 2.2 || *pf64p != 4.4 || *pdur != time.Minute || *pdurp != 10*time.Second {
		t.Errorf("unexpected float/duration pointer values")
	}

	// Test Error paths on setters
	errTests := []struct {
		name string
		val  string
	}{
		{"int8", "invalid"},
		{"int8", "128"}, // overflow
		{"int16", "invalid"},
		{"int32", "invalid"},
		{"int64", "invalid"},
		{"int", "invalid"},
		{"uint8", "invalid"},
		{"uint8", "256"}, // overflow
		{"uint16", "invalid"},
		{"uint32", "invalid"},
		{"uint64", "invalid"},
		{"uint", "invalid"},
		{"float32", "invalid"},
		{"float64", "invalid"},
		{"duration", "invalid"},
	}

	for _, et := range errTests {
		if err := f.Set(et.name, et.val); err == nil {
			t.Errorf("expected error setting %s to %s", et.name, et.val)
		}
	}
}

func TestPackageLevel_NumericAndBasic(t *testing.T) {
	oldCommandLine := CommandLine
	defer func() { CommandLine = oldCommandLine }()

	CommandLine = NewFlagSet("cmd_num", ContinueOnError)
	CommandLine.SetOutput(io.Discard)

	var (
		i8  int8
		i16 int16
		i32 int32
		i64 int64
		iv  int
		u8  uint8
		u16 uint16
		u32 uint32
		u64 uint64
		uv  uint
		f32 float32
		f64 float64
		dur time.Duration
		str string
	)

	Int8Var(&i8, "p-int8", 1, "usage")
	Int8VarP(&i8, "p-int8-p", "a", 1, "usage")
	Int16Var(&i16, "p-int16", 2, "usage")
	Int16VarP(&i16, "p-int16-p", "b", 2, "usage")
	Int32Var(&i32, "p-int32", 3, "usage")
	Int32VarP(&i32, "p-int32-p", "c", 3, "usage")
	Int64Var(&i64, "p-int64", 4, "usage")
	Int64VarP(&i64, "p-int64-p", "d", 4, "usage")
	IntVar(&iv, "p-int", 5, "usage")
	IntVarP(&iv, "p-int-p", "e", 5, "usage")

	Uint8Var(&u8, "p-uint8", 6, "usage")
	Uint8VarP(&u8, "p-uint8-p", "f", 6, "usage")
	Uint16Var(&u16, "p-uint16", 7, "usage")
	Uint16VarP(&u16, "p-uint16-p", "g", 7, "usage")
	Uint32Var(&u32, "p-uint32", 8, "usage")
	Uint32VarP(&u32, "p-uint32-p", "h", 8, "usage")
	Uint64Var(&u64, "p-uint64", 9, "usage")
	Uint64VarP(&u64, "p-uint64-p", "i", 9, "usage")
	UintVar(&uv, "p-uint", 10, "usage")
	UintVarP(&uv, "p-uint-p", "j", 10, "usage")

	Float32Var(&f32, "p-float32", 1.5, "usage")
	Float32VarP(&f32, "p-float32-p", "k", 1.5, "usage")
	Float64Var(&f64, "p-float64", 2.5, "usage")
	Float64VarP(&f64, "p-float64-p", "l", 2.5, "usage")
	DurationVar(&dur, "p-duration", time.Second, "usage")
	DurationVarP(&dur, "p-duration-p", "m", time.Second, "usage")
	StringVar(&str, "p-str", "def", "usage")
	StringVarP(&str, "p-str-p", "n", "def", "usage")

	// Package level pointers
	_ = Int8("pkg-int8", 1, "usage")
	_ = Int8P("pkg-int8-p", "o", 1, "usage")
	_ = Int16("pkg-int16", 2, "usage")
	_ = Int16P("pkg-int16-p", "p", 2, "usage")
	_ = Int32("pkg-int32", 3, "usage")
	_ = Int32P("pkg-int32-p", "q", 3, "usage")
	_ = Int64("pkg-int64", 4, "usage")
	_ = Int64P("pkg-int64-p", "r", 4, "usage")
	_ = Int("pkg-int", 5, "usage")
	_ = IntP("pkg-int-p", "s", 5, "usage")

	_ = Uint8("pkg-uint8", 6, "usage")
	_ = Uint8P("pkg-uint8-p", "t", 6, "usage")
	_ = Uint16("pkg-uint16", 7, "usage")
	_ = Uint16P("pkg-uint16-p", "u", 7, "usage")
	_ = Uint32("pkg-uint32", 8, "usage")
	_ = Uint32P("pkg-uint32-p", "v", 8, "usage")
	_ = Uint64("pkg-uint64", 9, "usage")
	_ = Uint64P("pkg-uint64-p", "w", 9, "usage")
	_ = Uint("pkg-uint", 10, "usage")
	_ = UintP("pkg-uint-p", "x", 10, "usage")

	_ = Float32("pkg-float32", 1.1, "usage")
	_ = Float32P("pkg-float32-p", "y", 1.1, "usage")
	_ = Float64("pkg-float64", 2.2, "usage")
	_ = Float64P("pkg-float64-p", "z", 2.2, "usage")
	_ = Duration("pkg-duration", time.Minute, "usage")
	_ = DurationP("pkg-duration-p", "A", time.Minute, "usage")
	_ = String("pkg-string", "val", "usage")
	_ = StringP("pkg-string-p", "B", "val", "usage")

	err := CommandLine.Parse([]string{"--p-int8=100", "-n", "hello"})
	if err != nil {
		t.Fatalf("CommandLine.Parse failed: %v", err)
	}
	if i8 != 100 || str != "hello" {
		t.Errorf("unexpected parsed values: i8=%d, str=%s", i8, str)
	}
}

func TestSpecialTypes_FlagSet_And_PackageLevel(t *testing.T) {
	f := NewFlagSet("test_special", ContinueOnError)
	f.SetOutput(io.Discard)

	// Count
	var cnt int
	f.CountVar(&cnt, "verbose", "verbose count")
	f.CountVarP(&cnt, "verbose-p", "v", "verbose count")
	pcnt := f.Count("cnt", "count ptr")
	pcntp := f.CountP("cnt-p", "c", "count ptr")

	// BytesHex & BytesBase64
	var bHex, bB64 []byte
	f.BytesHexVar(&bHex, "hex", []byte{0x01}, "hex usage")
	f.BytesHexVarP(&bHex, "hex-p", "H", []byte{0x01}, "hex usage")
	pHex := f.BytesHex("hex-ptr", []byte{0x02}, "usage")
	pHexp := f.BytesHexP("hex-ptr-p", "I", []byte{0x02}, "usage")

	f.BytesBase64Var(&bB64, "b64", []byte("a"), "b64 usage")
	f.BytesBase64VarP(&bB64, "b64-p", "J", []byte("a"), "b64 usage")
	pB64 := f.BytesBase64("b64-ptr", []byte("b"), "usage")
	pB64p := f.BytesBase64P("b64-ptr-p", "K", []byte("b"), "usage")

	// IPMask
	var mask net.IPMask
	defaultMask := net.IPv4Mask(255, 255, 255, 0)
	f.IPMaskVar(&mask, "mask", defaultMask, "mask usage")
	f.IPMaskVarP(&mask, "mask-p", "M", defaultMask, "mask usage")
	pMask := f.IPMask("mask-ptr", defaultMask, "usage")
	pMaskp := f.IPMaskP("mask-ptr-p", "N", defaultMask, "usage")

	// Time
	var tm time.Time
	formats := []string{time.RFC3339, "2006-01-02"}
	defaultTime, _ := time.Parse("2006-01-02", "2026-01-01")
	f.TimeVar(&tm, "time", defaultTime, formats, "time usage")
	f.TimeVarP(&tm, "time-p", "T", defaultTime, formats, "time usage")
	pTime := f.Time("time-ptr", defaultTime, formats, "usage")
	pTimep := f.TimeP("time-ptr-p", "U", defaultTime, formats, "usage")

	args := []string{
		"-v", "-v", "-v",
		"-c", "-c",
		"--hex=deadbeef",
		"-I", "cafebabe",
		"--b64=aGVsbG8=",
		"-K", "d29ybGQ=",
		"--mask=255.255.0.0",
		"-N", "ffff0000",
		"--time=2026-10-04T12:00:00Z",
		"-U", "2026-12-31",
	}

	if err := f.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	// Verify Count
	if cnt != 3 || *pcnt != 0 || *pcntp != 2 {
		t.Errorf("Count values: cnt=%d, *pcnt=%d, *pcntp=%d", cnt, *pcnt, *pcntp)
	}
	if v, err := f.GetCount("verbose"); err != nil || v != 3 {
		t.Errorf("GetCount = %v, %v", v, err)
	}

	// Verify BytesHex
	expectedHex, _ := hex.DecodeString("deadbeef")
	expectedHexPtr, _ := hex.DecodeString("cafebabe")
	if string(bHex) != string(expectedHex) || string(*pHexp) != string(expectedHexPtr) || string(*pHex) != string([]byte{0x02}) {
		t.Errorf("BytesHex mismatch")
	}
	if v, err := f.GetBytesHex("hex"); err != nil || string(v) != string(expectedHex) {
		t.Errorf("GetBytesHex = %v, %v", v, err)
	}

	// Verify BytesBase64
	if string(bB64) != "hello" || string(*pB64p) != "world" || string(*pB64) != "b" {
		t.Errorf("BytesBase64 mismatch")
	}
	if v, err := f.GetBytesBase64("b64"); err != nil || string(v) != "hello" {
		t.Errorf("GetBytesBase64 = %v, %v", v, err)
	}

	// Verify IPMask
	if mask.String() != "ffff0000" || pMaskp.String() != "ffff0000" || pMask.String() != "ffffff00" {
		t.Errorf("IPMask mismatch")
	}
	if v, err := f.GetIPv4Mask("mask"); err != nil || v.String() != "ffff0000" {
		t.Errorf("GetIPv4Mask = %v, %v", v, err)
	}

	// Verify Time
	if tm.Format(time.RFC3339) != "2026-10-04T12:00:00Z" || pTimep.Format("2006-01-02") != "2026-12-31" || pTime.Format("2006-01-02") != "2026-01-01" {
		t.Errorf("Time mismatch")
	}
	if v, err := f.GetTime("time"); err != nil || v.Format(time.RFC3339) != "2026-10-04T12:00:00Z" {
		t.Errorf("GetTime = %v, %v", v, err)
	}

	// Error paths
	if err := f.Set("hex", "not-a-hex"); err == nil {
		t.Errorf("expected error setting invalid hex")
	}
	if err := f.Set("b64", "not-a-base64!!"); err == nil {
		t.Errorf("expected error setting invalid base64")
	}
	if err := f.Set("mask", "invalid.mask.val"); err == nil {
		t.Errorf("expected error setting invalid mask")
	}
	if err := f.Set("time", "invalid-time"); err == nil {
		t.Errorf("expected error setting invalid time")
	}

	// Package-level helpers
	oldCommandLine := CommandLine
	defer func() { CommandLine = oldCommandLine }()

	CommandLine = NewFlagSet("cmd_special", ContinueOnError)
	CommandLine.SetOutput(io.Discard)

	var (
		pkgCount int
		pkgHex   []byte
		pkgB64   []byte
		pkgMask  net.IPMask
		pkgTime  time.Time
	)

	CountVar(&pkgCount, "p-count", "usage")
	CountVarP(&pkgCount, "p-count-p", "c", "usage")
	_ = Count("pkg-cnt", "usage")
	_ = CountP("pkg-cnt-p", "C", "usage")

	BytesHexVar(&pkgHex, "p-hex", nil, "usage")
	BytesHexVarP(&pkgHex, "p-hex-p", "x", nil, "usage")
	_ = BytesHex("pkg-hex", nil, "usage")
	_ = BytesHexP("pkg-hex-p", "X", nil, "usage")

	BytesBase64Var(&pkgB64, "p-b64", nil, "usage")
	BytesBase64VarP(&pkgB64, "p-b64-p", "b", nil, "usage")
	_ = BytesBase64("pkg-b64", nil, "usage")
	_ = BytesBase64P("pkg-b64-p", "B", nil, "usage")

	IPMaskVar(&pkgMask, "p-mask", defaultMask, "usage")
	IPMaskVarP(&pkgMask, "p-mask-p", "m", defaultMask, "usage")
	_ = IPMask("pkg-mask", defaultMask, "usage")
	_ = IPMaskP("pkg-mask-p", "M", defaultMask, "usage")

	TimeVar(&pkgTime, "p-time", defaultTime, formats, "usage")
	TimeVarP(&pkgTime, "p-time-p", "t", defaultTime, formats, "usage")
	_ = Time("pkg-time", defaultTime, formats, "usage")
	_ = TimeP("pkg-time-p", "T", defaultTime, formats, "usage")

	_ = CommandLine.Parse([]string{"-c", "-c", "--p-hex=abcd"})
	if pkgCount != 2 || hex.EncodeToString(pkgHex) != "abcd" {
		t.Errorf("package level special types mismatch: count=%d, hex=%x", pkgCount, pkgHex)
	}
}
