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
	mockHttpChecker = mocks.NewHttpChecker(t)

	handler = Handler{
		RunningAsRoot:     true,
		WindowsChecker:    mockWinChecker,
		DnsConfigurer:     mockDnsConfigurer,
		EnvVarPrinter:     mockEnvVarPrinter,
		InternalDnsServer: "42.42.42.42",
		PublicDnsServer:   "8.8.8.8",
		DirectAccess:      mockDirectAccess,
		ViaProxy:          mockViaProxy,
		InternetChecker:   InternetChecker{
			HttpChecker: mockHttpChecker,
			TimeoutInMilliseconds: 100,
		},
	}
}

func setupNoInternetConnection() {
	mockHttpChecker.On("HasDirectInternetAccess", 100).Return(false)
	mockHttpChecker.On("HasInternetAccessViaProxy", 100).Return(false)
}



func TestShortCircuitIfHttpConnectionAlreadyPossible(t *testing.T) {
	setupHandler(t)
	mockHttpChecker.On("HasDirectInternetAccess", 100).Return(true)
	mockHttpChecker.On("HasInternetAccessViaProxy", 100).Return(false)

	assert.NoError(t, handler.ConfigureNetwork())
}


func TestNetworkConfigRequiresRoot(t *testing.T) {
	setupHandler(t)
	setupNoInternetConnection()

	handler.RunningAsRoot = false
	assert.Error(t, handler.ConfigureNetwork())
}

func TestErrorWhenNotOnWsl(t *testing.T) {
	setupHandler(t)
	setupNoInternetConnection()

	mockWinChecker.On("IsRunningOnWsl2").Return(false)
	assert.Error(t, handler.ConfigureNetwork())
}

func TestErrorWhenNoDnsServerIsReached(t *testing.T) {
	setupHandler(t)
	setupNoInternetConnection()

	mockWinChecker.On("IsRunningOnWsl2").Return(true)
	mockDnsConfigurer.On("DisableResolveAutoConfGeneration").Return()
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(false)
	mockWinChecker.On("IsPingable", "8.8.8.8").Return(false)
	assert.Error(t, handler.ConfigureNetwork())
}

func TestPerformsDirectConfigWhenPublicDnsIsReachable(t *testing.T) {
	setupHandler(t)
	setupNoInternetConnection()

	mockWinChecker.On("IsRunningOnWsl2").Return(true)
	mockDnsConfigurer.On("DisableResolveAutoConfGeneration").Return()
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(false)
	mockWinChecker.On("IsPingable", "8.8.8.8").Return(true)
	mockDirectAccess.On("Configure").Return(nil)
	assert.NoError(t, handler.ConfigureNetwork())
}

func TestPerformsConfigViaProxyWhenPublicDnsIsReachable(t *testing.T) {
	setupHandler(t)
	setupNoInternetConnection()

	mockWinChecker.On("IsRunningOnWsl2").Return(true)
	mockDnsConfigurer.On("DisableResolveAutoConfGeneration").Return()
	mockWinChecker.On("IsPingable", "42.42.42.42").Return(true)
	mockViaProxy.On("Configure").Return(nil)
	assert.NoError(t, handler.ConfigureNetwork())
}

func TestWhenInternalDnsIsReachableExportStatementsArePrinted(t *testing.T) {
	setupHandler(t)

	mockWinChecker.On("IsPingable", "42.42.42.42").Return(true)
	mockEnvVarPrinter.On("PrintExportCommands")
	handler.PrintEnvVars()
}

func TestWhenPublicDnsIsReachableUnsetStatementsArePrinted(t *testing.T) {
	setupHandler(t)
	
	mockWinChecker.
		// internal DNS not reachable
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
