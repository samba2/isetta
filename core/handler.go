package core

import (
	"errors"

	log "org.samba/isetta/simplelogger"
)

type Handler struct {
	RunningAsRoot     bool
	InternalDnsServer string
	PublicDnsServer   string
	WindowsChecker    WindowsChecker
	DnsConfigurer     DnsConfigurer
	EnvVarPrinter     EnvVarPrinter
	DirectAccess      NetworkConfigurer
	ViaProxy          NetworkConfigurer
	InternetChecker   InternetChecker
}

func (h *Handler) PrintEnvVars() {
	log.Logger.CurrentLogLevel = log.LevelError
	if h.WindowsChecker.IsPingable(h.InternalDnsServer) {
		h.EnvVarPrinter.PrintExportCommands()
	} else if h.WindowsChecker.IsPingable(h.PublicDnsServer) {
		h.EnvVarPrinter.PrintUnsetCommands()
	}
}

func (h *Handler) ConfigureNetwork() error {
	log.Logger.Info("Checking if internet can already by reached via HTTP")
	if h.InternetChecker.HasInternetAccess() {
		log.Logger.Info("Internet is already accessible. No further setup needed")
		return nil
	}

	if !h.RunningAsRoot {
		return errors.New("to configure the network 'isetta' needs to run as root. Try running via sudo")
	}
	
	err := h.checkRunningOnWsl()
	if err != nil {
		return err
	}
	
	h.DnsConfigurer.DisableResolveAutoConfGeneration()

	log.Logger.Info("Detecting network connection")
	if h.WindowsChecker.IsPingable(h.InternalDnsServer) {
		log.Logger.Debug("Internal DNS server is reachable")
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
		return errors.New("neither the internal nor the public DNS server is reachable - are you offline?")
	}

	return nil
}

func (h *Handler) checkRunningOnWsl() error {
	if h.WindowsChecker.IsRunningOnWsl2() {
		log.Logger.Debug("Running on WSL2")
		return nil
	} else {
		return errors.New("isetta requires WSL2 but this Linux environment is running in something else. Run 'wsl.exe --list --verbose' for details")
	}
}
