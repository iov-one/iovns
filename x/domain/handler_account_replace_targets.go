package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// handlerMsgReplaceAccountTargets replaces account targets
func handlerMsgReplaceAccountTargets(ctx sdk.Context, k keeper.Keeper, msg types.MsgReplaceAccountTargets) (*sdk.Result, error) {
	// get configuration
	config := k.ConfigurationKeeper.GetConfiguration(ctx)
	// validate blockchain targets
	err := validateBlockchainTargets(msg.NewTargets, config)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidBlockchainTarget, err.Error())
	}
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// see if domain still valid
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "domain %s has expired", msg.Domain)
	}
	// get account
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found: %s", msg.Name)
	}
	// check if expired
	if ctx.BlockTime().After(iovns.SecondsToTime(account.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if account owner matches request signer
	if !msg.Owner.Equals(account.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "account %s is not authorized to perform actions on account owned by %s", msg.Owner, account.Owner)
	}
	// replace targets
	account.Targets = msg.NewTargets
	// update account
	k.SetAccount(ctx, account)
	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}
