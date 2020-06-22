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

	FeeCoinPriceKey = "fee_coin_price_key"
	FeeCoinDenom    = "fee_coin_denom"
	FeeDefault      = "fee_default"
)

var (
	// DomainStorePrefix is the prefix used to define the prefixed store containing domain data
	FeeStorePrefix = []byte{0x00}
)
