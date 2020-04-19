package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

// txRouteList clubs together all the transaction routes, which are the transactions
// // that return the bytes to sign to send a request that modifies state to the domain module
var txRoutesList = map[string]func(cliContext context.CLIContext) http.HandlerFunc{
	"registerDomain":         registerDomainHandler,
	"addAccountCertificates": addAccountCertificatesHandler,
	"delAccountCertificates": delAccountCertificateHandler,
	"deleteAccount":          deleteAccountHandler,
	"deleteDomain":           deleteDomainHandler,
	"flushDomain":            flushDomainHandler,
	"registerAccount":        registerAccountHandler,
	"renewAccount":           renewAccountHandler,
	"renewDomain":            renewDomainHandler,
	"replaceAccountTargets":  replaceAccountTargetsHandler,
	"transferAccountHandler": transferAccountHandler,
	"transferDomainHandler":  transferDomainHandler,
}

// registerTxRoutes registers all the transaction routes to the router
// the route will be exposed to storeName/handler, the handler will
// accept only post request with json codec
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	for route, handler := range txRoutesList {
		path := fmt.Sprintf("%s/%s", storeName, route)
		r.HandleFunc(path, handler(cliCtx))
	}
}

type queryHandler interface {
	Validate() error
	Route() string
	UnmarshalFromRest(r io.ReadCloser) error
	MarshalForApp() ([]byte, error)
}

func queryHandlerBuild(cliCtx context.CLIContext, storeName string, q queryHandler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// unmarshal request from the client to the query handler
		err := q.UnmarshalFromRest(request.Body)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
		}
		// marshal request to bytes understandable to the app query processor
		requestBytes, err := q.MarshalForApp()
		if err != nil {
			// this is an internal server error if we're not able to marshal a request TODO log
			rest.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
		}
		// build query path
		queryPath := fmt.Sprintf("custom/%s/%s", storeName, q.Route())
		// do query
		res, _, err := cliCtx.QueryWithData(queryPath, requestBytes)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
		}
		// success
		rest.PostProcessResponse(writer, cliCtx, res)
	}
}

// registerQueryRoutes registers all the routes used to query
// the domain module's keeper
func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {

}

func RegisterRoutes(cliContext context.CLIContext, r *mux.Router, storeName string) {
	// register tx routes
	registerTxRoutes(cliContext, r, storeName)
	// register query routes
	registerQueryRoutes(cliContext, r, storeName)
}
