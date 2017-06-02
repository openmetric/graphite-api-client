package graphiteapi

import (
	"io"
	"net"
	"net/http"
	"time"
)

var httpClient *http.Client

const userAgent = "graphite-api-client/0.1"

func init() {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 1 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: 5,
	}
	httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func newHttpRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}
