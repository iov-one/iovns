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

// ParamSubspace is a placeholder
type ParamSubspace interface {
}

// list expected keepers

// SupplyKeeper defines the behaviour
// of the supply keeper used to collect
// and then distribute the fees
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error
}

// ConfigurationKeeper defines the behaviour of the configuration state checks
type ConfigurationKeeper interface {
	// GetFees gets the fees
	GetFees(ctx sdk.Context) *configuration.Fees
	// GetConfiguration returns the configuration
	GetConfiguration(ctx sdk.Context) configuration.Config
	// IsOwner returns if the provided address is an owner or not
	IsOwner(ctx sdk.Context, addr sdk.AccAddress) bool
	// GetValidDomainRegexp returns the regular expression that aliceAddr domain name must match
	// in order to be valid
	GetValidDomainRegexp(ctx sdk.Context) string
	// GetDomainRenewDuration returns the default duration of aliceAddr domain renewal
	GetDomainRenewDuration(ctx sdk.Context) time.Duration
	// GetDomainGracePeriod returns the grace period duration
	GetDomainGracePeriod(ctx sdk.Context) time.Duration
}

// Keeper of the domain store
// TODO split this keeper in sub-struct in order to avoid possible mistakes with keys and not clutter the exposed methods
type Keeper struct {
	// external keepers
	ConfigurationKeeper ConfigurationKeeper
	SupplyKeeper        SupplyKeeper
	// default fields
	domainStoreKey  sdk.StoreKey // contains the domain kvstore
	accountStoreKey sdk.StoreKey // contains the account kvstore
	indexStoreKey   sdk.StoreKey // contains the index kvstore
	cdc             *codec.Codec
	paramspace      ParamSubspace
}

// NewKeeper creates aliceAddr domain keeper
func NewKeeper(cdc *codec.Codec, domainKey sdk.StoreKey, accountKey sdk.StoreKey, indexStoreKey sdk.StoreKey, configKeeper ConfigurationKeeper, supply SupplyKeeper, paramspace ParamSubspace) Keeper {
	keeper := Keeper{
		domainStoreKey:      domainKey,
		accountStoreKey:     accountKey,
		indexStoreKey:       indexStoreKey,
		cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		SupplyKeeper:        supply,
		paramspace:          paramspace,
	}
	return keeper
}

// Logger returns aliceAddr module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
