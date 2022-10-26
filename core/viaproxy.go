package core

import (
	"errors"
	"fmt"

	log "org.samba/isetta/simplelogger"
)

type ViaProxy struct {
	LinuxP2pIp        string
	WindowsP2pIp      string
	PxProxyPort       int
	InternalDnsServer string
	WindowsChecker    WindowsChecker
	WindowsConfigurer WindowsConfigurer
	DnsConfigurer     DnsConfigurer
	LinuxPinger       LinuxPinger
	LinuxConfigurer   LinuxConfigurer
	HttpChecker       HttpChecker
}

func (p *ViaProxy) Configure() error {
	err := p.checkPxProxyRunning()
	if err != nil {
		return err
	}

	p.DnsConfigurer.ActivateDnsServer(p.InternalDnsServer)
	err = p.setupLinuxP2pInterfaceIfNeeded()
	if err != nil {
		return err
	}

	if !p.isWindowsSideOk() {
		err := p.configureWindowsSide()
		if err != nil {
			return err
		}
	}

	err = p.configureDefaultGatewayIfNeeded()
	if err != nil {
		return err
	}

	err = p.checkAccessViaProxy()
	if err != nil {
		return err
	}

	return nil
}

// directly check on Windows if PX proxy is running at all
func (p *ViaProxy) checkPxProxyRunning() error {
	if p.WindowsChecker.IsPxProxyRunning() {
		log.Logger.Debug("PX proxy is running on Windows port %v", p.PxProxyPort)
		return nil
	} else {
		msg := fmt.Sprintf("Error: PX proxy is not running on Windows port %v", p.PxProxyPort)
		return errors.New(msg)
	}
}

func (p *ViaProxy) setupLinuxP2pInterfaceIfNeeded() error {
	if !p.isLinuxP2pIpUp() {
		log.Logger.Debug("Adding address %v to Linux", p.LinuxP2pIp)
		p.LinuxConfigurer.SetP2pInterface()

		// post condition
		if !p.isLinuxP2pIpUp() {
			return fmt.Errorf("failed to add P2P address %v to Linux", p.LinuxP2pIp)
		}
	}

	return nil
}

func (p *ViaProxy) isLinuxP2pIpUp() bool {
	if p.LinuxPinger.Ping(p.LinuxP2pIp) {
		log.Logger.Debug("Linux P2P address %v is up", p.LinuxP2pIp)
		return true
	} else {
		log.Logger.Debug("Linux P2P address %v is not up", p.LinuxP2pIp)
		return false
	}
}

func (p *ViaProxy) isWindowsSideOk() bool {
	return p.isWindowsP2pIpUp() && p.IsPxProxyReachable()
}

func (p *ViaProxy) isWindowsP2pIpUp() bool {
	if p.LinuxPinger.Ping(p.WindowsP2pIp) {
		log.Logger.Debug("Windows P2P address %v is up", p.WindowsP2pIp)
		return true
	} else {
		log.Logger.Debug("Windows P2P address %v is not up", p.WindowsP2pIp)
		return false
	}
}

func (p *ViaProxy) IsPxProxyReachable() bool {
	if p.HttpChecker.IsPxProxyReachable() {
		log.Logger.Debug("Px Proxy is reachable from within Linux")
		return true
	} else {
		log.Logger.Debug("Px Proxy not reachable from Linux side")
		return false
	}
}

func (p *ViaProxy) configureWindowsSide() error {
	p.WindowsConfigurer.Init()
	defer p.WindowsConfigurer.Cleanup()

	log.Logger.Debug("Adding Windows P2p address %v", p.WindowsP2pIp)
	windowsIpReachableFromWslChecker := func() bool { return p.LinuxPinger.Ping(p.WindowsP2pIp) }
	err := p.WindowsConfigurer.AddP2pAddress(windowsIpReachableFromWslChecker)
	if err != nil {
		return err
	}

	err = p.WindowsConfigurer.SetPortProxy(p.HttpChecker.HasInternetAccessViaProxy)
	if err != nil {
		return err
	}

	return nil
}

func (p *ViaProxy) configureDefaultGatewayIfNeeded() error {
	if !p.isInternalDnsServerUp() {
		log.Logger.Debug("Configuring default gateway")
		p.LinuxConfigurer.DeleteDefaultGateway()
		p.LinuxConfigurer.AddDefaultGateway()

		if !p.isInternalDnsServerUp() {
			return errors.New("failed to adjust default gateway 🤔")
		}
	}

	return nil
}

func (p *ViaProxy) isInternalDnsServerUp() bool {
	if p.LinuxPinger.Ping(p.InternalDnsServer) {
		log.Logger.Debug("Internal DNS %v can be reached from within Linux. Default gateway works", p.InternalDnsServer)
		return true
	} else {
		log.Logger.Debug("Internal DNS %v can't be reached from within Linux", p.InternalDnsServer)
		return false
	}
}

func (p *ViaProxy) checkAccessViaProxy() error {
	if p.HttpChecker.HasInternetAccessViaProxy() {
		log.Logger.Info("Done setting up Linux network via proxy")
		return nil
	} else {
		return errors.New("failed setting up Linux network via proxy")
	}
}
