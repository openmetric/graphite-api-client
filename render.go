package graphiteapi

import (
	"fmt"
	pb "github.com/go-graphite/carbonzipper/carbonzipperpb3"
	"net/url"
	"strconv"
	"strings"
)

// RenderQuery constructs a '/render/' query. RenderQuery implements Query
type RenderQuery struct {
	Targets       []string
	From          string
	Until         string
	MaxDataPoints int
}

// RenderResponse represents render query response
type RenderResponse struct {
	pb.MultiFetchResponse
}

// NewRenderQuery returns a RenderQuery instance
func NewRenderQuery(targets ...*QueryTarget) *RenderQuery {
	t := make([]string, len(targets))
	for i, target := range targets {
		t[i] = target.String()
	}
	q := &RenderQuery{
		Targets: t,
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

func (q *RenderQuery) AddTarget(target *QueryTarget) *RenderQuery {
	q.Targets = append(q.Targets, target.String())
	return q
}

func (q *RenderQuery) SetMaxDataPoints(maxDataPoints int) *RenderQuery {
	q.MaxDataPoints = maxDataPoints
	return q
}

// URL implements Query interface
func (q *RenderQuery) URL(urlbase string, format string) *url.URL {
	u, _ := url.Parse(urlbase + "/render/")
	v := url.Values{}

	for _, target := range q.Targets {
		v.Add("target", target)
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

	if format != "" {
		v.Set("format", format)
	}

	u.RawQuery = v.Encode()

	return u
}

// Request implements Query interface
func (q *RenderQuery) Request(urlbase string) (*RenderResponse, error) {
	u := q.URL(urlbase, "protobuf")

	response := &RenderResponse{}
	if err := get(u, &response.MultiFetchResponse); err != nil {
		return nil, err
	}

	return response, nil
}

// Target represents a "?target=" in query
type QueryTarget struct {
	str string
}

func (t *QueryTarget) String() string {
	return t.str
}

func NewQueryTarget(seriesList string) *QueryTarget {
	return &QueryTarget{
		str: seriesList,
	}
}

func formatFunction(name string, args ...interface{}) string {
	t := make([]string, len(args))
	for i, a := range args {
		t[i] = fmt.Sprintf("%v", a)
	}
	return fmt.Sprintf("%s(%s)", name, strings.Join(t, ","))
}

func (t *QueryTarget) applyFunction(name string, args ...interface{}) *QueryTarget {
	t.str = formatFunction(name, args...)
	return t
}

// function shortcuts, for code completion

func (t *QueryTarget) SumSeries() *QueryTarget {
	return t.applyFunction("sumSeries", t.str)
}

func (t *QueryTarget) ConstantLine(value interface{}) *QueryTarget {
	return t.applyFunction("constantLine", value)
}
