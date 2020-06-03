package configuration

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/tendermint/tendermint/libs/log"
)

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
	confBytes := store.Get([]byte(types.ConfigKey))
	if confBytes == nil {
		panic("no configuration available")
	}
	var conf types.Config
	k.cdc.MustUnmarshalBinaryBare(confBytes, &conf)
	// success
	return conf
}

// GetConfigurer returns the owner of domains with no superuser
func (k Keeper) GetConfigurer(ctx sdk.Context) sdk.AccAddress {
	return k.GetConfiguration(ctx).Configurer
}

// IsOwner checks if the provided address is an owner or not
func (k Keeper) IsOwner(ctx sdk.Context, addr sdk.AccAddress) bool {
	configurer := k.GetConfigurer(ctx)
	return configurer.Equals(addr)
}

// GetDomainRenewDuration returns the duration of a domain renewal period
func (k Keeper) GetDomainRenewDuration(ctx sdk.Context) time.Duration {
	return k.GetConfiguration(ctx).DomainRenewalPeriod
}

// GetValidDomainNameRegexp returns the regular expression used to match valid domain names
func (k Keeper) GetValidDomainNameRegexp(ctx sdk.Context) string {
	return k.GetConfiguration(ctx).ValidDomainName
}

// SetConfig updates or saves a new config in the store
func (k Keeper) SetConfig(ctx sdk.Context, conf types.Config) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.ConfigKey), k.cdc.MustMarshalBinaryBare(conf))
}

// GetDomainGracePeriod returns the default grace period before domains
// can be deleted by someone other than the owner him/herself
func (k Keeper) GetDomainGracePeriod(ctx sdk.Context) time.Duration {
	return k.GetConfiguration(ctx).DomainGracePeriod
}
