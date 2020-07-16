package starname

import (
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

// aliasing for naming constants
const (
	// ModuleNames aliases types.ModuleName
	ModuleName = types.ModuleName
	// DomainStoreKey aliases types.DomainStoreKey
	DomainStoreKey = types.DomainStoreKey
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
var (
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
)
