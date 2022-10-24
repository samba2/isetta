package core

import (
	"errors"
	log "org.samba/isetta/simplelogger"
)

type DirectAccess struct {
	PublicDnsServer string
	DnsConfigurer   DnsConfigurer
	LinuxPinger     LinuxPinger
	LinuxConfigurer LinuxConfigurer
	HttpChecker     HttpChecker
	EnvVarPrinter   EnvVarPrinter
}

func (d *DirectAccess) Configure() error {
	d.DnsConfigurer.ActivateDnsServer(d.PublicDnsServer)
	d.EnvVarPrinter.WarnIfProxyVarSet()
	err := d.configureDefaultGatewayIfNeeded()
	if err != nil {
		return err
	}

	err = d.checkDirectAccess()
	if err != nil {
		return err
	}

	return nil
}

func (d *DirectAccess) configureDefaultGatewayIfNeeded() error {
	if !d.isPublicDnsServerUp() {
		log.Logger.Debug("Configuring default gateway")
		d.LinuxConfigurer.DeleteDefaultGateway()
		d.LinuxConfigurer.AddDefaultGateway()

		if !d.isPublicDnsServerUp() {
			return errors.New("failed to adjust default gateway 🤔")
		}
	}

	return nil
}

func (d *DirectAccess) isPublicDnsServerUp() bool {
	if d.LinuxPinger.Ping(d.PublicDnsServer) {
		log.Logger.Trace("Public DNS server %v can be reached from within Linux. Default gateway works", d.PublicDnsServer)
		return true
	} else {
		log.Logger.Trace("Public DNS server %v can't be reached from within Linux", d.PublicDnsServer)
		return false
	}
}

func (d *DirectAccess) checkDirectAccess() error {
	if d.HttpChecker.HasDirectInternetAccess() {
		log.Logger.Info("Done setting up WSL network for direct internet access")
		return nil
	} else {
		return errors.New("failed setting up WSL network for direct internet access")
	}
}
