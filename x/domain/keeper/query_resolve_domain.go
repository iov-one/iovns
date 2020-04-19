package keeper

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryResolveDomain is the request made to
type QueryResolveDomain struct {
	Name string `json:"name"`
}

func (q *QueryResolveDomain) Validate() error {
	if q.Name == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	return nil
}

type QueryResolveDomainResponse struct {
	Domain types.Domain `json:"domain"`
}

func queryResolveDomainHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	q := new(QueryResolveDomain)
	err := json.Unmarshal(req.Data, q)
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
	return codec.MustMarshalJSONIndent(k.cdc, QueryResolveDomainResponse{
		Domain: domain,
	}), nil
}
