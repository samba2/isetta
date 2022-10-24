package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCidr(t *testing.T) {
	ipWithCidr := getCidrNotation("192.168.1.1", "255.255.255.0")
	assert.Equal(t, "192.168.1.1/24", ipWithCidr)
}

func TestGetBroadcastAddress(t *testing.T) {
	broadcast := getBroadcast("192.168.1.1", "255.255.255.0")
	assert.Equal(t, "192.168.1.255", broadcast)
}
