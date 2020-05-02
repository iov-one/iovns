package configuration

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

// configKey defines the key used for the configuration
// since the configuration is only one the key will always be one
const configKey = "config"

// feeKey defines the key used for fees
// since the fee params are only one
// this is the only key we will need
const feeKey = "fee"

// Keeper is the key value store handler for the configuration module
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace params.Subspace
}

// NewKeeper is Keeper constructor
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramspace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramspace: paramspace,
	}
}

// Logger provides logging facilities for Keeper
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

// GetOwners returns the owner of domains with no superuser
func (k Keeper) GetOwners(ctx sdk.Context) []sdk.AccAddress {
	return k.GetConfiguration(ctx).Owners
}

// IsOwner checks if the provided address is an owner or not
func (k Keeper) IsOwner(ctx sdk.Context, addr sdk.AccAddress) bool {
	owners := k.GetOwners(ctx)
	for _, owner := range owners {
		if owner.Equals(addr) {
			return true
		}
	}
	return false
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

// GetDomainGracePeriod returns the default grace period before domains
// can be deleted by someone other than the owner him/herself
func (k Keeper) GetDomainGracePeriod(ctx sdk.Context) time.Duration {
	return time.Duration(k.GetConfiguration(ctx).DomainGracePeriod) * time.Second
}
