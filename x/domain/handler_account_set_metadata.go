package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// handlerMsgSetAccountMetadata takes care of setting account metadata
func handlerMsgSetAccountMetadata(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSetAccountMetadata) (*sdk.Result, error) {
	// get domain
	domain, ok := k.GetDomain(ctx, msg.Domain)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if domain expired
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "domain %s has expired", domain.Name)
	}
	// get account
	account, ok := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if account expired
	if ctx.BlockTime().After(iovns.SecondsToTime(account.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if signer is account owner
	if !account.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "not allowed to change account metadata uri, invalid owner: %s", msg.Owner)
	}
	// update account
	account.MetadataURI = msg.NewMetadataURI
	// save to store
	k.SetAccount(ctx, account)
	// success
	return &sdk.Result{}, nil
}
