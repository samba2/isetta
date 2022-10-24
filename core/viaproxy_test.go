package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"org.samba/isetta/mocks"
)

var viaProxy ViaProxy

func setupViaProxy(t *testing.T) {
	mockWinChecker = mocks.NewWindowsChecker(t)
	mockWinConfigurer = mocks.NewWindowsConfigurer(t)
	mockHttpChecker = mocks.NewHttpChecker(t)
	mockDnsConfigurer = mocks.NewDnsConfigurer(t)
	mockEnvVarPrinter = mocks.NewEnvVarPrinter(t)
	mockLinuxPinger = mocks.NewLinuxPinger(t)
	mockLinuxConfigurer = mocks.NewLinuxConfigurer(t)

	viaProxy = ViaProxy{
		LinuxP2pIp: "linux-ip",
		WindowsP2pIp: "windows-ip",
		PxProxyPort: 3128,
		PrivateDnsServer: "42.42.42.42",
		WindowsChecker: mockWinChecker,
		WindowsConfigurer: mockWinConfigurer,
		DnsConfigurer: mockDnsConfigurer,
		LinuxPinger: mockLinuxPinger,
		LinuxConfigurer: mockLinuxConfigurer,
		HttpChecker: mockHttpChecker,
	}
}

func TestConfigureAccessViaProxy(t *testing.T) {
	setupViaProxy(t)

	// PX proxy is active on Windows
	mockWinChecker.On("IsPxProxyRunning").Return(true)

	// set private DNS server in resolve.conf
	mockDnsConfigurer.On("ActivateDnsServer", "42.42.42.42").Return()

	// Linux P2P IP is not up, but after SetP2pInterface was called, it will be up
	mockLinuxPinger.On("Ping", "linux-ip").Return(false).Once()
	mockLinuxConfigurer.On("SetP2pInterface").Return()
	mockLinuxPinger.On("Ping", "linux-ip").Return(true).Once()

	// Windows IP can't be reached...
	mockLinuxPinger.On("Ping", "windows-ip").Return(false)

	// ...need to configure Windows side
	mockWinConfigurer.On("Init").Return()
	mockWinConfigurer.On("Cleanup").Return()
	mockWinConfigurer.On("AddP2pAddress", mock.Anything).Return(nil)
	mockWinConfigurer.On("SetPortProxy", mock.Anything).Return(nil)

	// private DNS can't be reached, first need to configure
	// default gateway. After that, it can be reached
	mockLinuxPinger.On("Ping", "42.42.42.42").Return(false).Once()
	// default gateway on Linux side
	mockLinuxConfigurer.On("DeleteDefaultGateway").Return()
	mockLinuxConfigurer.On("AddDefaultGateway").Return()
	mockLinuxPinger.On("Ping", "42.42.42.42").Return(true).Once()

	// cool, setup worked
	mockHttpChecker.On("HasInternetAccessViaProxy").Return(true)
	assert.NoError(t, viaProxy.Configure())
}

func TestCheckHasAccessViaProxy(t *testing.T) {
	setupViaProxy(t)
	mockHttpChecker.On("HasInternetAccessViaProxy").Return(true)
	assert.NoError(t, viaProxy.checkAccessViaProxy())
}

func TestCheckHasNoAccessViaProxy(t *testing.T) {
	setupViaProxy(t)
	mockHttpChecker.On("HasInternetAccessViaProxy").Return(false)
	assert.Error(t, viaProxy.checkAccessViaProxy())
}

func TestCheckPxProxyIsRunning(t *testing.T) {
	setupViaProxy(t)
	mockWinChecker.On("IsPxProxyRunning").Return(true)
	assert.NoError(t, viaProxy.checkPxProxyRunning())
}

func TestCheckPxProxyIsNotRunning(t *testing.T) {
	setupViaProxy(t)
	mockWinChecker.On("IsPxProxyRunning").Return(false)
	assert.Error(t, viaProxy.checkPxProxyRunning())
}

