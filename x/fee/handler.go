package fee

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/fee/keeper"
	"github.com/iov-one/iovns/x/fee/types"
)

// NewHandler returns the handlers for the configuration module
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgUpdateFeeConfiguration:
			return handleUpdateConfiguration(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown request")
		}
	}
}

func handleUpdateConfiguration(ctx sdk.Context, msg types.MsgUpdateFeeConfiguration, k keeper.Keeper) (*sdk.Result, error) {
	configurer := k.GetFeeConfigurer(ctx)
	if !configurer.Equals(msg.Signer) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to update fees", msg.Signer)
	}
	k.SetFeeConfigurer(ctx, msg.NewFeeConfiguration.FeeConfigurer)
	k.SetFeeParameters(ctx, msg.NewFeeConfiguration.FeeParameters)
	for _, fs := range msg.NewFeeConfiguration.FeeSeeds {
		k.SetFeeSeed(ctx, fs)
	}
	return &sdk.Result{}, nil
}
