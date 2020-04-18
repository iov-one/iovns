package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryDomainRequest is the request made to
type QueryDomainRequest struct {
	Name string
}

type QueryDomainResponse struct {
	Domain types.Domain `json:"domain"`
}

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDomain:
			return queryGet(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query request", types.ModuleName)
		}
	}
}

func queryGet(ctx sdk.Context, args []string, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	resp, ok := k.GetDomain(ctx, args[0])
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, args[0])
	}
	return codec.MustMarshalJSONIndent(k.cdc, QueryDomainResponse{
		Domain: resp,
	}), nil
}

type QueryAccountsInDomain struct {
	Domain         string `json:"domain"`
	ResultsPerPage int    `json:"results_per_page"`
	Offset         int    `json:"offset"`
}

type QueryAccountsInDomainResponse struct {
}

// queryGetAccountsInDomain returns all accounts in a domain
func queryGetAccountsInDomain(ctx sdk.Context, args []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	panic("implement")
}
