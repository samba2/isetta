package config

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validLogLevels = []string{"trace", "info", "warn"}

func TestSimpleParse(t *testing.T) {
	var exampleConfig = `
[general]
internet_access_test_url = "https://www.foo.com/"
log_level = "trace"

[network]
wsl_to_windows_subnet = "169.254.254.0/24"
px_proxy_port = 3128

[dns]
private_server = "1.2.3.4"
public_server    = "8.8.8.8"
`

	cfg := FromByteBuffer(bytes.NewBufferString(exampleConfig), validLogLevels)
	assert.Equal(t, "https://www.foo.com/", cfg.General.InternetAccessTestUrl)
	assert.Equal(t, "trace", cfg.General.LogLevel)
	assert.Equal(t, "169.254.254.0/24", cfg.Network.WslToWindowsSubnet)
	assert.Equal(t, 3128, cfg.Network.PxProxyPort)
	assert.Equal(t, "1.2.3.4", cfg.Dns.PrivateServer)
	assert.Equal(t, "8.8.8.8", cfg.Dns.PublicServer)
}

func TestDefaults(t *testing.T) {
	var exampleConfig = `
[dns]
private_server = "1.2.3.4"
`

	cfg := FromByteBuffer(bytes.NewBufferString(exampleConfig), validLogLevels)
	assert.Contains(t, cfg.General.InternetAccessTestUrl, "https://")
	assert.NotEmpty(t, cfg.General.LogLevel)
	assert.NotEmpty(t, cfg.Network.PxProxyPort)
	assert.NotEmpty(t, cfg.Dns.PrivateServer)
	assert.NotEmpty(t, cfg.Dns.PublicServer)
}

func TestSubnetSplitting(t *testing.T) {
	var exampleConfig = `
[network]
wsl_to_windows_subnet = "169.254.254.0/24"

[dns]
private_server = "1.2.3.4"
`

	cfg := FromByteBuffer(bytes.NewBufferString(exampleConfig), validLogLevels)
	assert.Equal(t, "169.254.254.1", cfg.Network.P2p.WindowsIp)
	assert.Equal(t, "169.254.254.2", cfg.Network.P2p.LinuxIp)
	assert.Equal(t, "255.255.255.0", cfg.Network.P2p.SubnetMask)
}

func TestSubnetConfigWithIpStillWorks(t *testing.T) {
	var exampleConfig = `
[network]
wsl_to_windows_subnet = "169.254.254.5/24"

[dns]
private_server = "1.2.3.4"
`

	cfg := FromByteBuffer(bytes.NewBufferString(exampleConfig), validLogLevels)
	assert.Equal(t, "169.254.254.1", cfg.Network.P2p.WindowsIp)
	assert.Equal(t, "169.254.254.2", cfg.Network.P2p.LinuxIp)
	assert.Equal(t, "255.255.255.0", cfg.Network.P2p.SubnetMask)
}

func TestFromConfigFile(t *testing.T) {
	configFileDir := getTestConfigFileDir()

	cfg := FromConfigFile(configFileDir, validLogLevels)
	assert.Equal(t, "https://www.google.com/", cfg.General.InternetAccessTestUrl)
	assert.Equal(t, "trace", cfg.General.LogLevel)
	assert.Equal(t, "169.254.254.0/24", cfg.Network.WslToWindowsSubnet)
	assert.Equal(t, 3128, cfg.Network.PxProxyPort)
	assert.Equal(t, "1.2.3.4", cfg.Dns.PrivateServer)
	assert.Equal(t, "8.8.8.8", cfg.Dns.PublicServer)
}

func TestGetProxy(t *testing.T) {
	var exampleConfig = `
[network]
wsl_to_windows_subnet = "1.1.1.0/24"
px_proxy_port = 3128

[dns]
private_server = "1.2.3.4"
`

	cfg := FromByteBuffer(bytes.NewBufferString(exampleConfig), validLogLevels)
	proxyUrl := GetProxyUrl(cfg)
	assert.Equal(t, "http://1.1.1.1:3128", proxyUrl)
}

// TODO
// https://github.com/spf13/viper/issues/761
// env vars + config + unmarshal seems problematic
func TestFromConfigEnvironmentVariableAndConfigFile(t *testing.T) {
	t.Skip("Skipped as not implemented yet")
	configFileDir := getTestConfigFileDir()
	os.Setenv("ISETTA_DNS_PRIVATE_DNS_SERVER", "42.42.42.42")

	cfg := FromConfigFile(configFileDir, validLogLevels)
	// set via config file
	assert.Equal(t, "info", cfg.General.LogLevel)
	// set via env var
	assert.Equal(t, "42.42.42.42", cfg.Dns.PrivateServer)
}

func getTestConfigFileDir() string {
	_, pathToThisFile, _, _ := runtime.Caller(0)
	thisPackageDir := filepath.Dir(pathToThisFile)
	return filepath.Join(thisPackageDir, "../fixture/")
}
