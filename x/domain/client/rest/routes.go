package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/tutils"
	"io/ioutil"
	"net/http"
)

// txRouteList clubs together all the transaction routes, which are the transactions
// // that return the bytes to sign to send a request that modifies state to the domain module
var txRoutesList = map[string]func(cliContext context.CLIContext) http.HandlerFunc{
	"registerDomain":            registerDomainHandler,
	"addAccountCertificates":    addAccountCertificatesHandler,
	"delAccountCertificates":    delAccountCertificateHandler,
	"deleteAccount":             deleteAccountHandler,
	"deleteDomain":              deleteDomainHandler,
	"flushDomain":               flushDomainHandler,
	"registerAccount":           registerAccountHandler,
	"renewAccount":              renewAccountHandler,
	"renewDomain":               renewDomainHandler,
	"replaceAccountTargets":     replaceAccountTargetsHandler,
	"transferAccountHandler":    transferAccountHandler,
	"transferDomainHandler":     transferDomainHandler,
	"setAccountMetadataHandler": setAccountMetadataHandler,
}

// registerTxRoutes registers all the transaction routes to the router
// the route will be exposed to storeName/handler, the handler will
// accept only post request with json codec
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	for route, handler := range txRoutesList {
		path := fmt.Sprintf("%s/tx/%s", storeName, route)
		r.HandleFunc(path, handler(cliCtx))
	}
}

func queryHandlerBuild(cliCtx context.CLIContext, storeName string, queryType iovns.QueryHandler) http.HandlerFunc {
	// get query type
	typ := tutils.GetPtrType(queryType)
	// return function
	return func(writer http.ResponseWriter, request *http.Request) {
		// clone queryType so we can unmarshal data to it
		query := tutils.CloneFromType(typ).(iovns.QueryHandler)
		// read request bytes
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}
		// unmarshal request from the client to the query handler
		err = iovns.DefaultQueryDecode(b, query)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		// verify query correctness
		if err = query.Validate(); err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
		}
		// marshal request to bytes understandable to the app query processor
		requestBytes, err := iovns.DefaultQueryEncode(query)
		if err != nil {
			// this is an internal server error if we're not able to marshal a request TODO log
			rest.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}
		// build query path
		queryPath := fmt.Sprintf("custom/%s/%s", storeName, query.QueryPath())
		// do query
		res, _, err := cliCtx.QueryWithData(queryPath, requestBytes)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		// success
		rest.PostProcessResponse(writer, cliCtx, res)
	}
}

// registerQueryRoutes registers all the routes used to query
// the domain module's keeper
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string, queries []iovns.QueryHandler) {
	for _, query := range queries {
		path := fmt.Sprintf("%s/query/%s", storeName, query.QueryPath())
		r.HandleFunc(path, queryHandlerBuild(cliCtx, storeName, query))
	}
}

// RegisterRoutes clubs together the tx and query routes
func RegisterRoutes(cliContext context.CLIContext, r *mux.Router, storeName string, queries []iovns.QueryHandler) {
	// register tx routes
	registerTxRoutes(cliContext, r, storeName)
	// register query routes
	registerQueryRoutes(cliContext, r, storeName, queries)
}
