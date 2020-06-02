package domain

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/controllers/account"
	"github.com/iov-one/iovns/x/domain/controllers/domain"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerMsgAddAccountCertificates(ctx sdk.Context, k keeper.Keeper, msg *types.MsgAddAccountCertificates) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}

	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	d := domainCtrl.Domain()
	switch d.Type {
	case types.OpenDomain:
		if err := accountCtrl.Validate(account.NotExpired); err != nil {
			return nil, err
		}
	}
	if err := accountCtrl.Validate(
		account.MustExist,
		account.Owner(msg.Owner),
		account.CertificateLimitNotExceeded,
		account.CertificateSizeNotExceeded(msg.NewCertificate),
		account.CertificateNotExist(msg.NewCertificate)); err != nil {
		return nil, err
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}
	// add certificate
	k.AddAccountCertificate(ctx, accountCtrl.Account(), msg.NewCertificate)
	// success; TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgDeleteAccountCertificate(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccountCertificate) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(
		domain.MustExist,
		domain.NotExpired); err != nil {
		return nil, err
	}
	// perform account checks, save certificate index
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	d := domainCtrl.Domain()
	switch d.Type {
	case types.OpenDomain:
		if err := accountCtrl.Validate(account.NotExpired); err != nil {
			return nil, err
		}
	}
	certIndex := new(int)
	if err := accountCtrl.Validate(
		account.MustExist,
		account.Owner(msg.Owner),
		account.CertificateExists(msg.DeleteCertificate, certIndex)); err != nil {
		return nil, err
	}
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// delete cert
	k.DeleteAccountCertificate(ctx, accountCtrl.Account(), *certIndex)
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
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// delete account
	k.DeleteAccount(ctx, msg.Domain, msg.Name)
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
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(
		account.ValidName,
		account.MustNotExist,
		account.ValidTargets(msg.Targets),
	); err != nil {
		return nil, err
	}

	a := types.Account{
		Domain:       msg.Domain,
		Name:         msg.Name,
		Owner:        msg.Owner,
		ValidUntil:   ctx.BlockTime().Add(d.AccountRenew * time.Second).Unix(),
		Targets:      msg.Targets,
		Certificates: nil,
		Broker:       msg.Broker,
	}
	switch d.Type {
	case types.ClosedDomain:
		if err := domainCtrl.Validate(domain.Admin(msg.Registerer)); err != nil {
			return nil, err
		}
		a.ValidUntil = types.MaxValidUntil
	case types.OpenDomain:
		a.ValidUntil = ctx.BlockTime().Add(conf.AccountRenewalPeriod).Unix()
	}

	err := k.CollectFees(ctx, msg, msg.Registerer)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	k.CreateAccount(ctx, a)
	return &sdk.Result{}, nil
}

func handlerMsgRenewAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewAccount) (*sdk.Result, error) {
	// validate domain
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist); err != nil {
		return nil, err
	}
	// validate account
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(account.MustExist); err != nil {
		return nil, err
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Signer)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// renew account
	k.RenewAccount(ctx, accountCtrl.Account(), domainCtrl.Domain().AccountRenew)
	// success; todo emit event??
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountTargets replaces account targets
func handlerMsgReplaceAccountTargets(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountTargets) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(
		account.MustExist,
		account.Owner(msg.Owner),
		account.ValidTargets(msg.NewTargets),
		account.BlockchainTargetLimitNotExceeded(msg.NewTargets),
	); err != nil {
		return nil, err
	}

	d := domainCtrl.Domain()
	switch d.Type {
	case types.OpenDomain:
		if err := accountCtrl.Validate(account.NotExpired); err != nil {
			return nil, err
		}
	case types.ClosedDomain:
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// replace targets replaces accounts targets
	k.ReplaceAccountTargets(ctx, accountCtrl.Account(), msg.NewTargets)
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
		account.Owner(msg.Owner),
		account.MetadataSizeNotExceeded(msg.NewMetadataURI)); err != nil {
		return nil, err
	}
	d := domainCtrl.Domain()
	switch d.Type {
	case types.OpenDomain:
		if err := accountCtrl.Validate(account.NotExpired); err != nil {
			return nil, err
		}
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// save to store
	k.UpdateMetadataAccount(ctx, accountCtrl.Account(), msg.NewMetadataURI)
	// success TODO emit event
	return &sdk.Result{}, nil
}

// handlerMsgTransferAccount transfers account to a new owner
// after clearing targets and certificates
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

	d := domainCtrl.Domain()
	switch d.Type {
	case types.OpenDomain:

	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// transfer account
	k.TransferAccountWithReset(ctx, accountCtrl.Account(), msg.NewOwner, msg.Reset)
	// success, todo emit event?
	return &sdk.Result{}, nil
}
