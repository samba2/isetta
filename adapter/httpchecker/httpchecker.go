package httpchecker

import (
	"net/http"
	"net/url"
	"os"
	"time"

	log "org.samba/isetta/simplelogger"
)

type HttpCheckerImpl struct {
	InternetAccessTestUrl        string
	ProxyUrl                     *url.URL
	DefaultTimeoutInMilliseconds int
}

func New(InternetAccessTestUrl string, proxyUrl string) (HttpCheckerImpl, error) {
	proxyUrl2, err := url.Parse(proxyUrl)
	if err != nil {
		return HttpCheckerImpl{}, err
	}
	return HttpCheckerImpl{
		InternetAccessTestUrl:        InternetAccessTestUrl,
		ProxyUrl:                     proxyUrl2,
		DefaultTimeoutInMilliseconds: 5000,
	}, nil
}

func (h *HttpCheckerImpl) HasDirectInternetAccess(timeoutInMilliseconds ...int) bool {
	os.Unsetenv("https_proxy")
	os.Unsetenv("HTTPS_PROXY")

	timeout := determineTimeout(timeoutInMilliseconds, h.DefaultTimeoutInMilliseconds)
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	resp, err := client.Get(h.InternetAccessTestUrl)
	if err != nil {
		log.Logger.Info("Unable to directly access %v", h.InternetAccessTestUrl)
		log.Logger.Debug("Error was: %v", err)
		return false
	}

	if resp.StatusCode == 200 {
		log.Logger.Debug("Successfully connected directly to %v", h.InternetAccessTestUrl)
		return true
	} else {
		log.Logger.Debug("HTTP error when trying to directly connect to %v. HTTP status code was: %v", h.InternetAccessTestUrl, resp.StatusCode)
		return false
	}
}

func (h *HttpCheckerImpl) HasInternetAccessViaProxy(timeoutInMilliseconds ...int) bool {
	timeout := determineTimeout(timeoutInMilliseconds, h.DefaultTimeoutInMilliseconds)
	httpClientWithProxy := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(h.ProxyUrl)},
		Timeout:   time.Duration(timeout) * time.Millisecond,
	}
	resp, err := httpClientWithProxy.Get(h.InternetAccessTestUrl)

	if err != nil {
		log.Logger.Info("Unable to access %v via proxy", h.InternetAccessTestUrl)
		log.Logger.Debug("Error was: %v", err)
		return false
	}

	if resp.StatusCode == 200 {
		log.Logger.Debug("Successfully connected to %v via proxy %v", h.InternetAccessTestUrl, h.ProxyUrl)
		return true
	} else {
		log.Logger.Debug("HTTP error when connecting to %v via proxy %v. HTTP status code was: %v", h.InternetAccessTestUrl, h.ProxyUrl, resp.StatusCode)
		return false
	}
}

func determineTimeout(timeoutInMilliseconds []int, defaultTimeoutInMilliseconds int) int {
	if len(timeoutInMilliseconds) > 0 {
		return timeoutInMilliseconds[0]
	} else {
		return defaultTimeoutInMilliseconds
	}
}

func (h *HttpCheckerImpl) IsPxProxyReachable() bool {
	client := http.Client{
		Timeout: time.Duration(h.DefaultTimeoutInMilliseconds) * time.Millisecond,
	}

	_, err := client.Get(h.ProxyUrl.String())
	return err == nil
}
