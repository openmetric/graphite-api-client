package graphiteapi

import (
	"context"
	pb "github.com/go-graphite/carbonzipper/carbonzipperpb3"
)

// Query is interface for all api request
type Query interface {
	URL() string
	Request(ctx context.Context) (Response, error)
}

// Response is interface for all api request response types
type Response interface {
	Unmarshal([]byte) error
}

// RenderQuery is used to build `/render/` query
type RenderQuery struct {
	Base          string // base url of graphite server
	Targets       []*RenderTarget
	From          string
	Until         string
	MaxDataPoints int
}

// RenderResponse is response of `/render/` query
type RenderResponse struct {
	pb.MultiFetchResponse
}

// QueryTarget represents a `target=` arg in `/render/` query
type RenderTarget struct {
	str string
}

// GetLastNonNullValue searches for the latest non null value, and skips at most maxNullPoints.
// If the last maxNullPoints values are all absent, returns absent
func GetLastNonNullValue(m *pb.FetchResponse, maxNullPoints int) (v float64, t int32, absent bool) {
	l := len(m.Values)
	for i := 0; i < maxNullPoints && i < l; i++ {
		if m.IsAbsent[l-1-i] {
			continue
		}
		v = m.Values[l-1-i]
		t = m.StopTime - int32(i)*m.StepTime
		absent = false
		return v, t, absent
	}
	// if we didn't return in the loop body, there were too many null points
	v = 0
	t = m.StopTime
	absent = true
	return v, t, absent
}
