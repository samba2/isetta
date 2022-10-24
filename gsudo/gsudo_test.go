// +build wsl

package gsudo

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	winTempDir   string
	gsudoWslPath string
)

func init() {
	winTempDir = getWindowsTempDir()
	winTempDirInWsl := windowsPathToWsl(winTempDir)
	gsudoWslPath = path.Join(winTempDirInWsl, "gsudo-isetta.exe")
	os.Remove(gsudoWslPath)
}

func TestBinaryIsCopiedAndRemoved(t *testing.T) {
	assert.False(t, fileExists(gsudoWslPath), "binary should initially not exist")

	gsudo := Gsudo{}
	gsudo.Init()
	assert.True(t, fileExists(gsudoWslPath), "binary should exist")

	gsudo.Cleanup()
	assert.False(t, fileExists(gsudoWslPath), "binary should be removed")
}

func TestGsudoRunsViaCmdCall(t *testing.T) {
	gsudo := Gsudo{}
	gsudo.Init()
	defer gsudo.Cleanup()

	out := gsudo.run("status")
	assert.Contains(t, out, "Total active cache sessions")
}

func TestGsudoMultiArgumentWorks(t *testing.T) {
	gsudo := Gsudo{}
	gsudo.Init()
	defer gsudo.Cleanup()

	out := gsudo.run("config CopyEnvironmentVariables False")
	assert.Contains(t, out, "CopyEnvironmentVariables = \"False\"")
}

func TestSessionCacheIsNotActive(t *testing.T) {
	statusOutput := `
Credentials Cache:
  Mode: Auto
  Available for this process: False
  Total active cache sessions: 0`

	gsudo := Gsudo{}
	assert.False(t, gsudo.isCacheActive(statusOutput), "Session cache is not active")
}

func TestSessionCacheIsActive(t *testing.T) {
	statusOutput := `
	Credentials Cache:
	Mode: Auto
	Available for this process: True
	Total active cache sessions: 1
	  ProtectedPrefix\Administrators\gsudo_E47586A56563B34C06B7A20DB10A05A83245A9FEFB0D62C7187EEB48B157A9F1`

	gsudo := Gsudo{}
	assert.True(t, gsudo.isCacheActive(statusOutput), "Session cache is active")
}

func TestTooLongDefaultArgArray(t *testing.T) {
	_, err := isCheckError([]bool{false, false})
	assert.Error(t, err)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}
