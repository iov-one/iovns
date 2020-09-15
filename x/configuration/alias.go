package configuration

import "github.com/iov-one/iovns/x/configuration/types"

// alias for types
type (
	// Config aliases types.Config
	Config = types.Config
	// Fees aliases types.Fees
	Fees = types.Fees
)

// alias for consts
const (
	ModuleName        = types.ModuleName   // ModuleName aliases types.ModuleName
	RouterKey         = types.RouterKey    // RouterKey aliases types.RouterKey
	QuerierRoute      = types.QuerierRoute // QuerierRoute aliases types.QuerierRoute
	QueryConfig       = types.QueryConfig  // QueryConfig aliases types.QueryConfig
	StoreKey          = types.StoreKey     // StoreKey aliases types.StoreKey
	DefaultParamSpace = types.DefaultParamSpace
)

// function aliases

var (
	NewFees       = types.NewFees // NewFees aliases types.NewFees
	RegisterCodec = types.RegisterCodec
)
