package signutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type VerifyRequest struct {
	ChainID       string     `json:"chain_id"`
	AccountNumber uint64     `json:"account_number"`
	Sequence      uint64     `json:"sequence"`
	StdTx         auth.StdTx `json:"std_tx"`
}

type Error struct {
	Code   uint64 `json:"code"`
	Reason string `json:"error"`
}

func RegisterRestRoutes(ctx context.CLIContext, r *mux.Router) {
	cdc := ctx.Codec
	r.HandleFunc(fmt.Sprintf("%s/query/verify", ModuleName), func(writer http.ResponseWriter, request *http.Request) {
		b, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{
				Code:   http.StatusInternalServerError,
				Reason: err.Error(),
			}))
			return
		}
		var req VerifyRequest
		err = cdc.UnmarshalJSON(b, &req)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{
				Code:   http.StatusBadRequest,
				Reason: err.Error(),
			}))
			return
		}

		if req.ChainID == "" {
			req.ChainID = DefaultChainID
		}
		if req.AccountNumber == 0 {
			req.AccountNumber = DefaultAccountNumber
		}
		if req.Sequence == 0 {
			req.Sequence = DefaultSequence
		}

		err = Verify(req.StdTx, req.ChainID, req.AccountNumber, req.Sequence)
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = writer.Write(cdc.MustMarshalJSON(Error{
				Code:   http.StatusUnauthorized,
				Reason: err.Error(),
			}))
			return
		}
	})
}
