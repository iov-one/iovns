package starname

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/starname/controllers/fees"
	"github.com/iov-one/iovns/x/starname/keeper/executor"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/starname/controllers/domain"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

func handlerMsgDeleteDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	// do precondition and authorization checks
	if err := c.Validate(domain.MustExist, domain.DeletableBy(msg.Owner)); err != nil {
		return nil, err
	}
	// operation is allowed
	feeCtrl := fees.NewController(ctx, k, c.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// all checks passed delete domain
	executor.NewDomain(ctx, k, c.Domain()).Delete()
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
	// create new domain
	d := types.Domain{
		Name:       msg.Name,
		Admin:      msg.Admin,
		ValidUntil: ctx.BlockTime().Add(k.ConfigurationKeeper.GetDomainRenewDuration(ctx)).Unix(),
		Type:       msg.DomainType,
		Broker:     msg.Broker,
	}
	feeCtrl := fees.NewController(ctx, k, d)
	fee := feeCtrl.GetFee(msg)
	// collect fees
	if err := k.CollectFees(ctx, msg, fee); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// save domain
	ex := executor.NewDomain(ctx, k, d)
	ex.Create()
	// success TODO think here, can we emit any useful event
	return &sdk.Result{}, nil
}

// handlerMsgRenewDomain renews a domain
func handlerMsgRenewDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	err := c.Validate(domain.MustExist, domain.Renewable)
	if err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, c.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err = k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// update domain
	executor.NewDomain(ctx, k, c.Domain()).Renew()
	// success TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgTransferDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	err := c.Validate(
		domain.MustExist,
		domain.Type(types.ClosedDomain),
		domain.Admin(msg.Owner),
		domain.NotExpired,
		domain.Transferable(msg.TransferFlag),
	)
	if err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, c.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err = k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	ex := executor.NewDomain(ctx, k, c.Domain())
	ex.Transfer(msg.TransferFlag, msg.NewAdmin)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}
