package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// handlerMsgDelete account deletes the account from the system
func handlerMsgDeleteAccount(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteAccount) (*sdk.Result, error) {
	// check if domain exists
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if account exists
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found: %s", msg.Name)
	}
	// check if msg.Owner is either domain owner or account owner
	if !domain.Admin.Equals(msg.Owner) && !account.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner: %s and domain admin %s can delete the account", account.Owner, domain.Admin)
	}
	// delete account
	k.DeleteAccount(ctx, msg.Domain, msg.Name)
	// success; todo can we emit event?
	return &sdk.Result{}, nil
}
