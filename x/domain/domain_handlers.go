package domain

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/controllers/domain"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgDeleteDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	// do precondition checks
	if err := c.Validate(domain.MustExist, domain.Type(types.ClosedDomain)); err != nil {
		return nil, err
	}
	// check if signer is authorized to delete
	if err := c.Validate(domain.DeletableBy(msg.Owner)); err != nil {
		return nil, err
	}
	// operation is allowed
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
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
	// save domain
	k.CreateDomain(ctx, d)
	// generate empty name account
	acc := types.Account{
		Domain:       msg.Name,
		Name:         "",
		Owner:        msg.Admin, // TODO this is not clear, why the domain admin is zero address while this is msg.Admin
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
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// update domain
	k.RenewDomain(ctx, c.Domain())
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
