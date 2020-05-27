package domain

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/controllers/domain"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgDeleteDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	err := c.Validate(domain.MustExist, domain.Type(types.ClosedDomain))
	if errors.Is(err, types.ErrInvalidDomainType) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "user is unauthorized to delete domain %s with domain type: %s", msg.Domain, types.ClosedDomain)
	}
	if err != nil {
		return nil, err
	}
	// if domain is not over grace period and signer is not the owner of the domain then the operation is not allowed
	if err := c.Validate(domain.Owner(msg.Owner)); err != nil && !c.Condition(domain.GracePeriodFinished) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorized, "unable to delete domain not owned if grace period is not finished")
	}
	// operation is allowed
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// all checks passed delete domain
	_ = k.DeleteDomain(ctx, msg.Domain)
	// success TODO maybe emit event?
	return &sdk.Result{}, nil
}

// handleMsgRegisterDomain handles the domain registration process
func handleMsgRegisterDomain(ctx sdk.Context, k Keeper, msg *types.MsgRegisterDomain) (resp *sdk.Result, err error) {
	c := domain.NewController(ctx, k, msg.Name)
	err = c.Validate(domain.MustNotExist, domain.ValidName)
	if err != nil {
		return nil, err
	}
	// set new domain
	d := types.Domain{
		Name:         msg.Name,
		Admin:        msg.Admin,
		ValidUntil:   ctx.BlockTime().Add(k.ConfigurationKeeper.GetDomainRenewDuration(ctx)).Unix(),
		Type:         msg.DomainType,
		AccountRenew: msg.AccountRenew,
		Broker:       msg.Broker,
	}

	// generate empty name account
	acc := types.Account{
		Domain:       msg.Name,
		Name:         "",
		Owner:        msg.Admin,
		ValidUntil:   ctx.BlockTime().Add(d.AccountRenew).Unix(),
		Targets:      nil,
		Certificates: nil,
		Broker:       nil, // TODO ??
	}
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Admin)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// save domain
	k.CreateDomain(ctx, d)
	// save account
	k.CreateAccount(ctx, acc)
	// success TODO think here, can we emit any useful event
	return &sdk.Result{}, nil
}

// handlerMsgRenewDomain renews a domain
func handlerMsgRenewDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	err := c.Validate(domain.MustExist)
	if err != nil {
		return nil, err
	}
	domain := c.Domain()
	// get configuration
	renewDuration := k.ConfigurationKeeper.GetDomainRenewDuration(ctx)
	// update domain valid until
	domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// update domain
	k.SetDomain(ctx, domain)
	// success TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgTransferDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	err := c.Validate(
		domain.MustExist,
		domain.Type(types.ClosedDomain),
		domain.Owner(msg.Owner),
		domain.NotExpired,
	)
	if types.ErrInvalidDomainType.Is(err) {
		return nil, types.ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	// get domain
	d := c.Domain()
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// transfer domain and accounts ownership
	k.TransferDomain(ctx, msg.NewAdmin, d)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}
