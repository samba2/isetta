package core

type DnsConfigurer interface {
	// check if given IP is the active DNS server in /etc/resolve.conf and update if needed
	ActivateDnsServer(dnsServerIp string)

	// Ensure that in /etc/wsl.conf 'generateResolvConf' is set to 'false'
	// Creates /etc/wsl.conf if not exists
	DisableResolveAutoConfGeneration()
}

type EnvVarPrinter interface {
	PrintExportCommands()
	PrintUnsetCommands()
	WarnIfProxyVarSet()
}

type LinuxPinger interface {
	Ping(host string) bool
}

type LinuxConfigurer interface {
	SetP2pInterface()
	DeleteDefaultGateway()
	AddDefaultGateway()
}

type WindowsChecker interface {
	IsPingable(host string) bool
	IsPxProxyRunning() bool
	IsRunningOnWsl2() bool
}

type WindowsConfigurer interface {
	Init()    // deferred construction and object setup,
	Cleanup() // cleanup temporary resources
	AddP2pAddress(successChecker func() bool) error
	SetPortProxy(successChecker func() bool) error
}

type HttpChecker interface {
	HasDirectInternetAccess() bool
	HasInternetAccessViaProxy() bool
	IsPxProxyReachable() bool
}

type NetworkConfigurer interface {
	Configure() error
}