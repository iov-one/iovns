package domain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		// domain handlers
		case types.MsgRegisterDomain:
			return handleMsgRegisterDomain(ctx, k, msg)
		case types.MsgRenewDomain:
			return handlerMsgRenewDomain(ctx, k, msg)
		case types.MsgDeleteDomain:
			return handlerMsgDeleteDomain(ctx, k, msg)
		case types.MsgFlushDomain:
			return handlerMsgFlushDomain(ctx, k, msg)
		// account handlers
		case types.MsgRegisterAccount:
			return handleMsgRegisterAccount(ctx, k, msg)
		case types.MsgRenewAccount:
			return handlerMsgRenewAccount(ctx, k, msg)
		case types.MsgAddAccountCertificates:
			return handlerMsgAddAccountCertificates(ctx, k, msg)
		case types.MsgDeleteAccountCertificate:
			return handlerMsgDeleteAccountCertificate(ctx, k, msg)
		case types.MsgDeleteAccount:
			return handlerMsgDeleteAccount(ctx, k, msg)
		case types.MsgReplaceAccountTargets:
			return handlerMsgReplaceAccountTargets(ctx, k, msg)
		case types.MsgTransferAccount:
			return handlerMsgTransferAccount(ctx, k, msg)

		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unregonized request: %T", msg))
		}
	}
}
