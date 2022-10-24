package httpchecker

import (
	"net/http"
	"net/url"
	"os"
	"time"

	log "org.samba/isetta/simplelogger"
)

type HttpCheckerImpl struct {
	InternetAccessTestUrl string
	ProxyUrl              *url.URL
	TimeoutInMilliseconds int
}

func New(InternetAccessTestUrl string, proxyUrl string) (HttpCheckerImpl, error) {
	proxyUrl2, err := url.Parse(proxyUrl)
	if err != nil {
		return HttpCheckerImpl{}, err
	}
	return HttpCheckerImpl{
		InternetAccessTestUrl: InternetAccessTestUrl,
		ProxyUrl:              proxyUrl2,
		TimeoutInMilliseconds: 2000,
	}, nil
}

func (h *HttpCheckerImpl) HasDirectInternetAccess() bool {
	os.Unsetenv("https_proxy")
	os.Unsetenv("HTTPS_PROXY")

	client := http.Client{
		Timeout: time.Duration(h.TimeoutInMilliseconds) * time.Millisecond,
	}

	resp, err := client.Get(h.InternetAccessTestUrl)
	if err != nil {
		log.Logger.Warn("Error when trying access %v, error was: %v", h.InternetAccessTestUrl, err)
		return false
	}

	return resp.StatusCode == 200
}

func (h *HttpCheckerImpl) HasInternetAccessViaProxy() bool {
	httpClientWithProxy := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(h.ProxyUrl)},
		Timeout:   time.Duration(h.TimeoutInMilliseconds) * time.Millisecond,
	}
	resp, err := httpClientWithProxy.Get(h.InternetAccessTestUrl)

	if err != nil {
		log.Logger.Warn("Error when trying access %v, error was: %v", h.InternetAccessTestUrl, err)
		return false
	}

	if resp.StatusCode == 200 {
		log.Logger.Debug("Successfully connected to %v via proxy %v", h.InternetAccessTestUrl, h.ProxyUrl)
		return true
	} else {
		log.Logger.Warn("Failed to connected to %v via proxy %v", h.InternetAccessTestUrl, h.ProxyUrl)
		return false
	}
}

func (h *HttpCheckerImpl) IsPxProxyReachable() bool {
	client := http.Client{
		Timeout: time.Duration(h.TimeoutInMilliseconds) * time.Millisecond,
	}

	_, err := client.Get(h.ProxyUrl.String())
	return err == nil
}
