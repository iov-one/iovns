package domain

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/iov-one/iovns/x/domain/controllers/account"

	ctrl "github.com/iov-one/iovns/x/domain/controllers"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// NewHandler builds the tx requests handler for the domain module
func NewHandler(k Keeper) sdk.Handler {
	f := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		// domain handlers
		case *types.MsgRegisterDomain:
			return handleMsgRegisterDomain(ctx, k, msg)
		case *types.MsgRenewDomain:
			return handlerMsgRenewDomain(ctx, k, msg)
		case *types.MsgDeleteDomain:
			return handlerMsgDeleteDomain(ctx, k, msg)
		case *types.MsgTransferDomain:
			return handlerMsgTransferDomain(ctx, k, msg)
		// account handlers
		case *types.MsgRegisterAccount:
			return handleMsgRegisterAccount(ctx, k, msg)
		case *types.MsgRenewAccount:
			return handlerMsgRenewAccount(ctx, k, msg)
		case *types.MsgAddAccountCertificates:
			return handlerMsgAddAccountCertificates(ctx, k, msg)
		case *types.MsgDeleteAccountCertificate:
			return handlerMsgDeleteAccountCertificate(ctx, k, msg)
		case *types.MsgDeleteAccount:
			return handlerMsgDeleteAccount(ctx, k, msg)
		case *types.MsgReplaceAccountTargets:
			return handlerMsgReplaceAccountTargets(ctx, k, msg)
		case *types.MsgTransferAccount:
			return handlerMsgTransferAccount(ctx, k, msg)
		case *types.MsgSetAccountMetadata:
			return handlerMsgSetAccountMetadata(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unregonized request: %T", msg))
		}
	}

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		/*
			TODO
			remove when cosmos sdk decides that you are allowed to panic on errors that should not happen
			instead of returning random internal errors that mean actually nothing to a developer without
			a stacktrace or at least the error string of the panic itself, and also substitute 'log' stdlib
			with cosmos sdk logger when they make clear how you can use it and how to set up env to achieve so
		*/
		defer func() {
			if r := recover(); r != nil {
				log.Printf("FATAL-PANIC while executing message: %#v\nReason: %v", msg, r)
				// and lets panic again to throw it back to cosmos sdk yikes.
				panic(r)
			}
		}()
		resp, err := f(ctx, msg)
		if err != nil {
			msg := fmt.Sprintf("tx rejected %T: %s", msg, err)
			k.Logger(ctx).With("module", types.ModuleName).Info(msg)
		}
		return resp, err
	}
}

func handlerMsgAddAccountCertificates(ctx sdk.Context, k keeper.Keeper, msg *types.MsgAddAccountCertificates) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := ctrl.NewDomainController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(ctrl.DomainMustExist, ctrl.DomainNotExpired); err != nil {
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
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
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
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// delete cert
	k.DeleteAccountCertificate(ctx, accountCtrl.Account(), *certIndex)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}

// handlerMsgDelete account deletes the account from the system
func handlerMsgDeleteAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccount) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := ctrl.NewDomainController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(ctrl.DomainMustExist); err != nil {
		return nil, err
	}
	// perform account checks
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(account.MustExist); err != nil {
		return nil, err
	}
	// perform action authorization checks
	if (domainCtrl.Validate(ctrl.DomainOwner(msg.Owner)) != nil) && (accountCtrl.Validate(account.Owner(msg.Owner)) != nil) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner: %s and domain admin %s can delete the account", accountCtrl.Account().Owner, domainCtrl.Domain().Admin)
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
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
		return nil, sdkerrors.Wrap(types.ErrInvalidBlockchainTarget, err.Error())
	}
	// do validity checks on domain
	domainCtrl := ctrl.NewDomainController(ctx, k, msg.Domain)
	err := domainCtrl.Validate(ctrl.DomainMustExist, ctrl.DomainSuperuser(true), ctrl.DomainNotExpired, ctrl.DomainOwner(msg.Owner))
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
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
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
	domainCtrl := ctrl.NewDomainController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(ctrl.DomainMustExist); err != nil {
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
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
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
		return nil, sdkerrors.Wrapf(types.ErrInvalidBlockchainTarget, err.Error())
	}
	// perform domain checks
	domainCtrl := ctrl.NewDomainController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(ctrl.DomainMustExist, ctrl.DomainNotExpired); err != nil {
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
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// replace targets replaces accounts targets
	k.ReplaceAccountTargets(ctx, accountCtrl.Account(), msg.NewTargets)
	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}

