package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovnsd/x/account/types"
)

// NewQuerier creates a new querier for account clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		// TODOCHANGE
		case "TOODCHANGE":
			return queryParams(ctx, k)
			// TODO: Put the modules query routes
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown account query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	// params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, nil)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// TODO: Add the modules query functions
// They will be similar to the above one: queryParams()
