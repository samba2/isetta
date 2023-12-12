package main

import (
	"flag"
	"fmt"
	"os"

	"org.samba/isetta/adapter/dnsconfig"
	"org.samba/isetta/adapter/envvars"
	"org.samba/isetta/adapter/httpchecker"
	"org.samba/isetta/adapter/linux"
	"org.samba/isetta/adapter/windows"
	"org.samba/isetta/config"
	"org.samba/isetta/core"
	"org.samba/isetta/gsudo"
	"org.samba/isetta/helper"
	log "org.samba/isetta/simplelogger"
)

const version = "0.5.1"

func main() {
	envSettings := flag.Bool("env-settings", false, "Prints environment config. Handy if called via 'source'")
	printVersion := flag.Bool("version", false, "Print isetta version")
	flag.Parse()

	if *printVersion {
		fmt.Printf("Isetta version %v\n", version)
		os.Exit(0)
	}

	conf := config.FromConfigFile("$HOME", log.GetValidLogLevels())
	log.Logger.CurrentLogLevel = log.Levels[conf.General.LogLevel]
	handler := setupDependencies(conf)

	if *envSettings {
		handler.PrintEnvVars()
	} else {
		err := handler.ConfigureNetwork()
		helper.AssertNoError2(err)
	}
}

func setupDependencies(conf config.Config) core.Handler {
	envVarprinter := envvars.ConsoleEnvVarPrinter{
		WindowsIp:   conf.Network.P2p.WindowsIp,
		PxProxyPort: conf.Network.PxProxyPort,
		NoProxy:     conf.Network.NoProxy,
	}

	windowsChecker := windows.WindowsCheckerImpl{PxProxyPort: conf.Network.PxProxyPort}

	windowsConfigurer := windows.WindowsConfigurerImpl{
		WindowsIp:   conf.Network.P2p.WindowsIp,
		SubnetMask:  conf.Network.P2p.SubnetMask,
		PxProxyPort: conf.Network.PxProxyPort,
		Gsudo:       &gsudo.Gsudo{},
	}

	linuxPinger := linux.LinuxPingerImpl{}

	linuxConfigurer := linux.LinuxConfigurerImpl{
		WindowsIp:  conf.Network.P2p.WindowsIp,
		LinuxIp:    conf.Network.P2p.LinuxIp,
		SubnetMask: conf.Network.P2p.SubnetMask,
	}

	dnsConfigurer := dnsconfig.DnsConfigurerImpl{}
	httpchecker, err := httpchecker.New(conf.General.InternetAccessTestUrl, config.GetProxyUrl(conf))
	helper.AssertNoError2(err)

	directAccess := core.DirectAccess{
		PublicDnsServer: conf.Dns.PublicServer,
		DnsConfigurer:   dnsConfigurer,
		LinuxPinger:     &linuxPinger,
		LinuxConfigurer: &linuxConfigurer,
		HttpChecker:     &httpchecker,
		EnvVarPrinter:   &envVarprinter,
	}

	viaproxy := core.ViaProxy{
		// static
		LinuxP2pIp:        conf.Network.P2p.LinuxIp,
		WindowsP2pIp:      conf.Network.P2p.WindowsIp,
		PxProxyPort:       conf.Network.PxProxyPort,
		InternalDnsServer: conf.Dns.InternalServer,
		// objects
		WindowsChecker:    &windowsChecker,
		WindowsConfigurer: &windowsConfigurer,
		DnsConfigurer:     &dnsConfigurer,
		LinuxPinger:       &linuxPinger,
		LinuxConfigurer:   &linuxConfigurer,
		HttpChecker:       &httpchecker,
	}

	handler := core.Handler{
		RunningAsRoot:     os.Geteuid() == 0,
		InternalDnsServer: conf.Dns.InternalServer,
		PublicDnsServer:   conf.Dns.PublicServer,
		WindowsChecker:    &windowsChecker,
		DnsConfigurer:     &dnsConfigurer,
		EnvVarPrinter:     &envVarprinter,
		DirectAccess:      &directAccess,
		ViaProxy:          &viaproxy,
		InternetChecker:   core.NewInternetChecker(&httpchecker),
	}

	return handler
}
