package types

// QueryConfigResponse is the result returned after a query to the chain configuration
type QueryConfigResponse struct {
	Configuration Config `json:"configuration"`
}
