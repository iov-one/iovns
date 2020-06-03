package configuration

import (
	"encoding/json"

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
	err := json.Unmarshal(value, fees)
	if err != nil {
		panic(err)
	}
	return fees
}

// SetLengthFees sets the fee based on msg and length
func (k Keeper) SetLengthFees(ctx sdk.Context, msg sdk.Msg, length int, coin sdk.Coin) {
	fees := k.GetFees(ctx)
	fees.UpsertLevelFees(msg, length, coin)
	k.SetFees(ctx, fees)
}

// SetDefaultFees sets the default fees for a msg
func (k Keeper) SetDefaultFees(ctx sdk.Context, msg sdk.Msg, coin sdk.Coin) {
	fees := k.GetFees(ctx)
	fees.UpsertDefaultFees(msg, coin)
	k.SetFees(ctx, fees)
}

func (k Keeper) SetFees(ctx sdk.Context, fees *types.Fees) {
	store := ctx.KVStore(k.storeKey)
	b, err := json.Marshal(fees)
	if err != nil {
		panic(err)
	}
	store.Set([]byte(types.FeeKey), b)
}
