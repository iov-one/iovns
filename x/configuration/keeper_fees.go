package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration/types"
)

// GetFees returns the network fees
func (k Keeper) GetFees(ctx sdk.Context) types.Fees {
	store := ctx.KVStore(k.storeKey)
	value := store.Get([]byte(feeKey))
	if value == nil {
		panic("no length fees set")
	}
	var fees types.Fees
	k.cdc.MustUnmarshalBinaryBare(value, &fees)
	return fees
}

// SetLengthFees sets the fee based on msg and length
func (k Keeper) SetLengthFees(ctx sdk.Context, msg sdk.Msg, length int, coin sdk.Coin) {
	fees := k.GetFees(ctx)
	fees.UpsertLengthFees(msg, length, coin)
	k.SetFees(ctx, fees)
}

// SetDefaultFees sets the default fees for a msg
func (k Keeper) SetDefaultFees(ctx sdk.Context, msg sdk.Msg, coin sdk.Coin) {
	fees := k.GetFees(ctx)
	fees.UpsertDefaultFees(msg, coin)
	k.SetFees(ctx, fees)
}

func (k Keeper) SetFees(ctx sdk.Context, fees types.Fees) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(feeKey), k.cdc.MustMarshalBinaryBare(fees))
}
