package starname

import (
	"fmt"
	"strconv"

	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/controllers/fees"
	"github.com/iov-one/iovns/x/starname/keeper/executor"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/starname/controllers/account"
	"github.com/iov-one/iovns/x/starname/controllers/domain"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

func handlerMsgAddAccountCertificates(ctx sdk.Context, k keeper.Keeper, msg *types.MsgAddAccountCertificates) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)

	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		CertificateLimitNotExceeded().
		CertificateSizeNotExceeded(msg.NewCertificate).
		CertificateNotExist(msg.NewCertificate).
		Validate(); err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}
	// add certificate
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.AddCertificate(msg.NewCertificate)
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyNewCertificate, fmt.Sprintf("%x", msg.NewCertificate)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handlerMsgDeleteAccountCertificate(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccountCertificate) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}
	// perform account checks, save certificate index
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	certIndex := new(int)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		CertificateExists(msg.DeleteCertificate, certIndex).
		Validate(); err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// delete cert
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.DeleteCertificate(*certIndex)
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDeletedCertificate, fmt.Sprintf("%x", msg.DeleteCertificate)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handlerMsgDelete account deletes the account from the system
func handlerMsgDeleteAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccount) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.MustExist().Validate(); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)
	if err := accountCtrl.
		MustExist().
		DeletableBy(msg.Owner).
		Validate(); err != nil {
		return nil, err
	}
	// collect fees
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// delete account
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.Delete()
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handleMsgRegisterAccount registers the account
func handleMsgRegisterAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRegisterAccount) (*sdk.Result, error) {
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	domainCtrl := domain.NewController(ctx, k, msg.Domain).WithConfiguration(conf)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}
	d := domainCtrl.Domain()
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)
	if err := accountCtrl.
		ValidName().
		MustNotExist().
		ValidResources(msg.Resources).
		RegistrableBy(msg.Registerer).
		Validate(); err != nil {
		return nil, err
	}

	a := types.Account{
		Domain:       msg.Domain,
		Name:         utils.StrPtr(msg.Name),
		Owner:        msg.Owner,
		Resources:    msg.Resources,
		Certificates: nil,
		Broker:       msg.Broker,
	}
	switch d.Type {
	case types.ClosedDomain:
		a.ValidUntil = types.MaxValidUntil
	case types.OpenDomain:
		a.ValidUntil = ctx.BlockTime().Add(conf.AccountRenewalPeriod).Unix()
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	ex := executor.NewAccount(ctx, k, a)
	ex.Create()
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyBroker, msg.Broker.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

func handlerMsgRenewAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewAccount) (*sdk.Result, error) {
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	// validate domain
	domainCtrl := domain.NewController(ctx, k, msg.Domain).WithConfiguration(conf)
	if err := domainCtrl.MustExist().Type(types.OpenDomain).Validate(); err != nil {
		return nil, err
	}
	// validate account
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		Renewable().
		Validate(); err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// renew account
	// account valid until is extended here
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.Renew()
	// get grace period and expiration time
	d := domainCtrl.Domain()
	dgp := conf.DomainGracePeriod
	domainGracePeriodUntil := utils.SecondsToTime(d.ValidUntil).Add(dgp)
	accNewValidUntil := utils.SecondsToTime(ex.State().ValidUntil)
	if domainGracePeriodUntil.Before(accNewValidUntil) {
		dex := executor.NewDomain(ctx, k, domainCtrl.Domain())
		dex.Renew(accNewValidUntil.Unix())
	}
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Signer.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handlerMsgReplaceAccountResources replaces account resources
func handlerMsgReplaceAccountResources(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountResources) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		ValidResources(msg.NewResources).
		ResourceLimitNotExceeded(msg.NewResources).
		Validate(); err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// replace accounts resources
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.ReplaceResources(msg.NewResources)
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyNewResources, ""), // TODO stringify resources
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handlerMsgReplaceAccountMetadata takes care of setting account metadata
func handlerMsgReplaceAccountMetadata(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountMetadata) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		MetadataSizeNotExceeded(msg.NewMetadataURI).
		Validate(); err != nil {
		return nil, err
	}
	// collect fees
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// save to store
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.UpdateMetadata(msg.NewMetadataURI)
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyNewMetadata, msg.NewMetadataURI),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// handlerMsgTransferAccount transfers account to a new owner
// after clearing resources and certificates
func handlerMsgTransferAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferAccount) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}
	// check if account exists
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		TransferableBy(msg.Owner).
		ResettableBy(msg.Owner, msg.Reset).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// transfer account
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.Transfer(msg.NewOwner, msg.Reset)
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyTransferAccountNewOwner, msg.NewOwner.String()),
			sdk.NewAttribute(types.AttributeKeyTransferAccountReset, strconv.FormatBool(msg.Reset)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
