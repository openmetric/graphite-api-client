package graphiteapi

// Query is interface for all api queries, it is used to construct http query args
type Query interface {
	// Returns full request URL of this Query.
	URL(string) string
	// Sends requests to Graphite
	Request(string) *Response
}

// Response is interface for all api request responses
type Response interface {
}
