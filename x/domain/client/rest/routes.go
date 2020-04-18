package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
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
