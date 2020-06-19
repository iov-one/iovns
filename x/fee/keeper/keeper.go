package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/fee/types"
	"github.com/tendermint/tendermint/libs/log"
)

// ParamSubspace is a placeholder
type ParamSubspace interface {
}

// ConfigurationKeeper defines the behaviour of the configuration state checks
type ConfigurationKeeper interface {
	// GetConfiguration returns the configuration
	GetConfiguration(ctx sdk.Context) configuration.Config
	// IsOwner returns if the provided address is an owner or not
	IsOwner(ctx sdk.Context, addr sdk.AccAddress) bool
	// GetValidDomainNameRegexp returns the regular expression that aliceAddr domain name must match
	// in order to be valid
	GetValidDomainNameRegexp(ctx sdk.Context) string
	// GetDomainRenewDuration returns the default duration of aliceAddr domain renewal
	GetDomainRenewDuration(ctx sdk.Context) time.Duration
	// GetDomainGracePeriod returns the grace period duration
	GetDomainGracePeriod(ctx sdk.Context) time.Duration
}

// Keeper is the key value store handler for the configuration module
type Keeper struct {
	ConfigurationKeeper ConfigurationKeeper
	storeKey            sdk.StoreKey
	cdc                 *codec.Codec
	paramspace          ParamSubspace
}

// NewKeeper is Keeper constructor
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, configKeeper ConfigurationKeeper, paramspace params.Subspace) Keeper {
	return Keeper{
		ConfigurationKeeper: configKeeper,
		storeKey:            key,
		cdc:                 cdc,
		paramspace:          paramspace,
	}
}

// Logger provides logging facilities for Keeper
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf(types.ModuleName))
}

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
