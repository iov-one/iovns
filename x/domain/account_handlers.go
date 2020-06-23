package domain

import (
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/feecalculator"
	fee "github.com/iov-one/iovns/x/fee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/controllers/account"
	"github.com/iov-one/iovns/x/domain/controllers/domain"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgAddAccountCertificates(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgAddAccountCertificates) (*sdk.Result, error) {
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
	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// add certificate
	k.AddAccountCertificate(ctx, accountCtrl.Account(), msg.NewCertificate)
	// success; TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgDeleteAccountCertificate(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgDeleteAccountCertificate) (*sdk.Result, error) {
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
	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// delete cert
	k.DeleteAccountCertificate(ctx, accountCtrl.Account(), *certIndex)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}

// handlerMsgDelete account deletes the account from the system
func handlerMsgDeleteAccount(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgDeleteAccount) (*sdk.Result, error) {
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
	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// delete account
	k.DeleteAccount(ctx, msg.Domain, msg.Name)
	// success; todo can we emit event?
	return &sdk.Result{}, nil
}

// handleMsgRegisterAccount registers the account
func handleMsgRegisterAccount(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgRegisterAccount) (*sdk.Result, error) {
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
		account.ValidTargets(msg.Targets),
		account.RegistrableBy(msg.Registerer),
	); err != nil {
		return nil, err
	}

	a := types.Account{
		Domain:       msg.Domain,
		Name:         msg.Name,
		Owner:        msg.Owner,
		Targets:      msg.Targets,
		Certificates: nil,
		Broker:       msg.Broker,
	}
	switch d.Type {
	case types.ClosedDomain:
		a.ValidUntil = types.MaxValidUntil
	case types.OpenDomain:
		a.ValidUntil = ctx.BlockTime().Add(conf.AccountRenewalPeriod).Unix()
	}

	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	k.CreateAccount(ctx, a)
	return &sdk.Result{}, nil
}

func handlerMsgRenewAccount(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgRenewAccount) (*sdk.Result, error) {
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
		account.MaxRenewNotExceeded); err != nil {
		return nil, err
	}
	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// renew account
	a := accountCtrl.Account()
	// account valid until is extended here
	k.RenewAccount(ctx, &a, conf.AccountRenewalPeriod)
	// get grace period and expiration time
	d := domainCtrl.Domain()
	dgp := conf.DomainGracePeriod
	domainGracePeriodUntil := iovns.SecondsToTime(d.ValidUntil).Add(dgp)
	accNewValidUntil := iovns.SecondsToTime(a.ValidUntil)
	if domainGracePeriodUntil.Before(accNewValidUntil) {
		d.ValidUntil = accNewValidUntil.Unix()
		k.SetDomain(ctx, d)
	}
	// success; todo emit event??
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountTargets replaces account targets
func handlerMsgReplaceAccountTargets(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgReplaceAccountTargets) (*sdk.Result, error) {
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
		account.ValidTargets(msg.NewTargets),
		account.BlockchainTargetLimitNotExceeded(msg.NewTargets),
	); err != nil {
		return nil, err
	}
	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// replace targets replaces accounts targets
	k.ReplaceAccountTargets(ctx, accountCtrl.Account(), msg.NewTargets)
	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountMetadata takes care of setting account metadata
func handlerMsgReplaceAccountMetadata(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgReplaceAccountMetadata) (*sdk.Result, error) {
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
	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// save to store
	k.UpdateMetadataAccount(ctx, accountCtrl.Account(), msg.NewMetadataURI)
	// success TODO emit event
	return &sdk.Result{}, nil
}

// handlerMsgTransferAccount transfers account to a new owner
// after clearing targets and certificates
func handlerMsgTransferAccount(ctx sdk.Context, k keeper.Keeper, collector fee.Collector, msg *types.MsgTransferAccount) (*sdk.Result, error) {
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
	// calculate fees
	calc := feecalculator.NewFeeCalculator(ctx, k).WithDomain(domainCtrl.Domain()).WithAccount(accountCtrl.Account())
	f, err := calc.CalculateFee(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to calculate fee")
	}
	// collect fees
	if err := collector.CollectFee(ctx, f, msg.FeePayer()); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// transfer account
	k.TransferAccountWithReset(ctx, accountCtrl.Account(), msg.NewOwner, msg.Reset)
	// success, todo emit event?
	return &sdk.Result{}, nil
}
