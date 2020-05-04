package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgFlushDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgFlushDomain) (*sdk.Result, error) {
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if domain has superuser
	if !domain.HasSuperuser {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "domains without a superuser cannot be flushed")
	}
	// check if signer is also domain owner
	if !msg.Owner.Equals(domain.Admin) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to flush domain owned by %s", msg.Owner, domain.Admin)
	}
	// now flush
	_ = k.FlushDomain(ctx, msg.Domain)
	// success; TODO maybe emit event?
	return &sdk.Result{}, nil
}
