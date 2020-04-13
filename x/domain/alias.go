package domain

import (
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// aliasing for naming constants
const (
	ModuleName      = types.ModuleName
	DomainStoreKey  = types.DomainStoreKey
	AccountStoreKey = types.AccountStoreKey
	QuerierRoute    = types.QuerierRoute
	RouterKey       = types.RouterKey
)

// aliasing for types
type (
	Keeper = keeper.Keeper
)

// aliasing for funcs
var NewKeeper = keeper.NewKeeper
