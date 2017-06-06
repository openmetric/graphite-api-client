package graphiteapi

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

// httpClient is used everywhere internally, here we initialize a default http client,
// library users can call graphiteapi.UseHTTPClient() to set a customized http client.
var httpClient *http.Client = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 1 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost: 5,
	},
}

var (
	userAgent     string              = "graphite-api-client/0.1"
	customHeaders map[string][]string = make(map[string][]string)
	chl           sync.RWMutex
)

// SetHTTPClient sets the http client used to make requests
func SetHTTPClient(client *http.Client) {
	httpClient = client
}

// AddCustomHeader adds a custom header on all requests
func AddCustomHeader(key, value string) {
	if key == "User-Agent" {
		SetUserAgent(value)
		return
	}

	chl.Lock()
	defer chl.Unlock()

	if values, ok := customHeaders[key]; ok {
		values = append(values, value)
		customHeaders[key] = values
	} else {
		values = []string{value}
		customHeaders[key] = values
	}
}

// SetUserAgent sets a custom user agent
func SetUserAgent(ua string) {
	userAgent = ua
}

// httpNewRequest wraps http.NewRequest(), and set custom headers
func httpNewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	if req, err := http.NewRequest(method, url, body); err != nil {
		return req, err
	} else {
		req.Header.Set("User-Agent", userAgent)
		chl.Lock()
		defer chl.Unlock()
		for key := range customHeaders {
			values := customHeaders[key]
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		return req, nil
	}
}

// httpDo wraps http.Client.Do(), fetches response and unmarshals into r
func httpDo(ctx context.Context, req *http.Request, r Response) error {
	req = req.WithContext(ctx)

	var resp *http.Response
	var body []byte
	var err error

	if resp, err = httpClient.Do(req); err != nil {
		return err
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	if err = r.Unmarshal(body); err != nil {
		return err
	}

	return nil
}
