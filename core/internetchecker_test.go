package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"org.samba/isetta/mocks"
)

var internetChecker InternetChecker

func setupInternetChecker(t *testing.T) {
	mockHttpChecker = mocks.NewHttpChecker(t)

	internetChecker = InternetChecker{
		HttpChecker:     mockHttpChecker,
	}
}

func TestHasAccessViaDirectConnection(t *testing.T) {
	setupInternetChecker(t)
	mockHttpChecker.On("HasDirectInternetAccess").Return(true)
	assert.True(t, internetChecker.HasInternetAccess())
}

func TestHasAccessViaProxyConnection(t *testing.T) {
	setupInternetChecker(t)
	mockHttpChecker.On("HasDirectInternetAccess").Return(false)
	mockHttpChecker.On("HasInternetAccessViaProxy").Return(true)
	assert.True(t, internetChecker.HasInternetAccess())
}

func TestHasNoAccess(t *testing.T) {
	setupInternetChecker(t)
	mockHttpChecker.On("HasDirectInternetAccess").Return(false)
	mockHttpChecker.On("HasInternetAccessViaProxy").Return(false)
	assert.False(t, internetChecker.HasInternetAccess())
}
