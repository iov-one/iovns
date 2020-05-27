package domain

import (
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// aliasing for naming constants
const (
	// ModuleNames aliases types.ModuleName
	ModuleName = types.ModuleName
	// DomainStoreKey aliases types.DomainStoreKey
	DomainStoreKey = types.DomainStoreKey
	// AccountStoreKey aliases types.AccountStoreKey
	AccountStoreKey = types.AccountStoreKey
	// IndexStoreKey aliases types.IndexStoreKey
	IndexStoreKey = types.IndexStoreKey
	// QuerierRoute aliases types.QuerierRoute
	QuerierRoute = types.QuerierRoute
	// RouterKey aliases types.RouterKey
	RouterKey = types.RouterKey
	// DefaultParamSpace defines domain module default param space key
	DefaultParamSpace = types.DefaultParamSpace
)

// aliasing for types
type (
	Keeper = keeper.Keeper
)

// aliasing for funcs
var NewKeeper = keeper.NewKeeper
