package domain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// handlerMsgRenewDomain renews a domain
func handlerMsgRenewDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewDomain) (*sdk.Result, error) {
	// check if domain exists
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found %s", msg.Domain)
	}
	// get configuration
	renewDuration := k.ConfigurationKeeper.GetDomainRenewDuration(ctx)
	// update domain valid until
	domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// update domain
	k.SetDomain(ctx, domain)
	// success
	return &sdk.Result{}, nil
}
