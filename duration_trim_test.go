package pflag

import (
	"testing"
	"time"
)

func TestDurationTrimSpace(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	var d time.Duration
	f.DurationVar(&d, "timeout", 0, "")
	if err := f.Parse([]string{"--timeout", " 5s "}); err != nil {
		t.Fatal(err)
	}
	if d != 5*time.Second {
		t.Fatalf("got %v", d)
	}
}
