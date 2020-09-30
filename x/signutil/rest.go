package signutil

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
)

type Success struct {
	Message  string `json:"message"`
	Signer   string `json:"signer"`
	Verified bool   `json:"verified"`
	Signed   string `json:"signed"`
}

type Error struct { // TODO: use the sdk's REST utils
	Code   uint64 `json:"code"`
	Reason string `json:"error"`
}

func RegisterRestRoutes(ctx context.CLIContext, r *mux.Router) {
	cdc := ctx.Codec
	r.HandleFunc(fmt.Sprintf("/%s/query/verify", ModuleName), func(writer http.ResponseWriter, request *http.Request) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{ // TODO: log error on server
				Code:   http.StatusInternalServerError,
				Reason: err.Error(),
			}))
			return
		}
		var req auth.StdTx
		err = cdc.UnmarshalJSON(b, &req)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{ // TODO: log error on server
				Code:   http.StatusBadRequest,
				Reason: err.Error(),
			}))
			return
		}

		err = Verify(req, DefaultChainID, DefaultAccountNumber, DefaultSequence)
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{ // TODO: log error on server
				Code:   http.StatusUnauthorized,
				Reason: fmt.Sprintf("Did you sign with --chain-id '%s', --account-number %d, and --sequence %d?", DefaultChainID, DefaultAccountNumber, DefaultSequence),
			}))
			return
		}

		msgs := req.GetMsgs()
		if len(msgs) != 1 {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{ // TODO: log error on server
				Code:   http.StatusBadRequest,
				Reason: fmt.Sprintf("Expected 1 msg but got %d.", len(msgs)),
			}))
			return
		}

		var msg MsgSignText
		err = cdc.UnmarshalJSON(msgs[0].GetSignBytes(), &msg)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{ // TODO: log error on server
				Code:   http.StatusBadRequest,
				Reason: err.Error(),
			}))
			return
		}

		// success
		writer.WriteHeader(http.StatusOK) // TODO: unify response format with other successful REST responses in the starname module
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(cdc.MustMarshalJSON(Success{ // TODO: log success on server
			Message:  msg.Message,
			Signer:   msg.Signer.String(),
			Verified: true,
			Signed:   string(b[:]),
		}))
	})
}
