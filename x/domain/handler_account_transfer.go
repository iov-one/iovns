package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/iov-one/iovnsd/x/domain/types"
)

// handlerMsgTransferAccount transfers account to a new owner
// after clearing targets and certificates
func handlerMsgTransferAccount(ctx sdk.Context, k keeper.Keeper, msg types.MsgTransferAccount) (*sdk.Result, error) {
	// check if domain exists
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "%s does not exist", msg.Domain)
	}
	// check if domain has expired expired
	if iovnsd.SecondsToTime(domain.ValidUntil).Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "account transfer is not allowed for expired domains, expire date: %s", iovnsd.SecondsToTime(domain.ValidUntil))
	}
	// check if account exists
	account, exists := k.GetAccount(ctx, iovnsd.GetAccountKey(msg.Domain, msg.Name))
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "account %s does not exist", msg.Name)
	}
	// check if account has expired
	if iovnsd.SecondsToTime(account.ValidUntil).Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if domain has super user
	switch domain.HasSuperuser {
	// if it has a super user then only domain admin can transfer accounts
	case true:
		if !msg.Owner.Equals(domain.Admin) {
			return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only domain admin %s is allowed to transfer accounts", domain.Admin)
		}
	// if it has not a super user then only account owner can transfer the account
	case false:
		if !msg.Owner.Equals(account.Owner) {
			return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to transfer the account", account.Owner)
		}
	}
	// now transfer account
	account.Owner = msg.NewOwner // transfer owner
	account.Targets = nil        // clear blockchain targets
	account.Certificates = nil   // clear certificates
	// success, todo emit event?
	return &sdk.Result{}, nil
}
