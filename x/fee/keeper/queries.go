package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/fee/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryHandlerFunc defines the query handler for this module
type QueryHandlerFunc func(ctx sdk.Context, path []string, query abci.RequestQuery, k Keeper) ([]byte, error)

// AvailableQueries returns the list of available queries in the module
func AvailableQueries() []iovns.QueryHandler {
	queries := []iovns.QueryHandler{
		&QueryFees{},
	}
	return queries
}

// queryRouter defines a router for domain queries
type queryRouter map[string]QueryHandlerFunc

func buildRouter(queries []iovns.QueryHandler) queryRouter {
	// queryHandler extends the default query handler
	// to provide also an handler function required to
	// build a router
	type queryHandler interface {
		iovns.QueryHandler
		Handler() QueryHandlerFunc
	}
	// build router
	router := make(queryRouter, len(queries))
	for _, query := range queries {
		queryAndHandler, ok := query.(queryHandler)
		// if interface is not implemented then the query type formation is invalid
		if !ok {
			panic(fmt.Sprintf("invalid query type: %T", query))
		}
		router[queryAndHandler.QueryPath()] = queryAndHandler.Handler()
	}
	// return
	return router
}

// NewQuerier builds the query handler for the module
func NewQuerier(k Keeper) sdk.Querier {
	// get queries
	queries := AvailableQueries()
	router := buildRouter(queries)
	// return sdk.Querier
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		handler, ok := router[path[0]]
		// handler not found, query does not exist
		if !ok {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%s", path[0])
		}
		// handler
		return handler(ctx, path, req, k)
	}
}

type QueryFees struct{}

func (q *QueryFees) Use() string {
	return "query-fees"
}

func (q *QueryFees) Description() string {
	return "return the current fees"
}

func (q *QueryFees) Handler() QueryHandlerFunc {
	return queryFeeConfigurationHandler
}

func (q *QueryFees) Validate() error {
	return nil
}

func (q *QueryFees) QueryPath() string {
	return "fees"
}

func queryFeeConfigurationHandler(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	cfg := k.GetFeeConfiguration(ctx)
	// return response
	respBytes, err := iovns.DefaultQueryEncode(QueryFeesResponse{Fees: cfg})
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return respBytes, nil
}

type QueryFeesResponse struct {
	Fees types.FeeConfiguration `json:"fees"`
}
