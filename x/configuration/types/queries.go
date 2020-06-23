package types

import (
	"github.com/iov-one/iovns/x/fee/types"
)

// QueryConfigResponse is the result returned after a query to the chain configuration
type QueryConfigResponse struct {
	Configuration Config `json:"configuration"`
}

type QueryFeesResponse struct {
	Fees *types.FeeConfiguration `json:"fees"`
}
