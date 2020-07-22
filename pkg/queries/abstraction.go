package queries

// abstraction.go defines a series of interfaces to facilitate

// QueryHandler abstracts the functionality of a query handler
// CONTRACT: must be a struct pointer
type QueryHandler interface {
	// QueryPath defines the path of the query in the module to retrieve information
	QueryPath() string
	// Validate validates the correctness of the query formation in a stateless way
	Validate() error
}
