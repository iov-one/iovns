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

func (k Keeper) CollectFee(ctx sdk.Context, fee sdk.Coin, addr sdk.AccAddress) error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, auth.FeeCollectorName, sdk.NewCoins(fee))
}

// Logger provides logging facilities for Keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf(types.ModuleName))
}

func (k Keeper) feeStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) feeSeedStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(k.feeStore(ctx), types.FeeSeedPrefix)
}

// GetFees returns the network fees
func (k Keeper) GetFeeSeed(ctx sdk.Context, id string) *types.FeeSeed {
	store := k.feeSeedStore(ctx)
	value := store.Get([]byte(id))
	if value == nil {
		return nil
	}
	var amount sdk.Dec
	k.cdc.MustUnmarshalBinaryBare(value, &amount)
	return &types.FeeSeed{
		ID:     id,
		Amount: amount,
	}
}

func (k Keeper) GetAllFeeSeeds(ctx sdk.Context) []types.FeeSeed {
	store := k.feeSeedStore(ctx)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	var feeSeeds []types.FeeSeed
	for ; iterator.Valid(); iterator.Next() {
		var fs types.FeeSeed
		fsBytes := store.Get(iterator.Key())
		k.cdc.MustUnmarshalBinaryBare(fsBytes, &fs)
		feeSeeds = append(feeSeeds, fs)
	}
	return feeSeeds
}

func (k Keeper) GetFeeConfiguration(ctx sdk.Context) types.FeeConfiguration {
	feeSeeds := k.GetAllFeeSeeds(ctx)
	feeCfgr := k.GetFeeConfigurer(ctx)
	feeParams := k.GetFeeParams(ctx)
	return types.FeeConfiguration{
		FeeConfigurer: feeCfgr,
		FeeParamaters: feeParams,
		FeeSeeds:      feeSeeds,
	}
}

func (k Keeper) SetFeeSeed(ctx sdk.Context, fs types.FeeSeed) {
	store := k.feeSeedStore(ctx)
	store.Set([]byte(fs.ID), k.cdc.MustMarshalBinaryBare(fs.Amount))
}

func (k Keeper) GetFeeParams(ctx sdk.Context) types.FeeParamaters {
	store := k.feeStore(ctx)
	value := store.Get(types.FeeParametersKey)
	var fp types.FeeParamaters
	k.cdc.MustUnmarshalBinaryBare(value, &fp)
	return fp
}

func (k Keeper) SetFeeParameters(ctx sdk.Context, fs types.FeeParamaters) {
	store := k.feeStore(ctx)
	store.Set(types.FeeParametersKey, k.cdc.MustMarshalBinaryBare(fs))
}

func (k Keeper) GetFeeConfigurer(ctx sdk.Context) sdk.AccAddress {
	store := k.feeStore(ctx)
	value := store.Get(types.FeeConfigurerKey)
	var addr sdk.AccAddress
	k.cdc.MustUnmarshalBinaryBare(value, &addr)
	return addr
}

func (k Keeper) SetFeeConfigurer(ctx sdk.Context, cfgr sdk.AccAddress) {
	store := k.feeStore(ctx)
	store.Set(types.FeeConfigurerKey, k.cdc.MustMarshalBinaryBare(cfgr))
}
