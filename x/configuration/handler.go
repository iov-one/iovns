package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/configuration/types"
)

// NewHandler returns the handlers for the configuration module
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgUpdateConfig:
			return handleUpdateConfig(ctx, msg, k)
		case types.MsgUpdateFees:
			return handleUpdateFees(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown request")
		}
	}
}

func handleUpdateFees(ctx sdk.Context, msg types.MsgUpdateFees, k Keeper) (*sdk.Result, error) {
	configurer := k.GetConfigurer(ctx)
	if !configurer.Equals(msg.Configurer) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to update fees", msg.Configurer)
	}
	k.SetFees(ctx, msg.Fees)
	// TODO emit event
	return &sdk.Result{}, nil
}

func handleUpdateConfig(ctx sdk.Context, msg types.MsgUpdateConfig, k Keeper) (*sdk.Result, error) {
	configurer := k.GetConfigurer(ctx)
	if !configurer.Equals(msg.Signer) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to update configuration", msg.Signer)
	}
	// if allowed update configuration
	k.SetConfig(ctx, msg.NewConfiguration)
	// TODO emit event
	return &sdk.Result{}, nil
}
