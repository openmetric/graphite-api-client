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
