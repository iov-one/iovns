package types

// Module names
const (
	ModuleName      = "domain"
	DomainStoreKey  = "domain"
	AccountStoreKey = "account"
	RouterKey       = ModuleName
	QuerierRoute    = ModuleName
)

// Module Queries
const (
	// QueryDomain is the query route used to get a domain by its name
	QueryDomain = "get"
)
