package graphiteapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// NewRenderQuery returns a RenderQuery instance
func NewRenderQuery(base, from, until string, targets ...*RenderTarget) *RenderQuery {
	q := &RenderQuery{
		Base:          base,
		Targets:       targets,
		From:          from,
		Until:         until,
		MaxDataPoints: 0,
	}
	return q
}

func (q *RenderQuery) SetFrom(from string) *RenderQuery {
	q.From = from
	return q
}

func (q *RenderQuery) SetUntil(until string) *RenderQuery {
	q.Until = until
	return q
}

func (q *RenderQuery) AddTarget(target *RenderTarget) *RenderQuery {
	q.Targets = append(q.Targets, target)
	return q
}

func (q *RenderQuery) SetMaxDataPoints(maxDataPoints int) *RenderQuery {
	q.MaxDataPoints = maxDataPoints
	return q
}

// URL implements Query interface
func (q *RenderQuery) URL() *url.URL {
	u, _ := url.Parse(q.Base + "/render/")
	v := url.Values{}

	// force set format to protobuf
	v.Set("format", "protobuf")

	for _, target := range q.Targets {
		v.Add("target", target.String())
	}

	if q.From != "" {
		v.Set("from", q.From)
	}

	if q.Until != "" {
		v.Set("until", q.Until)
	}

	if q.MaxDataPoints != 0 {
		v.Set("maxDataPoints", strconv.Itoa(q.MaxDataPoints))
	}

	u.RawQuery = v.Encode()

	return u
}

// Request implements Query interface
func (q *RenderQuery) Request(ctx context.Context) (*RenderResponse, error) {
	var req *http.Request
	var err error
	response := &RenderResponse{}

	if req, err = httpNewRequest("GET", q.URL().String(), nil); err != nil {
		return nil, err
	}

	if err = httpDo(ctx, req, &response.MultiFetchResponse); err != nil {
		return nil, err
	}

	return response, nil
}

func (t *RenderTarget) String() string {
	return t.str
}

func NewRenderTarget(seriesList string) *RenderTarget {
	return &RenderTarget{
		str: seriesList,
	}
}

func (t *RenderTarget) ApplyFunction(name string, args ...interface{}) *RenderTarget {
	tmp := make([]string, len(args)+1)
	tmp[0] = t.String()
	for i, a := range args {
		tmp[i+1] = fmt.Sprintf("%v", a)
	}
	t.str = fmt.Sprintf("%s(%s)", name, strings.Join(tmp, ","))
	return t
}

func (t *RenderTarget) ApplyFunctionWithoutSeries(name string, args ...interface{}) *RenderTarget {
	tmp := make([]string, len(args))
	for i, a := range args {
		tmp[i] = fmt.Sprintf("%v", a)
	}
	t.str = fmt.Sprintf("%s(%s)", name, strings.Join(tmp, ","))
	return t
}

//
// function shortcuts, for code completion
//

func (t *RenderTarget) SumSeries() *RenderTarget {
	return t.ApplyFunction("sumSeries")
}

func (t *RenderTarget) ConstantLine(value interface{}) *RenderTarget {
	return t.ApplyFunctionWithoutSeries("constantLine", value)
}
