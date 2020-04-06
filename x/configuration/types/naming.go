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
)

// Query Routes const
const QueryConfig = "configuration"
