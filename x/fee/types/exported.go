package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProductFee interface {
	FeePayer() sdk.AccAddress
	FeeSeedPrefix() string
}

type ProductMsg interface {
	sdk.Msg
	ProductFee
}

type Calculator interface {
	CalculateFee(msg ProductMsg) (sdk.Coin, error)
}

type Collector interface {
	CollectFee(sdk.Context, sdk.Coin, sdk.AccAddress) error
}
