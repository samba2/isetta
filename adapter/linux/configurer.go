package linux

import (
	"net"
	"os/exec"
	"strings"

	"github.com/3th1nk/cidr"
	"org.samba/isetta/helper"
	log "org.samba/isetta/simplelogger"
)

type LinuxConfigurerImpl struct {
	WindowsIp        string
	LinuxIp          string
	SubnetMask       string
}

func (l *LinuxConfigurerImpl) SetP2pInterface() {
	linuxIpCidr := getCidrNotation(l.LinuxIp, l.SubnetMask)
	broadcast := getBroadcast(l.LinuxIp, l.SubnetMask)
	err := exec.Command("ip", "addr", "change", linuxIpCidr, "broadcast", broadcast, "dev", "eth0", "label", "eth0:1").Run()
	helper.AssertNoError2(err)
}

// returns IP address in CIDR notation like 192.168.2.1/24
func getCidrNotation(ip string, subnetMask string) string {
	ip2 := net.ParseIP(ip)
	assertNotNil(ip2)
	subnetMask2 := net.ParseIP(subnetMask)
	assertNotNil(subnetMask2)

	ipNet := net.IPNet{IP: ip2, Mask: net.IPMask(subnetMask2.To4())}
	return ipNet.String()
}

func getBroadcast(ip string, subnetMask string) string {
	cidr2, err := cidr.ParseCIDR(getCidrNotation(ip, subnetMask))
	helper.AssertNoError2(err)
	return cidr2.Broadcast()
}


func assertNotNil(ip net.IP) {
	if ip == nil {
		log.Logger.Error("Error parsing IP address %v", ip)
	}
}

func (l *LinuxConfigurerImpl) DeleteDefaultGateway() {
	cmd := exec.Command("ip", "route", "delete", "default")
	out, err := cmd.CombinedOutput()
	if err == nil {
		log.Logger.Trace("Deleted existing default route")
	} else {
		tmp := strings.TrimSpace(string(out))		
		log.Logger.Trace("Failed to deleted default route. Maybe it wasn't set? Output was: %v", tmp)
	}
}

func (l *LinuxConfigurerImpl) AddDefaultGateway() {
	cmd := exec.Command("ip", "route", "add", "default", "via", l.WindowsIp)
	out, err := cmd.CombinedOutput()
	helper.AssertNoError(err, "error configuring default gateway on Linux side: %v", string(out))
}
