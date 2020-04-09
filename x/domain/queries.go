package domain

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

const QueryDomainPath = "get"

// QueryDomainRequest is the request made to
type QueryDomainRequest struct {
	Name string
}

type QueryDomainResponse struct {
	Domain Domain `json:"domain"`
}

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryDomainPath:
			return queryGet(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query request", ModuleName)
		}
	}
}

func queryGet(ctx sdk.Context, args []string, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	resp, ok := k.GetDomain(ctx, args[0])
	if !ok {
		return nil, sdkerrors.Wrapf(ErrDomainDoesNotExist, args[0])
	}
	return codec.MustMarshalJSONIndent(k.cdc, QueryDomainResponse{
		Domain: resp,
	}), nil
}
