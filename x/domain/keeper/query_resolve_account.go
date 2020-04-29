package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryResolveAccount is the query model
// used to resolve an account
type QueryResolveAccount struct {
	// Domain is the name of the domain of the account
	Domain string `json:"domain"`
	// Name is the name of the account
	Name string `json:"name"`
}

// Use is a placeholder
func (q *QueryResolveAccount) Use() string {
	return "resolve-account"
}

// Description is a placeholder
func (q *QueryResolveAccount) Description() string {
	return "resolves the given account"
}

// Handler implements local queryHandler
func (q *QueryResolveAccount) Handler() QueryHandlerFunc {
	return queryResolveAccountHandler
}

// Validate implements iovns.QueryHandler
func (q *QueryResolveAccount) Validate() error {
	if q.Domain == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	if q.Name == "" {
		return sdkerrors.Wrapf(types.ErrInvalidAccountName, "empty")
	}
	return nil
}

// QueryPath implements iovns.QueryHandler
func (q *QueryResolveAccount) QueryPath() string {
	return "resolveAccount"
}

// QueryResolveAccountResponse is the response
// returned by the QueryResolveAccount query
type QueryResolveAccountResponse struct {
	// Account contains the resolved account
	Account types.Account `json:"account"`
}

// queryResolveAccountHandler is the query handler that takes care of resolving accounts
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
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", q.Domain, q.Name)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryResolveAccountResponse{Account: account})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
