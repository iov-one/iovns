package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgTransferDomain(ctx sdk.Context, k keeper.Keeper, msg types.MsgTransferDomain) (*sdk.Result, error) {
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if has superuser
	if !domain.HasSuperuser {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "domain %s without superuser cannot be transferred", msg.Domain)
	}
	// check if signer is domain owner
	if !msg.Owner.Equals(domain.Admin) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to transfer domain owned by %s", msg.Owner, domain.Admin)
	}
	// check if domain is valid
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "%s has expired", msg.Domain)
	}
	// change domain admin
	domain.Admin = msg.NewAdmin
	// save domain
	k.SetDomain(ctx, domain)
	// get account keys related to the domain
	accountKeys := k.GetAccountsInDomain(ctx, domain.Name)
	// iterate over accounts
	for _, accountKey := range accountKeys {
		// get account; TODO might change the need to convert back from bytes to string to bytes
		account, _ := k.GetAccount(ctx, string(accountKey))
		// update account
		account.Certificates = nil   // delete certs
		account.Targets = nil        // delete targets
		account.Owner = msg.NewAdmin // update admin
		// save to kvstore
		k.SetAccount(ctx, account)
	}
	// success; TODO emit event?
	return &sdk.Result{}, nil
}
