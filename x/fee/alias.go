package fee

import (
	"github.com/iov-one/iovns/x/fee/keeper"
	"github.com/iov-one/iovns/x/fee/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamSpace = types.DefaultParamSpace
	StoreKey          = types.StoreKey // StoreKey aliases types.StoreKey
)

// aliasing for types
type (
	Keeper = keeper.Keeper
)

// aliasing for funcs
var NewKeeper = keeper.NewKeeper
