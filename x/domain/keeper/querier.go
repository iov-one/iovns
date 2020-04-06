package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovnsd/x/domain/types"
)

// the list of query endpoints supported
const (
	QueryDomain = "domain"
)

// NewQuerier creates a new querier for domain clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryDomain:
			return queryDomain(ctx, path[1:], req, k)
		/*
			case types.QueryParams:
				return queryParams(ctx, k)
				// TODO: Put the modules query routes
		*/
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown domain query endpoint")
		}
	}
}

/*
func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}
*/

// TODO: Add the modules query functions
// queryDomain returns the domain
func queryDomain(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper Keeper) ([]byte, error) {
	// get domain
	domain, ok := keeper.GetDomain(ctx, path[0])
	// check if it exists
	if !ok {
		return nil, sdkerrors.Wrap(types.ErrDomainDoesNotExist, path[0])
	}
	// return response
	return codec.MustMarshalJSONIndent(keeper.cdc, domain), nil
}

// They will be similar to the above one: queryParams()
