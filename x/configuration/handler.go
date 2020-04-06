package configuration

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewHandler returns the handlers for the configuration
// since configuration has no active handlers, it does nothing
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		return &sdk.Result{
			Data:   nil,
			Log:    "",
			Events: nil,
		}, nil
	}
}
