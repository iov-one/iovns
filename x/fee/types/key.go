package types

const (
	// ModuleName is the name of the module
	ModuleName = "fee"
	StoreKey   = ModuleName
	// RouterKey is the key used to process transactions for the module
	RouterKey = ModuleName
	// QuerierRoute is used to process queries for the module
	QuerierRoute = ModuleName

	DefaultParamSpace = ModuleName
)
const (
	// FeeKey defines the key used for fees
	// since the fee params are only one
	// this is the only key we will need
	FeeKey = "fee"
)
