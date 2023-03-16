// +build wsl

// tests here need a WSL underneath

package windows

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunInPowershellOk(t *testing.T) {
	result := runInPowerShell("echo HELLOWORLD")
	assert.Equal(t, "HELLOWORLD", result)
}

func TestIsPortOpenOnWindows(t *testing.T) {
	assert.True(t, isPortOpenOnWindows("445"))
	assert.False(t, isPortOpenOnWindows("42"))
}

func TestIsPingableOnWindows(t *testing.T) {
	checker := WindowsCheckerImpl{PxProxyPort: 3128}
	assert.True(t, checker.IsPingable("127.0.0.1"))
	assert.False(t, checker.IsPingable("42.42.42.42"))
}
