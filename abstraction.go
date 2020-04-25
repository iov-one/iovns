package iovns

// abstraction.go defines a series of interfaces to facilitate

// QueryHandler abstracts the functionality of a query handler
// CONTRACT: must be a struct pointer
type QueryHandler interface {
	// QueryPath defines the path of the query in the module to retrieve information
	QueryPath() string
	// Validate validates the correctness of the query formation in a stateless way
	Validate() error
}

// Command defines CLI commands functionality
// CONTRACT: must be a struct pointer
type Command interface {
	// Use defines the command string required to run the command
	Use() string
	// Description returns a string containing the information
	// of what the command actually does, ex: register a domain
	Description() string
}

// QueryCommand defines the query CLI commands functionality
// CONTRACT: must be a struct pointer
type QueryCommand interface {
	Command
	// QueryPath defines the path in which the query will be run
	// CONTRACT: it must not contain the module name, only the path
	// for the query inside the module.
	QueryPath() string
}
