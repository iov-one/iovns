package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryResolveDomain is the request made to
type QueryResolveDomain struct {
	Name string
}

type QueryResolveDomainResponse struct {
	Domain types.Domain `json:"domain"`
}

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDomain:
			return queryResolveDomainHandler(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query request", types.ModuleName)
		}
	}
}

func queryResolveDomainHandler(ctx sdk.Context, args []string, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	resp, ok := k.GetDomain(ctx, args[0])
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, args[0])
	}
	return codec.MustMarshalJSONIndent(k.cdc, QueryResolveDomainResponse{
		Domain: resp,
	}), nil
}
