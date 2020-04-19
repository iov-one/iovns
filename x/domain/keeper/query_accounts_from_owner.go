package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type QueryAccountsFromOwner struct {
	Owner          sdk.AccAddress `json:"owner"`
	ResultsPerPage int            `json:"results_per_page"`
	Offset         int            `json:"offset"`
}

func (q *QueryAccountsFromOwner) Handler() QueryHandlerFunc {
	return queryAccountsFromOwnerHandler
}

func (q *QueryAccountsFromOwner) QueryPath() string {
	return "accountsFromOwner"
}

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

type QueryAccountsFromOwnerResponse struct {
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
	// get keys from owner
	accKeys := k.iterAccountToOwner(ctx, query.Owner)
	nKeys := len(accKeys) // total number of keys
	// no results
	if nKeys == 0 {
		// return response
		respBytes, err := iovns.DefaultQueryEncode(QueryAccountsFromOwnerResponse{})
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return respBytes, nil
	}
	// get the index of the first object we want
	firstObjectIndex := query.Offset*query.ResultsPerPage - query.ResultsPerPage
	// check if there is at least one object at that index
	if nKeys < firstObjectIndex+1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid offset")
	}

	// get the index for the last object
	lastObjectIndex := firstObjectIndex + query.ResultsPerPage - 1
	// check if last object index would outbound our acc keys slice
	if lastObjectIndex > nKeys {
		lastObjectIndex = nKeys - 1 // if it does then set last index as the last element of our slice
	}
	accounts := make([]types.Account, 0, lastObjectIndex-firstObjectIndex+1)
	// fill accounts
	for currIndex := firstObjectIndex; currIndex <= lastObjectIndex; currIndex++ {
		_, domain, accountName := splitOwnerToAccountKey(accKeys[currIndex])
		// get account
		account, _ := k.GetAccount(ctx, domain, accountName)
		// append
		accounts = append(accounts, account)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryAccountsFromOwnerResponse{Accounts: accounts})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
