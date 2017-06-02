package graphiteapi

import (
	"fmt"
	"strings"
)

// RenderQuery constructs a '/render/' query. RenderQuery implements Query
type RenderQuery struct {
	Targets       []string
	From          interface{}
	Until         interface{}
	MaxDataPoints int
}

// RenderResponse represents render query response
type RenderResponse struct {
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

func (q *RenderQuery) SetFrom(from interface{}) *RenderQuery {
	q.From = from
	return q
}

func (q *RenderQuery) SetUntil(until interface{}) *RenderQuery {
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
func (q *RenderQuery) URL(urlbase string, format string) string {
	args := []string{}
	for _, t := range q.Targets {
		args = append(args, "target="+t)
	}
	if q.From != nil {
		args = append(args, fmt.Sprintf("from=%v", q.From))
	}
	if q.Until != nil {
		args = append(args, fmt.Sprintf("until=%v", q.Until))
	}
	if q.MaxDataPoints != 0 {
		args = append(args, fmt.Sprintf("maxDataPoints=%d", q.MaxDataPoints))
	}
	if format != "" {
		args = append(args, fmt.Sprintf("format=%s", format))
	}
	return fmt.Sprintf("%s/render/?%s", urlbase, strings.Join(args, "&"))
}

// Request implements Query interface
func (q *RenderQuery) Request(urlbase string) Response {
	return &RenderResponse{}
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

// function shortcuts, for code completion

func (t *QueryTarget) SumSeries() *QueryTarget {
	t.str = formatFunction("sumSeries", t.str)
	return t
}

func (t *QueryTarget) ConstantLine(value interface{}) *QueryTarget {
	t.str = formatFunction("constantLine", value)
	return t
}
