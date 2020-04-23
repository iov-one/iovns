package domain

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgAddAccountCertificates(ctx sdk.Context, k keeper.Keeper, msg types.MsgAddAccountCertificates) (*sdk.Result, error) {
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if current time is after domain validity time
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "domain %s has expired", msg.Domain)
	}
	// get account
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if current time is after account validity time
	if ctx.BlockTime().After(iovns.SecondsToTime(account.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if signer is account owner
	if !msg.Owner.Equals(account.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s cannot add certificates to account owned by %s", msg.Owner, account.Owner)
	}
	// check if certificate is already present in account
	for _, cert := range account.Certificates {
		if bytes.Equal(cert, msg.NewCertificate) {
			return nil, sdkerrors.Wrapf(types.ErrCertificateExists, "certificate is already present")
		}
	}
	// add certificate
	k.AddAccountCertificate(ctx, account, msg.NewCertificate)
	// success; TODO emit event
	return &sdk.Result{}, nil
}
