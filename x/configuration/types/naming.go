package types

// ModuleConst
const (
	// ModuleName defines the name of the module
	ModuleName = "configuration"
	// StoreKey is the key used to identify the module in the KVStore
	StoreKey = ModuleName
	// RouterKey is the key used to process transactions for the module
	RouterKey = ModuleName
	// QuerierRoute is used to process queries for the module
	QuerierRoute = ModuleName
	// DefaultParamSpace defines the key for the configuration paramspace
	DefaultParamSpace = ModuleName
)

// QueryConfig is the route key used to query a config
const QueryConfig = "configuration"
