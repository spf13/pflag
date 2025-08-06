package pflag

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"regexp"
	"time"
)

// AddFlags can add flags to FlagSet by struct
func (f *FlagSet) AddFlags(obj interface{}) error {
	objT := reflect.TypeOf(obj)
	if isStructPtr(objT) {
		return fmt.Errorf("%v must be a struct not a struct pointer!", obj)
	}
	return parseFlagsFromTag("", f, objT)
}

// SetValues can set struct values from FlagSet
func (f *FlagSet) SetValues(obj interface{}) error {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if !isStructPtr(objT) {
		return fmt.Errorf("%v must be a struct pointer", obj)
	}
	objT = objT.Elem()
	objV = objV.Elem()

	return setValuesFromFlagSet("", f, objT, objV)
}

type setFlagFunc func(*FlagSet, string, string, string, string) error
type setValueFunc func(*FlagSet, string, reflect.Value) error

var typeOf = reflect.TypeOf

var typeSetFlagMap = map[reflect.Type]setFlagFunc{
	typeOf(bool(true)):        setBoolFlag,
	typeOf([]bool{}):          setBoolSliceFlag,
	typeOf(string("")):        setStringFlag,
	typeOf([]string{}):        setStringSliceFlag,
	typeOf(int(1)):            setIntFlag,
	typeOf([]int{}):           setIntSliceFlag,
	typeOf(int8(1)):           setInt8Flag,
	typeOf(int16(1)):          setInt16Flag,
	typeOf(int32(1)):          setInt32Flag,
	typeOf([]int32{}):         setInt32SliceFlag,
	typeOf(int64(1)):          setInt64Flag,
	typeOf([]int64{}):         setInt64SliceFlag,
	typeOf(uint(1)):           setUintFlag,
	typeOf([]uint{}):          setUintSliceFlag,
	typeOf(uint8(1)):          setUint8Flag,
	typeOf(uint16(1)):         setUint16Flag,
	typeOf(uint32(1)):         setUint32Flag,
	typeOf(uint64(1)):         setUint64Flag,
	typeOf(float32(1)):        setFloat32Flag,
	typeOf([]float32{}):       setFloat32SliceFlag,
	typeOf(float64(1)):        setFloat64Flag,
	typeOf([]float64{}):       setFloat64SliceFlag,
	typeOf(time.Second):       setDurationFlag,
	typeOf([]time.Duration{}): setDurationSliceFlag,
	typeOf(net.IP{}):          setIPFlag,
	typeOf([]net.IP{}):        setIPSliceFlag,
	typeOf(net.IPMask{}):      setIPMaskFlag,
	typeOf(net.IPNet{}):       setIPNetFlag,
}

var typeSetValueMap = map[reflect.Type]setValueFunc{
	typeOf(bool(true)):        setBoolValue,
	typeOf([]bool{}):          setBoolSliceValue,
	typeOf(string("")):        setStringValue,
	typeOf([]string{}):        setStringSliceValue,
	typeOf(int(1)):            setIntValue,
	typeOf([]int{}):           setIntSliceValue,
	typeOf(int8(1)):           setInt8Value,
	typeOf(int16(1)):          setInt16Value,
	typeOf(int32(1)):          setInt32Value,
	typeOf([]int32{}):         setInt32SliceValue,
	typeOf(int64(1)):          setInt64Value,
	typeOf([]int64{}):         setInt64SliceValue,
	typeOf(uint(1)):           setUintValue,
	typeOf([]uint{}):          setUintSliceValue,
	typeOf(uint8(1)):          setUint8Value,
	typeOf(uint16(1)):         setUint16Value,
	typeOf(uint32(1)):         setUint32Value,
	typeOf(uint64(1)):         setUint64Value,
	typeOf(float32(1)):        setFloat32Value,
	typeOf([]float32{}):       setFloat32SliceValue,
	typeOf(float64(1)):        setFloat64Value,
	typeOf([]float64{}):       setFloat64SliceValue,
	typeOf(time.Second):       setDurationValue,
	typeOf([]time.Duration{}): setDurationSliceValue,
	typeOf(net.IP{}):          setIPValue,
	typeOf([]net.IP{}):        setIPSliceValue,
	typeOf(net.IPMask{}):      setIPMaskValue,
	typeOf(net.IPNet{}):       setIPNetValue,
}

func setValuesFromFlagSet(prefix string, flagSet *FlagSet, objT reflect.Type, objV reflect.Value) error {
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		fieldT := objT.Field(i)

		flag := fieldT.Tag.Get("flag")
		if prefix != "" {
			flag = prefix + "." + flag
		}

		setFunc, ok := typeSetValueMap[fieldT.Type]
		if ok {
			if err := setFunc(flagSet, flag, fieldV); err != nil {
				return err
			}
		} else {
			// do recursion process when filed is struct
			if !fieldT.Anonymous && fieldT.Type.Kind() == reflect.Struct {
				if err := setValuesFromFlagSet(flag, flagSet, fieldT.Type, fieldV); err != nil {
					return err
				}
			}
			// not support type, do nothing
		}

	}
	return nil
}

func parseFlagsFromTag(prefix string, flagSet *FlagSet, objT reflect.Type) error {
	for i := 0; i < objT.NumField(); i++ {
		fieldT := objT.Field(i)
		tag := fieldT.Tag
		flag := tag.Get("flag")
		shorthand := tag.Get("short")
		if prefix != "" {
			flag = prefix + "." + flag
			shorthand = ""
		}
		def := tag.Get("default")
		// def from env
		def = defFromEnv(def)
		desc := tag.Get("desc")

		setFunc, ok := typeSetFlagMap[fieldT.Type]
		if ok {
			if err := setFunc(flagSet, flag, shorthand, def, desc); err != nil {
				return err
			}
		} else {
			// do recursion process when filed is struct
			if !fieldT.Anonymous && fieldT.Type.Kind() == reflect.Struct {
				return parseFlagsFromTag(flag, flagSet, fieldT.Type)
			}
			// not support type, do nothing
		}
	}
	return nil
}

func defFromEnv(def string) string {
	matchedParenthesis, err := regexp.MatchString(`^\$\(\w+\)$`, def)
	if err != nil {
		return def
	}

	matchedBrace, err := regexp.MatchString(`^\$\{\w+\}$`, def)
	if err != nil {
		return def
	}

	if matchedParenthesis || matchedBrace {
		return os.Getenv(def[2 : len(def)-1])
	}

	return def
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
