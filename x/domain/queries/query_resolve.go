package queries

// RequestBuilder defines the behaviour of a request
// that gets filled from either bytes, coming from example
// from a request body or arguments coming from the CLI
type RequestBuilder interface {
	// FromArgs builds the request from command line arguments
	FromArgs([]string) error
	// Unmarshal builds the request from bytes
	Unmarshal([]byte) error
}

// ResponseHandler takes care of unmarshalling the response
// coming from queries
type ResponseHandler interface {
	// Unmarshal fills the ResponseHandlers fields with the responses body
	Unmarshal([]byte) error
	// Marshal marshals the ResponseHandler to bytes
	Marshal() ([]byte, error)
}

// Query defines a query type in the
type Query struct {
	Builder RequestBuilder
	Handler ResponseHandler
}
