package configuration

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/configuration/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

// configKey defines the key used for the configuration
// since the configuration is only one the key will always be one
const configKey = "config"

type paramSubspace interface {
}

type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace paramSubspace // TODO define what this is
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace paramSubspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramspace: paramspace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf(types.ModuleName))
}

// GetConfiguration returns the configuration of the blockchain
func (k Keeper) GetConfiguration(ctx sdk.Context) types.Config {
	store := ctx.KVStore(k.storeKey)
	confBytes := store.Get([]byte(configKey))
	if confBytes == nil {
		panic("no configuration available")
	}
	var conf types.Config
	k.cdc.MustUnmarshalBinaryBare(confBytes, &conf)
	// success
	return conf
}

// GetOwner returns the owner of domains with no superuser
func (k Keeper) GetOwner(ctx sdk.Context) sdk.AccAddress {
	return k.GetConfiguration(ctx).Owner
}

// GetDomainRenewDuration returns the duration of a domain renewal period
func (k Keeper) GetDomainRenewDuration(ctx sdk.Context) time.Duration {
	return time.Duration(k.GetConfiguration(ctx).DomainRenew) * time.Second
}

// GetValidDomainRegexp returns the regular expression used to match valid domain names
func (k Keeper) GetValidDomainRegexp(ctx sdk.Context) string {
	return k.GetConfiguration(ctx).ValidDomain
}

// SetConfig updates or saves a new config in the store
func (k Keeper) SetConfig(ctx sdk.Context, conf types.Config) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(configKey), k.cdc.MustMarshalBinaryBare(conf))
}
