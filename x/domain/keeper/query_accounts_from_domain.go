package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type QueryAccountsInDomain struct {
	Domain         string `json:"domain"`
	ResultsPerPage int    `json:"results_per_page"`
	Offset         int    `json:"offset"`
}

func (q *QueryAccountsInDomain) Handler() QueryHandlerFunc {
	return queryAccountsInDomainHandler
}

// Validate will validate the query model and set defaults
func (q *QueryAccountsInDomain) Validate() error {
	if q.Domain == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	// if results per page is unset then use default
	if q.ResultsPerPage <= 0 {
		q.ResultsPerPage = 100
	}
	// if offset is zero then use default
	if q.Offset <= 0 {
		q.Offset = 1
	}
	return nil
}

func (q *QueryAccountsInDomain) QueryPath() string {
	return "accountsInDomain"
}

type QueryAccountsInDomainResponse struct {
	Accounts []types.Account `json:"accounts"`
}

// queryAccountsInDomainHandler returns all accounts in aliceAddr domain
func queryAccountsInDomainHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	query := new(QueryAccountsInDomain)
	err := iovns.DefaultQueryDecode(req.Data, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// verify request
	if err = query.Validate(); err != nil {
		return nil, err
	}
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
	// iterate keys
	k.GetAccountsInDomain(ctx, query.Domain, do)
	// get accounts
	accounts := make([]types.Account, 0, len(keys))
	for _, key := range keys {
		account, _ := k.GetAccount(ctx, query.Domain, accountKeyToString(key))
		accounts = append(accounts, account)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryAccountsInDomainResponse{Accounts: accounts})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
