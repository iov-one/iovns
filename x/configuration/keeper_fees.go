package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration/types"
)

// GetFees returns the network fees
func (k Keeper) GetFees(ctx sdk.Context) *types.Fees {
	store := ctx.KVStore(k.storeKey)
	value := store.Get([]byte(types.FeeKey))
	if value == nil {
		panic("no length fees set")
	}
	var fees = new(types.Fees)
	k.cdc.MustUnmarshalBinaryBare(value, fees)
	return fees
}

func (k Keeper) SetFees(ctx sdk.Context, fees *types.Fees) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.FeeKey), k.cdc.MustMarshalBinaryBare(fees))
}
