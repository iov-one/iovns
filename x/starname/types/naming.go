package types

// Module names
const (
	// ModuleName is the name of the module
	ModuleName = "starname"
	// DomainStore key defines the store key used to store domains information
	DomainStoreKey = "starname"
	// RouterKey defines the path used to interact with the domain module
	RouterKey    = ModuleName
	QuerierAlias = "starname"
	// QuerierRoute defines the query path used to interact with the domain module
	QuerierRoute = ModuleName
	// DefaultParamSpace defines the key for the default param space
	DefaultParamSpace = ModuleName
)

// Events attribute keys
const (
	AttributeKeyDomainName = "domain_name"
	AttributeKeyDomainType = "domain_type"
)
