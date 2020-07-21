package starname

import (
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
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}

	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)

	if err := accountCtrl.Validate(
		account.MustExist,
		account.NotExpired,
		account.Owner(msg.Owner),
		account.CertificateLimitNotExceeded,
		account.CertificateSizeNotExceeded(msg.NewCertificate),
		account.CertificateNotExist(msg.NewCertificate),
	); err != nil {
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
	// success; TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgDeleteAccountCertificate(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccountCertificate) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(
		domain.MustExist,
		domain.NotExpired,
	); err != nil {
		return nil, err
	}
	// perform account checks, save certificate index
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	certIndex := new(int)
	if err := accountCtrl.Validate(
		account.MustExist,
		account.NotExpired,
		account.Owner(msg.Owner),
		account.CertificateExists(msg.DeleteCertificate, certIndex),
	); err != nil {
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
	// success; TODO emit event?
	return &sdk.Result{}, nil
}

// handlerMsgDelete account deletes the account from the system
func handlerMsgDeleteAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccount) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)
	if err := accountCtrl.Validate(account.MustExist, account.DeletableBy(msg.Owner)); err != nil {
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
	// success; todo can we emit event?
	return &sdk.Result{}, nil
}

// handleMsgRegisterAccount registers the account
func handleMsgRegisterAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRegisterAccount) (*sdk.Result, error) {
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	domainCtrl := domain.NewController(ctx, k, msg.Domain).WithConfiguration(conf)
	if err := domainCtrl.Validate(
		domain.MustExist,
		domain.NotExpired,
	); err != nil {
		return nil, err
	}
	d := domainCtrl.Domain()
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)
	if err := accountCtrl.Validate(
		account.ValidName,
		account.MustNotExist,
		account.ValidResources(msg.Resources),
		account.RegistrableBy(msg.Registerer),
	); err != nil {
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
	return &sdk.Result{}, nil
}

func handlerMsgRenewAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewAccount) (*sdk.Result, error) {
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	// validate domain
	domainCtrl := domain.NewController(ctx, k, msg.Domain).WithConfiguration(conf)
	if err := domainCtrl.Validate(domain.MustExist, domain.Type(types.OpenDomain)); err != nil {
		return nil, err
	}
	// validate account
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).WithConfiguration(conf)
	if err := accountCtrl.Validate(
		account.MustExist,
		account.Renewable); err != nil {
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
	// success; todo emit event??
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountResources replaces account resources
func handlerMsgReplaceAccountResources(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountResources) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(
		account.MustExist,
		account.NotExpired,
		account.Owner(msg.Owner),
		account.ValidResources(msg.NewResources),
		account.ResourceLimitNotExceeded(msg.NewResources),
	); err != nil {
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
	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountMetadata takes care of setting account metadata
func handlerMsgReplaceAccountMetadata(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountMetadata) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(
		account.MustExist,
		account.NotExpired,
		account.Owner(msg.Owner),
		account.MetadataSizeNotExceeded(msg.NewMetadataURI)); err != nil {
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
	// success TODO emit event
	return &sdk.Result{}, nil
}

// handlerMsgTransferAccount transfers account to a new owner
// after clearing resources and certificates
func handlerMsgTransferAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferAccount) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(
		domain.MustExist,
		domain.NotExpired,
	); err != nil {
		return nil, err
	}
	// check if account exists
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name).
		WithDomainController(domainCtrl)
	if err := accountCtrl.Validate(
		account.MustExist,
		account.NotExpired,
		account.TransferableBy(msg.Owner),
		account.ResettableBy(msg.Owner, msg.Reset),
	); err != nil {
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
	// success, todo emit event?
	return &sdk.Result{}, nil
}
