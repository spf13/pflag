package pflag

import (
	"net"
	"testing"
)

func TestGetIPNilDefault(t *testing.T) {
	t.Run("IP", func(t *testing.T) {
		f := NewFlagSet("test", ContinueOnError)
		f.IP("ip", nil, "IP address")

		ip, err := f.GetIP("ip")
		if err != nil {
			t.Fatalf("GetIP returned error: %v", err)
		}
		if ip != nil {
			t.Fatalf("expected nil IP, got %v", ip)
		}
	})

	t.Run("IPVar", func(t *testing.T) {
		f := NewFlagSet("test", ContinueOnError)
		var ip net.IP
		f.IPVar(&ip, "ip", nil, "IP address")

		got, err := f.GetIP("ip")
		if err != nil {
			t.Fatalf("GetIP returned error: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil IP, got %v", got)
		}
	})
}
