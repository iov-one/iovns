package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns the handlers for the configuration module
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg.(type) {

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown request")
		}
	}
}
