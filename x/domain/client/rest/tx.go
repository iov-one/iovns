package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/iov-one/iovns/x/domain/types"
	"net/http"
)

// handleTxRequest is a helper function that takes care of checking base requests, sdk messages, after verifying
// requests it forwards an error to the client in case of error, otherwise it will return a transaction to sign
// and send to the /tx endpoint to do a request
func handleTxRequest(cliCtx context.CLIContext, baseReq rest.BaseReq, msg sdk.Msg, writer http.ResponseWriter) {
	baseReq = baseReq.Sanitize()
	if !baseReq.ValidateBasic(writer) {
		return
	}
	// validate request
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
	}
	// write tx
	utils.WriteGenerateStdTxResponse(writer, cliCtx, baseReq, []sdk.Msg{msg})
}

// registerDomain defines the request model used for registerDomainHandler
type registerDomain struct {
	BaseReq rest.BaseReq             `json:"base_req"`
	Message *types.MsgRegisterDomain `json:"message"`
}

// registerDomainHandler builds the transaction to sign
func registerDomainHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req registerDomain
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, "failed to parse request")
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// addAccountCertificates is the request model for addAccountCertificatesHandler
type addAccountCertificates struct {
	BaseReq rest.BaseReq                     `json:"base_req"`
	Message *types.MsgAddAccountCertificates `json:"message"`
}

// addAccountCertificatesHandler builds the transaction to sign to add account certificates
func addAccountCertificatesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req addAccountCertificates
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// delAccountCertificate is the request model for delAccountCertificateHandler
type delAccountCertificate struct {
	BaseReq rest.BaseReq                       `json:"base_req"`
	Message *types.MsgDeleteAccountCertificate `json:"message"`
}

// delAccountCertificateHandler builds the transaction to sign to delete account certificates
func delAccountCertificateHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req delAccountCertificate
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// deleteAccount is the request
type deleteAccount struct {
	BaseReq rest.BaseReq            `json:"base_req"`
	Message *types.MsgDeleteAccount `json:"message"`
}

// deleteAccountHandler builds the transaction to sign to delete an account
func deleteAccountHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req deleteAccount
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// deleteDomain is the request model for deleteDomainHandler
type deleteDomain struct {
	BaseReq rest.BaseReq           `json:"base_req"`
	Message *types.MsgDeleteDomain `json:"message"`
}

// deleteDomainHandler builds the transaction to sign to delete a domain
func deleteDomainHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req deleteDomain
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// flushDomain is the request model for flushDomainHandler
type flushDomain struct {
	BaseReq rest.BaseReq          `json:"base_req"`
	Message *types.MsgFlushDomain `json:"message"`
}

// flushDomainHandler builds the transaction to sign to flush a domain
func flushDomainHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req flushDomain
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// registerAccount is the request model used for registerAccountHandler
type registerAccount struct {
	BaseReq rest.BaseReq              `json:"base_req"`
	Message *types.MsgRegisterAccount `json:"message"`
}

// registerAccountHandler builds the transaction to sign to register an account
func registerAccountHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req registerAccount
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// renewAccount is the request model for renewAccountHandler
type renewAccount struct {
	BaseReq rest.BaseReq           `json:"base_req"`
	Message *types.MsgRenewAccount `json:"message"`
}

// renewAccountHandler builds the transaction request to sign to renew a domain
func renewAccountHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req renewAccount
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// renewDomain is the request model for renewDomainHandler
type renewDomain struct {
	BaseReq rest.BaseReq          `json:"base_req"`
	Message *types.MsgRenewDomain `json:"message"`
}

// renewDomainHandler builds the transaction to sign to renew a domain
func renewDomainHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req renewDomain
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// replaceAccountTargets is the request model for replaceAccountTargets
type replaceAccountTargets struct {
	BaseReq rest.BaseReq                    `json:"base_req"`
	Message *types.MsgReplaceAccountTargets `json:"message"`
}

// replaceAccountTargets builds the transaction to sign to replace account targets
func replaceAccountTargetsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req renewDomain
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// transferAccount is the request model for transferAccountHandler
type transferAccount struct {
	BaseReq rest.BaseReq              `json:"base_req"`
	Message *types.MsgTransferAccount `json:"message"`
}

// transferAccountHandler builds the transaction to sign to transfer accounts
func transferAccountHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req transferAccount
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

// transferDomain is the request model for transferDomainHandler
type transferDomain struct {
	BaseReq rest.BaseReq             `json:"base_req"`
	Message *types.MsgTransferDomain `json:"message"`
}

// transferDomainHandler builds the transaction to sign to transfer domains
func transferDomainHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req transferDomain
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}
