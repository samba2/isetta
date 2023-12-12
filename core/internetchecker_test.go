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

func TestInternetChecker(t *testing.T) {
	testCases := []struct {
		desc	string
		hasDirectInternetAccess bool
		hasInternetAccessViaProxy bool
		hasInternetAccess bool
	}{
		{
			desc: "has access via direct connection",
			hasDirectInternetAccess: true,
			hasInternetAccessViaProxy: false,
			hasInternetAccess: true,
		},
		{
			desc: "has access via proxy connection",
			hasDirectInternetAccess: false,
			hasInternetAccessViaProxy: true,
			hasInternetAccess: true,
		},
		{
			desc: "has not access",
			hasDirectInternetAccess: false,
			hasInternetAccessViaProxy: false,
			hasInternetAccess: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			setupInternetChecker(t)
			mockHttpChecker.On("HasDirectInternetAccess", 100).Return(tC.hasDirectInternetAccess)
			mockHttpChecker.On("HasInternetAccessViaProxy", 100).Return(tC.hasInternetAccessViaProxy)
			assert.Equal(t, tC.hasInternetAccess, internetChecker.HasInternetAccess())
		})
	}
}
