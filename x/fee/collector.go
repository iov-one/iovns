package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CollectorI interface {
	CollectFee(sdk.Context, sdk.Coin, sdk.AccAddress) error
}

type Collector struct {
	CollectorI
	k Keeper
}

func NewCollector(k Keeper) Collector {
	return Collector{
		k: k,
	}
}

func (gc Collector) CollectFee(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error {
	return gc.k.CollectFee(ctx, fee, addr)
}
