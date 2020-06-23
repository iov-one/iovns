package domain

import (
	"fmt"
	"log"

	fee "github.com/iov-one/iovns/x/fee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
)

// NewHandler builds the tx requests handler for the domain module
func NewHandler(k Keeper, feeCollector fee.Collector) sdk.Handler {
	f := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		// domain handlers
		case *types.MsgRegisterDomain:
			return handleMsgRegisterDomain(ctx, k, feeCollector, msg)
		case *types.MsgRenewDomain:
			return handlerMsgRenewDomain(ctx, k, feeCollector, msg)
		case *types.MsgDeleteDomain:
			return handlerMsgDeleteDomain(ctx, k, feeCollector, msg)
		case *types.MsgTransferDomain:
			return handlerMsgTransferDomain(ctx, k, feeCollector, msg)
		// account handlers
		case *types.MsgRegisterAccount:
			return handleMsgRegisterAccount(ctx, k, feeCollector, msg)
		case *types.MsgRenewAccount:
			return handlerMsgRenewAccount(ctx, k, feeCollector, msg)
		case *types.MsgAddAccountCertificates:
			return handlerMsgAddAccountCertificates(ctx, k, feeCollector, msg)
		case *types.MsgDeleteAccountCertificate:
			return handlerMsgDeleteAccountCertificate(ctx, k, feeCollector, msg)
		case *types.MsgDeleteAccount:
			return handlerMsgDeleteAccount(ctx, k, feeCollector, msg)
		case *types.MsgReplaceAccountTargets:
			return handlerMsgReplaceAccountTargets(ctx, k, feeCollector, msg)
		case *types.MsgTransferAccount:
			return handlerMsgTransferAccount(ctx, k, feeCollector, msg)
		case *types.MsgReplaceAccountMetadata:
			return handlerMsgReplaceAccountMetadata(ctx, k, feeCollector, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unregonized request: %T", msg))
		}
	}

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		/*
			TODO
			remove when cosmos sdk decides that you are allowed to panic on errors that should not happen
			instead of returning random internal errors that mean actually nothing to a developer without
			a stacktrace or at least the error string of the panic itself, and also substitute 'log' stdlib
			with cosmos sdk logger when they make clear how you can use it and how to set up env to achieve so
		*/
		defer func() {
			if r := recover(); r != nil {
				log.Printf("FATAL-PANIC while executing message: %#v\nReason: %v", msg, r)
				// and lets panic again to throw it back to cosmos sdk yikes.
				panic(r)
			}
		}()
		resp, err := f(ctx, msg)
		if err != nil {
			msg := fmt.Sprintf("tx rejected %T: %s", msg, err)
			k.Logger(ctx).With("module", types.ModuleName).Info(msg)
		}
		return resp, err
	}
}
