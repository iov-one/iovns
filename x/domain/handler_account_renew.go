package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgRenewAccount(ctx sdk.Context, k keeper.Keeper, msg types.MsgRenewAccount) (*sdk.Result, error) {
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// get account
	account, exists := k.GetAccount(ctx, iovns.GetAccountKey(msg.Domain, msg.Name))
	if exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found: %s", msg.Name)
	}
	// update account time
	account.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(account.ValidUntil).Add(domain.AccountRenew),
	)
	// update account in kv store
	k.SetAccount(ctx, account)
	// success; todo emit event??
	return &sdk.Result{}, nil
}
