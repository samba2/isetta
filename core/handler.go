package core

import (
	"errors"

	log "org.samba/isetta/simplelogger"
)

type Handler struct {
	RunningAsRoot    bool
	PrivateDnsServer string
	PublicDnsServer  string
	WindowsChecker   WindowsChecker
	DnsConfigurer    DnsConfigurer
	EnvVarPrinter    EnvVarPrinter
	DirectAccess     NetworkConfigurer
	ViaProxy         NetworkConfigurer
}

func (h *Handler) PrintEnvVars() {
	log.Logger.CurrentLogLevel = log.LevelError
	if h.WindowsChecker.IsPingable(h.PrivateDnsServer) {
		h.EnvVarPrinter.PrintExportCommands()
	} else if h.WindowsChecker.IsPingable(h.PublicDnsServer) {
		h.EnvVarPrinter.PrintUnsetCommands()
	}
}

func (h *Handler) ConfigureNetwork() error {
	if !h.RunningAsRoot {
		return errors.New("to configure the network 'isetta' needs to run as root. Try running via sudo")
	}

	err := h.checkRunningOnWsl()
	if err != nil {
		return err
	}

	h.DnsConfigurer.DisableResolveAutoConfGeneration()

	log.Logger.Info("Detecting network connection")
	if h.WindowsChecker.IsPingable(h.PrivateDnsServer) {
		log.Logger.Debug("Private DNS server is reachable")
		log.Logger.Info("Found internet access via proxy")
		err = h.ViaProxy.Configure()
		if err != nil {
			return err
		}

	} else if h.WindowsChecker.IsPingable(h.PublicDnsServer) {
		log.Logger.Debug("Public DNS server is reachable")
		log.Logger.Info("Found direct internet connection")
		err = h.DirectAccess.Configure()
		if err != nil {
			return err
		}
	} else {
		return errors.New("neither the private nor the public DNS server is reachable - are you offline?")
	}

	return nil
}

func (h *Handler) checkRunningOnWsl() error {
	if h.WindowsChecker.IsRunningOnWsl2() {
		log.Logger.Debug("Running on WSL2")
		return nil
	} else {
		return errors.New("did not detect required WSL version 2")
	}
}
