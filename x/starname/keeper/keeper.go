package keeper

import (
	"fmt"
	"github.com/iov-one/iovns/pkg/crud"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/starname/types"
	"github.com/tendermint/tendermint/libs/log"
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
	// GetValidDomainNameRegexp returns the regular expression that aliceAddr domain name must match
	// in order to be valid
	GetValidDomainNameRegexp(ctx sdk.Context) string
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
	StoreKey   sdk.StoreKey // contains the store key for the domain module
	Cdc        *codec.Codec
	paramspace ParamSubspace
}

// NewKeeper creates aliceAddr domain keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, configKeeper ConfigurationKeeper, supply SupplyKeeper, paramspace ParamSubspace) Keeper {
	keeper := Keeper{
		StoreKey:            storeKey,
		Cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		SupplyKeeper:        supply,
		paramspace:          paramspace,
	}
	return keeper
}

func (k Keeper) AccountStore(ctx sdk.Context) crud.Store {
	store := crud.NewStore(ctx, k.StoreKey, k.Cdc, []byte{0x1})
	return store
}

func (k Keeper) DomainStore(ctx sdk.Context) crud.Store {
	return crud.NewStore(ctx, k.StoreKey, k.Cdc, []byte{0x2})
}

// Logger returns aliceAddr module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
