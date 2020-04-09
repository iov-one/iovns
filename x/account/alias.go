package account

import (
	"github.com/iov-one/iovnsd/x/account/types"
)

const (
	// TODO: define constants that you would like exposed from the internal package

	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	DefaultParamspace = types.DefaultParamspace
	QueryParams       = types.QueryParams
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases
	RegisterCodec = types.RegisterCodec
	// TODO: Fill out function aliases

	// variable aliases
	ModuleCdc = types.ModuleCdc
	// TODO: Fill out variable aliases
)

type (
	Params = types.Params

	// TODO: Fill out module types
	Account = types.Account
)
