package windows

// tests here run with all unit tests and don't need a Windows/ WSL underneath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseListOutputOk(t *testing.T) {
	output := `
NAME      STATE           VERSION
* Ubuntu    Running         2
  foo
`

	assert.True(t, parseListOutput(output))
}

func TestParseListOutputNotOk(t *testing.T) {
	output := `
NAME      STATE           VERSION
* Ubuntu    Running         1
  foo
`

	assert.False(t, parseListOutput(output))
}