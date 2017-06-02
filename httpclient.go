package graphiteapi

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
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

type unmarshaler interface {
	Unmarshal([]byte) error
}

func get(u *url.URL, r unmarshaler) error {
	var err error
	var req *http.Request
	var resp *http.Response
	var body []byte

	if req, err = newHttpRequest("GET", u.String(), nil); err != nil {
		return err
	}

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
