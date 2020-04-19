package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type QueryDomainsFromOwner struct {
	Owner          sdk.AccAddress `json:"owner"`
	ResultsPerPage int            `json:"results_per_page"`
	Offset         int            `json:"offset"`
}

func (q *QueryDomainsFromOwner) QueryPath() string {
	return "domainsFromOwner"
}

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

type QueryDomainsFromOwnerResponse struct {
	Domains []types.Domain
}

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
	domainKeys := k.iterDomainToOwner(ctx, query.Owner)
	nKeys := len(domainKeys) // total number of keys
	// no results
	if nKeys == 0 {
		// return response
		respBytes, err := iovns.DefaultQueryEncode(QueryDomainsFromOwnerResponse{})
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
	domains := make([]types.Domain, 0, lastObjectIndex-firstObjectIndex+1)
	// fill accounts
	for currIndex := firstObjectIndex; currIndex <= lastObjectIndex; currIndex++ {
		// get domainName
		_, domainName := splitOwnerToDomainKey(domainKeys[currIndex])
		domain, _ := k.GetDomain(ctx, domainName)
		// append
		domains = append(domains, domain)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryDomainsFromOwnerResponse{Domains: domains})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
