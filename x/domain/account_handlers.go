package domain

import (
	"fmt"
	"regexp"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/configuration"
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
	if err := accountCtrl.Validate(account.MustExist, account.NotExpired, account.Owner(msg.Owner), account.CertificateNotExist(msg.NewCertificate)); err != nil {
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
	// perform account checks, save certificate index
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	certIndex := new(int)
	if err := accountCtrl.Validate(account.MustExist, account.Owner(msg.Owner), account.CertificateExists(msg.DeleteCertificate, certIndex)); err != nil {
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
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(account.MustExist); err != nil {
		return nil, err
	}
	// perform action authorization checks
	if (domainCtrl.Validate(domain.Owner(msg.Owner)) != nil) && (accountCtrl.Validate(account.Owner(msg.Owner)) != nil) {
		return nil, errors.Wrapf(types.ErrUnauthorized, "only account owner: %s and domain admin %s can delete the account", accountCtrl.Account().Owner, domainCtrl.Domain().Admin)
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

// handleMsgRegisterAccount registers the domain
func handleMsgRegisterAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRegisterAccount) (*sdk.Result, error) {
	// verify request
	// get config
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	// validate blockchain targets
	if err := validateBlockchainTargets(msg.Targets, conf); err != nil {
		return nil, errors.Wrap(types.ErrInvalidBlockchainTarget, err.Error())
	}
	// do validity checks on domain
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	err := domainCtrl.Validate(domain.MustExist, domain.Type(types.ClosedDomain), domain.NotExpired, domain.Owner(msg.Owner))
	if err != nil {
		return nil, err
	}
	// get domain
	domain := domainCtrl.Domain()
	// accounts validity checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	err = accountCtrl.Validate(account.ValidName, account.MustNotExist)
	if err != nil {
		return nil, err
	}
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// create account struct
	account := types.Account{
		Domain:       msg.Domain,
		Name:         msg.Name,
		Owner:        msg.Owner,
		ValidUntil:   ctx.BlockTime().Add(domain.AccountRenew * time.Second).Unix(), // add curr block time + domain account renew and convert to unix seconds
		Targets:      msg.Targets,
		Certificates: nil,
		Broker:       msg.Broker,
	}
	// save account
	k.CreateAccount(ctx, account)
	// success; TODO can we emit events?
	return &sdk.Result{}, nil
}

// validateBlockchainTargets validates different blockchain targets address and ID
func validateBlockchainTargets(targets []types.BlockchainAddress, conf configuration.Config) error {
	validBlockchainID := regexp.MustCompile(conf.ValidBlockchainID)
	validBlockchainAddress := regexp.MustCompile(conf.ValidBlockchainAddress)
	// create blockchain targets set to identify duplicates
	sets := make(map[string]struct{}, len(targets))
	// iterate over targets to check their validity
	for _, target := range targets {
		// check if blockchain ID was already specified
		if _, ok := sets[target.ID]; ok {
			return fmt.Errorf("duplicate blockchain ID: %s", target)
		}
		sets[target.ID] = struct{}{}
		// is blockchain id valid?
		if !validBlockchainID.MatchString(target.ID) {
			return fmt.Errorf("%s is not a valid blockchain ID", target.ID)
		}
		// is blockchain address valid?
		if !validBlockchainAddress.MatchString(target.Address) {
			return fmt.Errorf("%s is not a valid blockchain address", target.Address)
		}
	}
	// success
	return nil
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
	k.UpdateAccountValidity(ctx, accountCtrl.Account(), domainCtrl.Domain().AccountRenew)
	// success; todo emit event??
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountTargets replaces account targets
func handlerMsgReplaceAccountTargets(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountTargets) (*sdk.Result, error) {
	// get configuration
	config := k.ConfigurationKeeper.GetConfiguration(ctx)
	// validate blockchain targets
	err := validateBlockchainTargets(msg.NewTargets, config)
	if err != nil {
		return nil, errors.Wrapf(types.ErrInvalidBlockchainTarget, err.Error())
	}
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(account.MustExist, account.NotExpired, account.Owner(msg.Owner)); err != nil {
		return nil, err
	}
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// replace targets replaces accounts targets
	k.ReplaceAccountTargets(ctx, accountCtrl.Account(), msg.NewTargets)
	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}

// handlerMsgSetAccountMetadata takes care of setting account metadata
func handlerMsgSetAccountMetadata(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSetAccountMetadata) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(account.MustExist, account.NotExpired, account.Owner(msg.Owner)); err != nil {
		return nil, err
	}
	// update account
	account := accountCtrl.Account()
	account.MetadataURI = msg.NewMetadataURI
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// save to store
	k.SetAccount(ctx, account)
	// success TODO emit event
	return &sdk.Result{}, nil
}

// handlerMsgTransferAccount transfers account to a new owner
// after clearing targets and certificates
func handlerMsgTransferAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferAccount) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := domain.NewController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(domain.MustExist, domain.NotExpired); err != nil {
		return nil, err
	}
	// check if account exists
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(account.MustExist, account.NotExpired); err != nil {
		return nil, err
	}
	// check if domain has super user
	switch domainCtrl.Domain().Type {
	// if it has a super user then only domain admin can transfer accounts
	case types.ClosedDomain:
		if domainCtrl.Validate(domain.Owner(msg.Owner)) != nil {
			return nil, errors.Wrapf(types.ErrUnauthorized, "only domain admin %s is allowed to transfer accounts", domainCtrl.Domain().Admin)
		}
	// if it has not a super user then only account owner can transfer the account
	case types.OpenDomain:
		if accountCtrl.Validate(account.Owner(msg.Owner)) != nil {
			return nil, errors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to transfer the account", accountCtrl.Account().Owner)
		}
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// transfer account
	k.TransferAccount(ctx, accountCtrl.Account(), msg.NewOwner)
	// success, todo emit event?
	return &sdk.Result{}, nil
}
