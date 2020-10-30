package starname

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/starname/controllers/domain"
	"github.com/iov-one/iovns/x/starname/controllers/fees"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/keeper/executor"
	"github.com/iov-one/iovns/x/starname/types"
)

func handlerMsgDeleteDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteDomain) (*sdk.Result, error) {
	ctrl := domain.NewController(ctx, k, msg.Domain)
	// do precondition and authorization checks
	if err := ctrl.
		MustExist().
		DeletableBy(msg.Owner).
		Validate(); err != nil {
		return nil, err
	}
	// operation is allowed
	feeCtrl := fees.NewController(ctx, k, ctrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// all checks passed delete domain
	executor.NewDomain(ctx, k, ctrl.Domain()).Delete()
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgRegisterDomain handles the domain registration process
func handleMsgRegisterDomain(ctx sdk.Context, k Keeper, msg *types.MsgRegisterDomain) (resp *sdk.Result, err error) {
	ctrl := domain.NewController(ctx, k, msg.Name)
	err = ctrl.
		MustNotExist().
		ValidName().
		Validate()
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
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Admin.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDomainType, (string)(msg.DomainType)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Admin.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handlerMsgRenewDomain renews a domain
func handlerMsgRenewDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewDomain) (*sdk.Result, error) {
	ctrl := domain.NewController(ctx, k, msg.Domain)
	err := ctrl.
		MustExist().
		Renewable().
		Validate()
	if err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, ctrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err = k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// update domain
	executor.NewDomain(ctx, k, ctrl.Domain()).Renew()
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Signer.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handlerMsgTransferDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferDomain) (*sdk.Result, error) {
	c := domain.NewController(ctx, k, msg.Domain)
	err := c.
		MustExist().
		Admin(msg.Owner).
		NotExpired().
		Transferable(msg.TransferFlag).
		Validate()
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
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyTransferDomainNewOwner, msg.NewAdmin.String()),
			sdk.NewAttribute(types.AttributeKeyTransferDomainFlag, fmt.Sprintf("%d", msg.TransferFlag)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
