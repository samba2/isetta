package linux

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingOk(t *testing.T) {
	if ! isRoot() {
		t.Skip("Test requires root rights")
	}

	pinger := LinuxPingerImpl{}
	assert.True(t, pinger.Ping("127.0.0.1"))
}

func TestPingNotOk(t *testing.T) {
	if ! isRoot() {
		t.Skip("Test requires root rights")
	}

	pinger := LinuxPingerImpl{}
	assert.False(t, pinger.Ping("42.42.42.42"))
}

func TestPingNotOkSinceNetworkIsNotReachable(t *testing.T) {
	if ! isRoot() {
		t.Skip("Test requires root rights")
	}
	pinger := LinuxPingerImpl{}
	assert.False(t, pinger.Ping("42.43.44.45"))
}


func isRoot() bool {
	return os.Getuid() == 0
}