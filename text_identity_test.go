package pflag

import (
	"encoding"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetTextRejectsIdenticalTypeNames(t *testing.T) {
	type textCodec interface {
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	}
	constructors := []func(time.Time) textCodec{
		func(value time.Time) textCodec {
			type stamp struct{ time.Time }
			return &stamp{value}
		},
		func(value time.Time) textCodec {
			type stamp struct{ time.Time }
			return &stamp{value}
		},
	}
	first, second := reflect.TypeOf(constructors[0](time.Time{})), reflect.TypeOf(constructors[1](time.Time{}))
	if first == second || first.String() != second.String() {
		t.Fatalf("test requires distinct types with the same display name: %v, %v", first, second)
	}
	for i, source := range constructors {
		for j, target := range constructors {
			t.Run(fmt.Sprintf("%d-to-%d", i, j), func(t *testing.T) {
				f := NewFlagSet("test", ContinueOnError)
				value := source(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
				f.TextVar(value, "stamp", value, "time stamp")
				output := target(time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC))
				before, _ := output.MarshalText()
				err := f.GetText("stamp", output)
				after, _ := output.MarshalText()
				if i != j {
					if err == nil {
						t.Errorf("source %d, target %d: expected a type mismatch", i, j)
					}
					if string(after) != string(before) {
						t.Errorf("source %d, target %d: rejected output changed from %q to %q", i, j, before, after)
					}
				} else {
					want, _ := value.MarshalText()
					if err != nil || string(after) != string(want) {
						t.Errorf("source %d, target %d: got %q, %v; want %q", i, j, after, err, want)
					}
				}
			})
		}
	}
}

func TestGetTextWithWrappedValue(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	value := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	f.TextVar(&value, "stamp", value, "time stamp")
	flag := f.Lookup("stamp")
	flag.Value = struct{ Value }{flag.Value}
	var output time.Time
	if err := f.GetText("stamp", &output); err != nil || !output.Equal(value) {
		t.Errorf("wrapped value: got %v, %v; want %v", output, err, value)
	}
}
