package account

import (
	"fmt"
	"github.com/iov-one/iovnsd/x/account/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the account type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgRegisterDomain:
			return handleMsgRegisterDomain(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgRegisterDomain registers the domain
func handleMsgRegisterDomain(ctx sdk.Context, k Keeper, msg sdk.Msg) (*sdk.Result, error) {
	// insert rules
	// success
	return &sdk.Result{}, nil
}
