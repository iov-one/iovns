package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/iov-one/bnsd/x/account"
	"time"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

/*
When a module wishes to interact with another module, it is good practice to define what it will use
as an interface so the module cannot use things that are not permitted.
TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
*/

// ConfigurationKeeper defines the behaviour of the configuration state checks
type ConfigurationKeeper interface {
	// GetOwner returns the owner
	GetOwner(ctx sdk.Context) sdk.AccAddress
	// GetValidDomainRegexp returns the regular expression that a domain name must match
	// in order to be valid
	GetValidDomainRegexp(ctx sdk.Context) string
	// GetDomainRenewDuration returns the default duration of a domain renewal
	GetDomainRenewDuration(ctx sdk.Context) time.Duration
}

// AccountKeeper defines the behaviour of the account module required by the domain
// module to interact with it
type AccountKeeper interface {
	// SetAccount saves the account in the state
	SetAccount(ctx sdk.Context, account account.Account)
}
