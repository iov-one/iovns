package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryResolveDomain is the request model
// used to resolve a domain
type QueryResolveDomain struct {
	// Name is the domain name
	Name string `json:"name" arg:"positional"`
}

// Use is a placeholder
func (q *QueryResolveDomain) Use() string {
	return "resolve-domain"
}

// Description is a placeholder
func (q *QueryResolveDomain) Description() string {
	return "resolves a domain"
}

// Handler implements the local queryHandler
func (q *QueryResolveDomain) Handler() QueryHandlerFunc {
	return queryResolveDomainHandler
}

// QueryPath implements iovns.QueryHandler
func (q *QueryResolveDomain) QueryPath() string {
	return "resolveDomain"
}

// Validate implements iovns.QueryHandler
func (q *QueryResolveDomain) Validate() error {
	if q.Name == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	return nil
}

// QueryResolveDomainResponse is response returned
// by the QueryResolveDomain query
type QueryResolveDomainResponse struct {
	// Domain contains the queried domain information
	Domain types.Domain `json:"domain"`
}

// queryResolveDomainHandler takes care of resolving domains
func queryResolveDomainHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	q := new(QueryResolveDomain)
	err := iovns.DefaultQueryDecode(req.Data, q)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if err = q.Validate(); err != nil {
		return nil, err
	}
	domain, ok := k.GetDomain(ctx, q.Name)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", q.Name)
	}
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryResolveDomainResponse{Domain: domain})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}
