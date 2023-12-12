package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"org.samba/isetta/mocks"
)

var internetChecker InternetCheckerImpl

func setupInternetChecker(t *testing.T) {
	mockHttpChecker = mocks.NewHttpChecker(t)

	internetChecker = InternetCheckerImpl{
		HttpChecker: mockHttpChecker,
		TimeoutInMilliseconds: 100,
	}
}

func TestHasAccessViaDirectConnection(t *testing.T) {
	setupInternetChecker(t)
	mockHttpChecker.On("HasDirectInternetAccess", 100).Return(true)
	mockHttpChecker.On("HasInternetAccessViaProxy", 100).Return(false)
	assert.True(t, internetChecker.HasInternetAccess())
}

func TestHasAccessViaProxyConnection(t *testing.T) {
	setupInternetChecker(t)
	mockHttpChecker.On("HasDirectInternetAccess", 100).Return(false)
	mockHttpChecker.On("HasInternetAccessViaProxy", 100).Return(true)
	assert.True(t, internetChecker.HasInternetAccess())
}

func TestHasNoAccess(t *testing.T) {
	setupInternetChecker(t)
	mockHttpChecker.On("HasDirectInternetAccess", 100).Return(false)
	mockHttpChecker.On("HasInternetAccessViaProxy", 100).Return(false)
	assert.False(t, internetChecker.HasInternetAccess())
}

