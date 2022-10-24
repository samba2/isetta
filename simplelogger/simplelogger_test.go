package simplelogger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidLogLevels(t *testing.T) {

	validLogLevels := GetValidLogLevels()
	assert.NotEmpty(t, validLogLevels)
	assert.Contains(t, validLogLevels, "trace")
}