func TestWindowsSideOk1(t *testing.T) {
	setupViaProxy(t)
	mockLinuxPinger.On("Ping", "windows-ip").Return(true)
	mockHttpChecker.On("IsPxProxyReachable").Return(true)
	assert.Equal(t, true, viaProxy.isWindowsSideOk())
}

func TestWindowsSideOk2(t *testing.T) {
	setupViaProxy(t)
	mockLinuxPinger.On("Ping", "windows-ip").Return(false)
	assert.Equal(t, false, viaProxy.isWindowsSideOk())
}

func TestWindowsSideOk3(t *testing.T) {
	setupViaProxy(t)
	mockLinuxPinger.On("Ping", "windows-ip").Return(true)
	mockHttpChecker.On("IsPxProxyReachable").Return(false)
	assert.Equal(t, false, viaProxy.isWindowsSideOk())
}

func TestSetupWslPspInterfaceIsNotNeeded(t *testing.T) {
	setupViaProxy(t)
	mockLinuxPinger.On("Ping", "linux-ip").Return(true)
	viaProxy.setupLinuxP2pInterfaceIfNeeded()
	mockLinuxConfigurer.AssertNotCalled(t, "SetP2pInterface")
}

func TestSetupWslPspInterfaceIfNeeded(t *testing.T) {
	setupViaProxy(t)

	mockLinuxPinger.On("Ping", "linux-ip").Return(false).Once()
	mockLinuxConfigurer.On("SetP2pInterface").Return()
	// SetP2pInterface fixed it, now ping is successful
	mockLinuxPinger.On("Ping", "linux-ip").Return(true).Once()
	viaProxy.setupLinuxP2pInterfaceIfNeeded()
}

func TestSuccessfullyConfigureAccessViaProxy(t *testing.T) {
	setupViaProxy(t)
	mockWinConfigurer.On("Init").Return()
	mockWinConfigurer.On("Cleanup").Return()
	mockWinConfigurer.On("AddP2pAddress", mock.Anything).Return(nil)
	mockWinConfigurer.On("SetPortProxy", mock.Anything).Return(nil)
	assert.NoError(t, viaProxy.configureWindowsSide())
}

func TestConfigureAccessViaProxyHasError1(t *testing.T) {
	setupViaProxy(t)
	mockWinConfigurer.On("Init").Return()
	mockWinConfigurer.On("Cleanup").Return()
	mockWinConfigurer.On("AddP2pAddress", mock.Anything).Return(errors.New(""))

	assert.Error(t, viaProxy.configureWindowsSide())
	mockWinConfigurer.AssertNotCalled(t, "AddP2pAddress")
}

func TestConfigureAccessViaProxyHasError2(t *testing.T) {
	setupViaProxy(t)
	mockWinConfigurer.On("Init").Return()
	mockWinConfigurer.On("Cleanup").Return()
	mockWinConfigurer.On("AddP2pAddress", mock.Anything).Return(nil)
	mockWinConfigurer.On("SetPortProxy", mock.Anything).Return(errors.New(""))

	assert.Error(t, viaProxy.configureWindowsSide())
}


func TestErrorWhenSettingP2pAddressFailed(t *testing.T) {
	setupViaProxy(t)
	mockLinuxPinger.On("Ping", "linux-ip").Return(false)
	mockLinuxConfigurer.On("SetP2pInterface").Return()
	assert.Error(t, viaProxy.setupLinuxP2pInterfaceIfNeeded())
}

func TestErrorWhenDefaultGatewayConfigFailed(t *testing.T) {
	setupViaProxy(t)
	mockLinuxPinger.On("Ping", "42.42.42.42").Return(false)
	mockLinuxConfigurer.On("DeleteDefaultGateway").Return()
	mockLinuxConfigurer.On("AddDefaultGateway").Return()
	assert.Error(t, viaProxy.configureDefaultGatewayIfNeeded())
}