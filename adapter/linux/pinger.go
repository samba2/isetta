package linux

import (
	"time"

	"github.com/go-ping/ping"
	"org.samba/isetta/helper"
)

type LinuxPingerImpl struct{}

func (LinuxPingerImpl) Ping(host string) bool {
	pinger, err := ping.NewPinger(host)
	helper.AssertNoError2(err)

	// required to use ICMP
	pinger.SetPrivileged(true)
	pinger.Count = 2
	pinger.Interval = 300 * time.Millisecond
	pinger.Timeout = 2 * time.Second

	err = pinger.Run()
	return err == nil && pinger.Statistics().PacketsRecv != 0
}
