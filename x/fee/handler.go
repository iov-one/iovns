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
		case types.MsgUpdateFees:
			return handleUpdateFees(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown request")
		}
	}
}

func handleUpdateFees(ctx sdk.Context, msg types.MsgUpdateFees, k keeper.Keeper) (*sdk.Result, error) {
	configurer := k.ConfigurationKeeper.GetConfiguration(ctx).Configurer
	if !configurer.Equals(msg.Configurer) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to update fees", msg.Configurer)
	}
	k.SetFees(ctx, msg.Fees)
	return &sdk.Result{}, nil
}
