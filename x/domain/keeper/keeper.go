package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

type ParamSubspace interface {
}

// list expected keepers

// ConfigurationKeeper defines the behaviour of the configuration state checks
type ConfigurationKeeper interface {
	// GetConfiguration returns the configuration
	GetConfiguration(ctx sdk.Context) configuration.Config
	// GetOwner returns the owner
	GetOwner(ctx sdk.Context) sdk.AccAddress
	// GetValidDomainRegexp returns the regular expression that a domain name must match
	// in order to be valid
	GetValidDomainRegexp(ctx sdk.Context) string
	// GetDomainRenewDuration returns the default duration of a domain renewal
	GetDomainRenewDuration(ctx sdk.Context) time.Duration
}

// Keeper of the domain store
// TODO split this keeper in sub-struct in order to avoid possible mistakes with keys and not clutter the exposed methods
type Keeper struct {
	// external keepers
	ConfigurationKeeper ConfigurationKeeper
	// default fields
	domainKey  sdk.StoreKey // contains the domain kvstore
	accountKey sdk.StoreKey // contains the account kvstore
	cdc        *codec.Codec
	paramspace ParamSubspace
}

// NewKeeper creates a domain keeper
func NewKeeper(cdc *codec.Codec, domainKey sdk.StoreKey, accountKey sdk.StoreKey, configKeeper ConfigurationKeeper, paramspace ParamSubspace) Keeper {
	keeper := Keeper{
		domainKey:           domainKey,
		accountKey:          accountKey,
		cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		paramspace:          nil,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
