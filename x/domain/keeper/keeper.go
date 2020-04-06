package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/domain/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the domain store
type Keeper struct {
	// external keepers
	ConfigurationKeeper types.ConfigurationKeeper
	AccountKeeper       types.AccountKeeper
	// default fields
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
}

// NewKeeper creates a domain keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, accountKeeper types.AccountKeeper, configKeeper types.ConfigurationKeeper, paramspace types.ParamSubspace) Keeper {
	keeper := Keeper{
		storeKey:            key,
		cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		AccountKeeper:       accountKeeper,
		paramspace:          paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetDomain returns the domain based on its name, if domain is not found ok will be false
func (k Keeper) GetDomain(ctx sdk.Context, domainName string) (domain types.Domain, ok bool) {
	store := ctx.KVStore(k.storeKey)
	// get domain in form of bytes
	domainBytes := store.Get([]byte(domainName))
	// if nothing is returned, return nil
	if domainBytes == nil {
		return
	}
	// if domain exists then unmarshal
	k.cdc.MustUnmarshalBinaryBare(domainBytes, &domain)
	// success
	return domain, true
}

// SetDomain saves the domain inside the KVStore with its name as key
func (k Keeper) SetDomain(ctx sdk.Context, domain types.Domain) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(domain.Name), k.cdc.MustMarshalBinaryBare(domain))
}

func (k Keeper) delete(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}
