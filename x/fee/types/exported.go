package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
)

type ProductFee interface {
	CalculateFee(calculator types.FeeCalculator) (sdk.Dec, error)
	FeePayer() sdk.AccAddress
}

type ProductMsg interface {
	sdk.Msg
	ProductFee
}

type Calculator interface {
	CalculateFee(msg ProductMsg) (sdk.Coin, error)
}
