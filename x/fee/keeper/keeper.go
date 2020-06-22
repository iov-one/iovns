package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/iov-one/iovns/x/fee/types"
	"github.com/tendermint/tendermint/libs/log"
)

// ParamSubspace is a placeholder
type ParamSubspace interface {
}

type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error
}

// Keeper is the key value store handler for the configuration module
type Keeper struct {
	supplyKeeper SupplyKeeper
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	paramspace   ParamSubspace
}

// NewKeeper is Keeper constructor
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper SupplyKeeper, paramspace params.Subspace) Keeper {
	return Keeper{
		supplyKeeper: supplyKeeper,
		storeKey:     key,
		cdc:          cdc,
		paramspace:   paramspace,
	}
}

// Logger provides logging facilities for Keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf(types.ModuleName))
}

func feeStore(store sdk.KVStore) sdk.KVStore {
	return prefix.NewStore(store, types.FeeStorePrefix)
}

func (k Keeper) CollectFee(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, auth.FeeCollectorName, sdk.NewCoins(fee))
}

// GetFees returns the network fees
func (k Keeper) GetFeeSeed(ctx sdk.Context, id string) types.FeeSeed {
	store := feeStore(ctx.KVStore(k.storeKey))
	value := store.Get([]byte(id))
	var d types.FeeSeed
	k.cdc.MustUnmarshalBinaryBare(value, &d)
	return d
}

func (k Keeper) GetFeeCoinPrice(ctx sdk.Context) sdk.Dec {
	store := feeStore(ctx.KVStore(k.storeKey))
	value := store.Get([]byte(types.FeeCoinPriceKey))
	var d sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(value, &d)
	return d
}

func (k Keeper) GetFeeCoinDenom(ctx sdk.Context) string {
	store := feeStore(ctx.KVStore(k.storeKey))
	value := store.Get([]byte(types.FeeCoinDenom))
	var cd string
	k.cdc.MustUnmarshalBinaryBare(value, &cd)
	return cd
}

func (k Keeper) GetDefaultFee(ctx sdk.Context) sdk.Dec {
	store := feeStore(ctx.KVStore(k.storeKey))
	value := store.Get([]byte(types.FeeDefault))
	var d sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(value, &d)
	return d
}

func (k Keeper) SetFee(ctx sdk.Context, id string, fees *types.FeeSeed) {
	store := feeStore(ctx.KVStore(k.storeKey))
	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(fees))
}