// handlerMsgSetAccountMetadata takes care of setting account metadata
func handlerMsgSetAccountMetadata(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSetAccountMetadata) (*sdk.Result, error) {
	// perform domain checks
	domainCtrl := ctrl.NewDomainController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(ctrl.DomainMustExist, ctrl.DomainNotExpired); err != nil {
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
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
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
	domainCtrl := ctrl.NewDomainController(ctx, k, msg.Domain)
	if err := domainCtrl.Validate(ctrl.DomainMustExist, ctrl.DomainNotExpired); err != nil {
		return nil, err
	}
	// check if account exists
	accountCtrl := account.NewController(ctx, k, msg.Domain, msg.Name)
	if err := accountCtrl.Validate(account.MustExist, account.NotExpired); err != nil {
		return nil, err
	}
	// check if domain has super user
	switch domainCtrl.Domain().HasSuperuser {
	// if it has a super user then only domain admin can transfer accounts
	case true:
		if domainCtrl.Validate(ctrl.DomainOwner(msg.Owner)) != nil {
			return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only domain admin %s is allowed to transfer accounts", domainCtrl.Domain().Admin)
		}
	// if it has not a super user then only account owner can transfer the account
	case false:
		if accountCtrl.Validate(account.Owner(msg.Owner)) != nil {
			return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to transfer the account", accountCtrl.Account().Owner)
		}
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// transfer account
	k.TransferAccount(ctx, accountCtrl.Account(), msg.NewOwner)
	// success, todo emit event?
	return &sdk.Result{}, nil
}

func handlerMsgDeleteDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteDomain) (*sdk.Result, error) {
	c := ctrl.NewDomainController(ctx, k, msg.Domain)
	err := c.Validate(ctrl.DomainMustExist, ctrl.DomainSuperuser(true))
	if err != nil {
		return nil, err
	}
	// if domain is not over grace period and signer is not the owner of the domain then the operation is not allowed
	if err := c.Validate(ctrl.DomainOwner(msg.Owner)); err != nil && !c.Condition(ctrl.DomainGracePeriodFinished) {
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
	c := ctrl.NewDomainController(ctx, k, msg.Name)
	err = c.Validate(ctrl.DomainMustNotExist, ctrl.DomainValidName)
	if err != nil {
		return nil, err
	}
	// set new domain
	domain := types.Domain{
		Name:         msg.Name,
		Admin:        msg.Admin,
		ValidUntil:   ctx.BlockTime().Add(k.ConfigurationKeeper.GetDomainRenewDuration(ctx)).Unix(),
		HasSuperuser: msg.HasSuperuser,
		AccountRenew: msg.AccountRenew,
		Broker:       msg.Broker,
	}
	// if domain has not a super user then set domain to 0 address
	if !domain.HasSuperuser {
		domain.Admin = iovns.ZeroAddress // TODO change with module address
	}
	// save domain
	k.CreateDomain(ctx, domain)
	// generate empty name account
	acc := types.Account{
		Domain:       msg.Name,
		Name:         "",
		Owner:        msg.Admin, // TODO this is not clear, why the domain admin is zero address while this is msg.Admin
		ValidUntil:   ctx.BlockTime().Add(domain.AccountRenew).Unix(),
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
	c := ctrl.NewDomainController(ctx, k, msg.Domain)
	err := c.Validate(ctrl.DomainMustExist)
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
	c := ctrl.NewDomainController(ctx, k, msg.Domain)
	err := c.Validate(
		ctrl.DomainMustExist,
		ctrl.DomainSuperuser(true),
		ctrl.DomainOwner(msg.Owner),
		ctrl.DomainNotExpired,
	)
	if err != nil {
		return nil, err
	}
	// get domain
	domain := c.Domain()
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// transfer domain and accounts ownership
	k.TransferDomain(ctx, msg.NewAdmin, domain)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}
