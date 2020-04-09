package domain

import (
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/iov-one/iovnsd/x/domain/types"
)

// aliasing for naming constants
const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
)

// aliasing for types
type (
	Keeper = keeper.Keeper
)

// aliasing for funcs
var NewKeeper = keeper.NewKeeper
