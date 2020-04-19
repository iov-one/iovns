package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type QueryResolveAccount struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
}

func (q *QueryResolveAccount) Validate() error {
	if q.Domain == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	if q.Name == "" {
		return sdkerrors.Wrapf(types.ErrInvalidAccountName, "empty")
	}
	return nil
}

func (q *QueryResolveAccount) QueryPath() string {
	return "resolveAccount"
}

type QueryResolveAccountResponse struct {
	Account types.Account `json:"account"`
}

func queryResolveAccountHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	q := new(QueryResolveAccount)
	err := iovns.DefaultQueryDecode(req.Data, q)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// validate
	if err = q.Validate(); err != nil {
		return nil, err
	}
	// do query
	account, exists := k.GetAccount(ctx, q.Domain, q.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found: account %s in domain %s", q.Name, q.Domain)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryResolveAccountResponse{Account: account})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
