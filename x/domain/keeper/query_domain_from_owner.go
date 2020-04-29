package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryDomainsFromOwner is the request model used
// to query domains owned by a sdk.AccAddress
type QueryDomainsFromOwner struct {
	// Owner is the address of the owner of the domains
	Owner sdk.AccAddress `json:"owner"`
	// ResultsPerPage is the number of results displayed in a page
	ResultsPerPage int `json:"results_per_page"`
	// Offset is the page number
	Offset int `json:"offset"`
}

// Use is a placeholder
func (q *QueryDomainsFromOwner) Use() string {
	return "owner-domains"
}

// Description is a placeholder
func (q *QueryDomainsFromOwner) Description() string {
	return "gets all the domains owned by the given address"
}

// Handler implements the local queryHandler
func (q *QueryDomainsFromOwner) Handler() QueryHandlerFunc {
	return queryDomainsFromOwnerHandler
}

// QueryPath implements iovns.QueryHandler
func (q *QueryDomainsFromOwner) QueryPath() string {
	return "domainsFromOwner"
}

// Validate implements iovns.QueryHandler
func (q *QueryDomainsFromOwner) Validate() error {
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

// QueryDomainsFromOwnerResponse is the response
// returned by the QueryDomainsFromOwner query
type QueryDomainsFromOwnerResponse struct {
	// Domains is a slice of the domains
	// found by the query
	Domains []types.Domain
}

// queryDomainsFromOwnerHandler is the query handler used to get all the domains owned by an sdk.AccAddress
func queryDomainsFromOwnerHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	query := new(QueryDomainsFromOwner)
	err := iovns.DefaultQueryDecode(req.Data, query)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	// validate request
	if err = query.Validate(); err != nil {
		return nil, err
	}
	// get domain keys
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
	// fill domain keys
	k.iterDomainToOwner(ctx, query.Owner, do)
	// get domains
	domains := make([]types.Domain, 0, len(keys))
	for _, key := range keys {
		_, domainName := splitOwnerToDomainKey(key)
		domain, _ := k.GetDomain(ctx, domainName)
		domains = append(domains, domain)
	}
	respBytes, err := iovns.DefaultQueryEncode(QueryDomainsFromOwnerResponse{Domains: domains})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
