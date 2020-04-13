package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/configuration/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier generates the queries handler for the configuration module
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryConfig:
			return queryConfig(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query for module: %s", types.ModuleName)
		}
	}
}

// queryConfig returns the configuration
func queryConfig(ctx sdk.Context, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	config := k.GetConfiguration(ctx)
	return ModuleCdc.MustMarshalJSON(types.QueryConfigResponse{Configuration: config}), nil
}
