package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"org.samba/isetta/mocks"
)

var direct DirectAccess

func setupDirect(t *testing.T) {
	mockDnsConfigurer = mocks.NewDnsConfigurer(t)
	mockLinuxPinger = mocks.NewLinuxPinger(t)
	mockLinuxConfigurer = mocks.NewLinuxConfigurer(t)
	mockHttpChecker = mocks.NewHttpChecker(t)
	mockEnvVarPrinter = mocks.NewEnvVarPrinter(t)

	direct = DirectAccess{
		PublicDnsServer: "8.8.8.8",
		DnsConfigurer:   mockDnsConfigurer,
		LinuxPinger: mockLinuxPinger,
		LinuxConfigurer: mockLinuxConfigurer,
		HttpChecker:     mockHttpChecker,
		EnvVarPrinter:   mockEnvVarPrinter,
	}
}

func TestConfigureDirectInternetAccess(t *testing.T) {
	setupDirect(t)
	mockDnsConfigurer.On("ActivateDnsServer", "8.8.8.8").Return()
	mockEnvVarPrinter.On("WarnIfProxyVarSet").Return()
	// default gateway is ok as we have access to public DNS server
	mockLinuxPinger.On("Ping", "8.8.8.8").Return(true)

	// test url can be reached directly
	mockHttpChecker.On("HasDirectInternetAccess").Return(true)

	assert.NoError(t, direct.Configure())
}

func TestConfigureDirectInternetAccessWithDefaultGatewaySetup(t *testing.T) {
	setupDirect(t)
	mockDnsConfigurer.On("ActivateDnsServer", "8.8.8.8").Return()
	mockEnvVarPrinter.On("WarnIfProxyVarSet").Return()
	// default gateway is ok as we have access to public DNS server
	mockLinuxPinger.On("Ping", "8.8.8.8").Return(false).Once()
	mockLinuxConfigurer.On("DeleteDefaultGateway")
	mockLinuxConfigurer.On("AddDefaultGateway")
	// default GW setup was ok
	mockLinuxPinger.On("Ping", "8.8.8.8").Return(true).Once()

	// test url can be reached directly
	mockHttpChecker.On("HasDirectInternetAccess").Return(true)

	assert.NoError(t, direct.Configure())
}

func TestCheckHasDirectAccess(t *testing.T) {
	setupDirect(t)
	mockHttpChecker.On("HasDirectInternetAccess").Return(true)
	assert.NoError(t, direct.checkDirectAccess())
}

func TestCheckHasNoDirectAccess(t *testing.T) {
	setupDirect(t)
	mockHttpChecker.On("HasDirectInternetAccess").Return(false)
	assert.Error(t, direct.checkDirectAccess())
}

