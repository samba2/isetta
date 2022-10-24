package httpchecker

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO rename proxyUrl  -> pxProxyUrl

func TestDirectInternetAccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	httpChecker, err := New(ts.URL, "")
	assert.NoError(t, err)
	httpChecker.TimeoutInMilliseconds = 100
	assert.True(t, httpChecker.HasDirectInternetAccess())
}

func TestExitOnWrongAddress(t *testing.T) {
	httpChecker, err := New("http://non-existing", "")
	assert.NoError(t, err)
	httpChecker.TimeoutInMilliseconds = 100
	assert.False(t, httpChecker.HasDirectInternetAccess())
}

func TestInternetAccessViaProxy(t *testing.T) {
	// arrange
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ts.Close()

	targetUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	proxyHandler := httputil.NewSingleHostReverseProxy(targetUrl)
	proxy := httptest.NewServer(proxyHandler)
	defer proxy.Close()

	// act
	httpChecker, err := New(ts.URL, proxy.URL)
	assert.NoError(t, err)
	httpChecker.TimeoutInMilliseconds = 100

	// assert
	assert.True(t, httpChecker.HasInternetAccessViaProxy())
}

func TestInternetAccessViaProxyFailed(t *testing.T) {
	httpChecker, err := New("http://foo", "http://127.0.0.1:1")
	assert.NoError(t, err)
	httpChecker.TimeoutInMilliseconds = 100

	assert.False(t, httpChecker.HasInternetAccessViaProxy())
}

func TestInvalidProxyAddress(t *testing.T) {
	invalidProxyUrl := ":" // parser is really forgiving, this one works ;-)
	_, err := New("http://foo", invalidProxyUrl)
	assert.Error(t, err)
}

func TestPxProxyReachable(t *testing.T) {
	fakePxProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer fakePxProxy.Close()

	httpChecker, err := New("", fakePxProxy.URL)
	assert.NoError(t, err)
	assert.True(t, httpChecker.IsPxProxyReachable())
}

func TestPxProxyNotReachable(t *testing.T) {
	httpChecker, err := New("", "http://127.0.0.1:9999")
	assert.NoError(t, err)
	assert.False(t, httpChecker.IsPxProxyReachable())
}
