package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryAccountsFromOwner queries all the accounts
// owned by a certain sdk.AccAddress
type QueryAccountsFromOwner struct {
	// Owner is the owner of the accounts
	Owner sdk.AccAddress `json:"owner"`
	// ResultsPerPage is the number of results returned in each page
	ResultsPerPage int `json:"results_per_page"`
	// Offset is the page number
	Offset int `json:"offset"`
}

// Use is a placeholder
func (q *QueryAccountsFromOwner) Use() string {
	return "owner-accounts"
}

// Description is a placeholder
func (q *QueryAccountsFromOwner) Description() string {
	return "gets all the accounts owned by a given address"
}

// Handler implements local queryHandler
func (q *QueryAccountsFromOwner) Handler() QueryHandlerFunc {
	return queryAccountsFromOwnerHandler
}

// QueryPath implements iovns.QueryHandler
func (q *QueryAccountsFromOwner) QueryPath() string {
	return "accountsFromOwner"
}

// Validate implements iovns.QueryHandler
func (q *QueryAccountsFromOwner) Validate() error {
	if q.Owner == nil {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "empty")
	}
	if q.ResultsPerPage == 0 {
		q.ResultsPerPage = 100
	}
	if q.Offset == 0 {
		q.Offset = 1
	}
	return nil
}

// QueryAccountsFromOwnerResponse is the response model
// returned by QueryAccountsFromOwner
type QueryAccountsFromOwnerResponse struct {
	// Accounts is a slice containing the accounts
	// returned by the query
	Accounts []types.Account `json:"accounts"`
}

// queryAccountsFromOwnerHandler gets all the accounts related to an account address
func queryAccountsFromOwnerHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	query := new(QueryAccountsFromOwner)
	err := iovns.DefaultQueryDecode(req.Data, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// validate request
	if err = query.Validate(); err != nil {
		return nil, err
	}
	// generate expected keys
	keys := make([][]byte, 0, query.ResultsPerPage)
	index := 0
	// calculate index range
	indexStart := query.ResultsPerPage*query.Offset - query.ResultsPerPage // this is the start
	indexEnd := indexStart + query.ResultsPerPage - 1                      // this is the end
	do := func(key []byte) bool {
		// check if our index is grater-equal than our start
		if index >= indexStart {
			keys = append(keys, key)
		}
		if index == indexEnd {
			return false
		}
		// increase index
		index++
		return true
	}
	// iterate account keys
	k.iterAccountToOwner(ctx, query.Owner, do)
	// check if there are any keys
	if len(keys) == 0 {
		respBytes, err := iovns.DefaultQueryEncode(QueryAccountsFromOwnerResponse{})
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return respBytes, nil
	}
	// fill accounts
	accounts := make([]types.Account, 0, len(keys))
	for _, accKey := range keys {
		_, domainName, accountName := splitOwnerToAccountKey(accKey)
		account, _ := k.GetAccount(ctx, domainName, accountName)
		accounts = append(accounts, account)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryAccountsFromOwnerResponse{Accounts: accounts})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
