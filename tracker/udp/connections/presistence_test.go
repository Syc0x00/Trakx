package connections

import (
	"net/netip"
	"testing"
)

func TestMarshallUnmarshall(t *testing.T) {
	connections := NewConnections(1, testTimeNever, testTimeNever)
	ipv4 := netip.MustParseAddrPort("1.1.1.1:1234")
	ipv6 := netip.MustParseAddrPort("[::1]:1234")

	id4 := connections.Create(ipv4)
	id6 := connections.Create(ipv6)

	data, err := connections.Marshall()
	if err != nil {
		t.Fatal("failed to marshall connections", err)
	}

	err = connections.Unmarshall(data)
	if err != nil {
		t.Fatal("failed to unmarshall connections", err)
	}

	if connections.associations[ipv4].ID != id4 {
		t.Errorf("ipv4 connection id = %v; want %v", connections.associations[ipv4].ID, id4)
	}
	if connections.associations[ipv6].ID != id6 {
		t.Errorf("ipv6 connection id = %v; want %v", connections.associations[ipv6].ID, id6)
	}
}
