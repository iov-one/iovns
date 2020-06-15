package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	. "github.com/iov-one/iovns/x/configuration/types"
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

type updateConfig struct {
	BaseReq rest.BaseReq     `json:"base_req"`
	Message *MsgUpdateConfig `json:"message"`
}

func updateConfigHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req updateConfig
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, "failed to parse request")
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}

type updateFees struct {
	BaseReq rest.BaseReq   `json:"base_req"`
	Message *MsgUpdateFees `json:"message"`
}

func updateFeesHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req updateFees
		if !rest.ReadRESTReq(writer, request, cliCtx.Codec, &req) {
			return
		}
		handleTxRequest(cliCtx, req.BaseReq, req.Message, writer)
	}
}
