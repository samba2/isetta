package envvars

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildPrintExportCommands(t *testing.T) {

	uut := ConsoleEnvVarPrinter{
		WindowsIp: "1.1.1.1",
		PxProxyPort: 4242,
	}

	assert.Regexp(t, "^export HTTPS_PROXY=http://1.1.1.1:4242", uut.buildPrintExportCommands())	
}

func TestHttpEnvVarsAreNotSet(t *testing.T) {
	uut := ConsoleEnvVarPrinter{}
	assert.False(t, uut.areHttpEnvVarsSet())	
}

func TestHttpEnvVarsAreSet(t *testing.T) {
	uut := ConsoleEnvVarPrinter{}
	os.Setenv("http_proxy", "foo")

	assert.True(t, uut.areHttpEnvVarsSet())	
}
