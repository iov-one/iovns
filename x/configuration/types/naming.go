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

// QueryConfig is the route key used to query configuration data
const QueryConfig = "configuration"

// QueryFees is the route key used to query fees data
const QueryFees = "fees"

const (
	// ConfigKey defines the key used for the configuration
	// since the configuration is only one the key will always be one
	ConfigKey = "config"

	// FeeKey defines the key used for fees
	// since the fee params are only one
	// this is the only key we will need
	FeeKey = "fee"
)
