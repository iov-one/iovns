package configuration

import (
	"github.com/iov-one/iovns/x/configuration/types"
	feeTypes "github.com/iov-one/iovns/x/fee/types"
)

// alias for types
type (
	Config = types.Config // Config aliases types.Config
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
	NewFees = feeTypes.NewFeeConfiguration // NewFeeConfiguration aliases types.NewFeeConfiguration
)
