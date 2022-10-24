package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"org.samba/isetta/mocks"
)

var mockDirectAccess *mocks.NetworkConfigurer
var mockViaProxy *mocks.NetworkConfigurer

var handler Handler

func setupHandler(t *testing.T) {
	mockWinChecker = mocks.NewWindowsChecker(t)
	mockDnsConfigurer = mocks.NewDnsConfigurer(t)
	mockEnvVarPrinter = mocks.NewEnvVarPrinter(t)
	mockDirectAccess = mocks.NewNetworkConfigurer(t)
	mockViaProxy = mocks.NewNetworkConfigurer(t)

	handler = Handler {
		RunningAsRoot: true,
		WindowsChecker: mockWinChecker,
		DnsConfigurer: mockDnsConfigurer,
		EnvVarPrinter: mockEnvVarPrinter,
		PrivateDnsServer: "42.42.42.42",
		PublicDnsServer: "8.8.8.8",
		DirectAccess: mockDirectAccess,
		ViaProxy: mockViaProxy,
	}
}

func TestErrorWhenNotOnWsl(t *testing.T) {
	setupHandler(t)

	mockWinChecker.On("IsRunningOnWsl2").Return(false)
	assert.Error(t, handler.ConfigureNetwork())
}

func TestErrorWhenNoDnsServerIsReached(t *testing.T) {
	setupHandler(t)

	mockWinChecker.On("IsRunningOnWsl2").Return(true)
	mockDnsConfigurer.On("DisableResolveAutoConfGeneration").Return()
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(false)
	mockWinChecker.On("IsPingable", "8.8.8.8").Return(false)
	assert.Error(t, handler.ConfigureNetwork())
}

func TestPerformsDirectConfigWhenPublicDnsIsReachable(t *testing.T) {
	setupHandler(t)

	mockWinChecker.On("IsRunningOnWsl2").Return(true)
	mockDnsConfigurer.On("DisableResolveAutoConfGeneration").Return()
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(false)
	mockWinChecker.On("IsPingable", "8.8.8.8").Return(true)
	mockDirectAccess.On("Configure").Return(nil)
	assert.NoError(t, handler.ConfigureNetwork())
}

func TestPerformsConfigViaProxyWhenPublicDnsIsReachable(t *testing.T) {
	setupHandler(t)

	mockWinChecker.On("IsRunningOnWsl2").Return(true)
	mockDnsConfigurer.On("DisableResolveAutoConfGeneration").Return()
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(true)
	mockViaProxy.On("Configure").Return(nil)
	assert.NoError(t, handler.ConfigureNetwork())
}

func TestWhenPrivateDnsIsReachableExportStatementsArePrinted(t *testing.T) {	
	setupHandler(t)
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(true)
	mockEnvVarPrinter.On("PrintExportCommands")
	handler.PrintEnvVars()
}

func TestWhenPublicDnsIsReachableUnsetStatementsArePrinted(t *testing.T) {	
	setupHandler(t)
	mockWinChecker.
		// private DNS not reachable
		On("IsPingable", "42.42.42.42").Return(false).
		// public DNS reachable
		On("IsPingable", "8.8.8.8").Return(true)

	mockEnvVarPrinter.On("PrintUnsetCommands")
	handler.PrintEnvVars()
}

func TestEnvVarsArePrintedIfNonRoot(t *testing.T) {
	setupHandler(t)
	handler.RunningAsRoot = false
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(true)
	mockEnvVarPrinter.On("PrintExportCommands")
	handler.PrintEnvVars()
}

func TestNetworkConfigRequiresRoot(t *testing.T) {
	setupHandler(t)
	handler.RunningAsRoot = false
	assert.Error(t, handler.ConfigureNetwork())
}