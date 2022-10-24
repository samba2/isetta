// +build wsl
// +build interactive

package gsudo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// should only ask 1x for a password
func TestRunElevated(t *testing.T) {
	gsudo := Gsudo{}
	gsudo.Init()
	defer gsudo.Cleanup()

	// this asks for admin credentials
	out := gsudo.RunElevated("net session && echo I am elevated")
	assert.Contains(t, out, "I am elevated")

	// this doesn't, the credential cache is used
	out = gsudo.RunElevated("net session && echo I am elevated")
	assert.Contains(t, out, "I am elevated")

	// disable cache
	gsudo.run("--reset-timestamp")
}
