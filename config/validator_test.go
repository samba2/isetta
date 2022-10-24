package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTooHighPortNumber(t *testing.T) {
	var exampleConfig = `
[network]
# too high port number
px_proxy_port = 70000

[dns]
private_server = "1.2.3.4"
`

	err := validate(exampleConfig)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "PxProxyPort")
}

func TestValidationErrorWhenMandatoryPrivateDnsServerIsMissing(t *testing.T) {
	err := validate("")
	assert.Error(t, err)
}

func TestTranslationIsWorking(t *testing.T) {
	var exampleConfig = `
[general]
internet_access_test_url = "foo"
	
[dns]
private_server = "1.2.3.4"
`
	err := validate(exampleConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternetAccessTestUrl: foo is an an invalid URL")
}

func TestErrorOnTooSmallNetwork(t *testing.T) {
	var exampleConfig = `
[network]
wsl_to_windows_subnet = "169.254.254.0/31"

[dns]
private_server = "1.2.3.4"
`

	err := validate(exampleConfig)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too small")
}

func TestErrorOnWrongLogLevel(t *testing.T) {
	var exampleConfig = `
[general]
log_level = "foolevel"
	
[dns]
private_server = "1.2.3.4"
`

	conf := readConfigFromBuffer(bytes.NewBufferString(exampleConfig))
	myValidator := NewValidator(&conf, []string{"trace"})
	err := myValidator.DoValidate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "foolevel")
	assert.Contains(t, err.Error(), "trace")
}

func validate(config string) error {
	conf := readConfigFromBuffer(bytes.NewBufferString(config))
	myValidator := NewValidator(&conf, []string{"info"})
	err := myValidator.DoValidate()
	return err
}
