package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/fee/types"
)

type DefaultCollector struct {
	types.Collector
	k Keeper
}

func NewCollector(k Keeper) DefaultCollector {
	return DefaultCollector{
		k: k,
	}
}

func (dc DefaultCollector) CollectFee(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error {
	return dc.k.CollectFee(ctx, fee, addr)
}
