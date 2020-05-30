package domain

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func handlerFn(ctx sdk.Context, k Keeper, msg sdk.Msg) (*sdk.Result, error) {
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

// NewHandler builds the tx requests handler for the domain module
func NewHandler(k Keeper) sdk.Handler {
	f := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case *types.MsgWithFeePayer:
			res, err := handlerFn(ctx, k, msg)
			if err != nil {
				if err := k.CollectFees(ctx, msg.Msg, msg.FeePayer); err != nil {
					return res, err
				}
			}
			return res, err
		default:
			res, err := handlerFn(ctx, k, msg)
			if err != nil {
				if err := k.CollectFees(ctx, msg, sdk.AccAddress{}); err != nil {
					return res, err
				}
			}
			return res, err
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
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if current time is after domain validity time
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "domain %s has expired", msg.Domain)
	}
	// get account
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if current time is after account validity time
	if ctx.BlockTime().After(iovns.SecondsToTime(account.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if signer is account owner
	if !msg.Owner.Equals(account.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s cannot add certificates to account owned by %s", msg.Owner, account.Owner)
	}
	// check if certificate is already present in account
	for _, cert := range account.Certificates {
		if bytes.Equal(cert, msg.NewCertificate) {
			return nil, sdkerrors.Wrapf(types.ErrCertificateExists, "certificate is already present")
		}
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}
	// add certificate
	k.AddAccountCertificate(ctx, account, msg.NewCertificate)
	// success; TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgDeleteAccountCertificate(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccountCertificate) (*sdk.Result, error) {
	// get account
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if signer is account owner
	if !msg.Owner.Equals(account.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s cannot delete certificates from account owned by %s", msg.Owner, account.Owner)
	}
	// check if certificate exists
	var found bool
	var certIndex int
	// iterate over certs
	for i, cert := range account.Certificates {
		// if found
		if bytes.Equal(cert, msg.DeleteCertificate) {
			found = true  // set found to true
			certIndex = i // save index of cert for removal
			break
		}
	}
	// check if found
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCertificateDoesNotExist, "not found")
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// delete cert
	k.DeleteAccountCertificate(ctx, account, certIndex)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}

// handlerMsgDelete account deletes the account from the system
func handlerMsgDeleteAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccount) (*sdk.Result, error) {
	// check if domain exists
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if account exists
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if msg.Owner is either domain owner or account owner
	if !domain.Admin.Equals(msg.Owner) && !account.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner: %s and domain admin %s can delete the account", account.Owner, domain.Admin)
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
	// validate account name
	if !regexp.MustCompile(conf.ValidName).MatchString(msg.Name) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAccountName, "invalid name: %s", msg.Name)
	}
	// check if domain name exists
	domain, ok := k.GetDomain(ctx, msg.Domain)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if domain has super user that owner equals to the domain admin
	if domain.HasSuperuser && !domain.Admin.Equals(msg.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "address %s is not authorized to register an account in a domain with superuser", msg.Owner)
	}
	// check if domain is still valid
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrap(types.ErrDomainExpired, "account registration is not allowed")
	}
	// check account does not exist already
	if _, ok := k.GetAccount(ctx, msg.Domain, msg.Name); ok {
		return nil, sdkerrors.Wrapf(types.ErrAccountExists, "account: %s exists for domain %s", msg.Name, msg.Domain)
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
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
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// get account
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// renew account
	k.UpdateAccountValidity(ctx, account, domain.AccountRenew)
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
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// see if domain still valid
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "domain %s has expired", msg.Domain)
	}
	// get account
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if expired
	if ctx.BlockTime().After(iovns.SecondsToTime(account.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if account owner matches request signer
	if !msg.Owner.Equals(account.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "account %s is not authorized to perform actions on account owned by %s", msg.Owner, account.Owner)
	}
	// collect fees
	err = k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// replace targets replaces accounts targets
	k.ReplaceAccountTargets(ctx, account, msg.NewTargets)
	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}

// handlerMsgSetAccountMetadata takes care of setting account metadata
func handlerMsgSetAccountMetadata(ctx sdk.Context, k keeper.Keeper, msg *types.MsgSetAccountMetadata) (*sdk.Result, error) {
	// get domain
	domain, ok := k.GetDomain(ctx, msg.Domain)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if domain expired
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "domain %s has expired", domain.Name)
	}
	// get account
	account, ok := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if account expired
	if ctx.BlockTime().After(iovns.SecondsToTime(account.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if signer is account owner
	if !account.Owner.Equals(msg.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "not allowed to change account metadata uri, invalid owner: %s", msg.Owner)
	}
	// update account
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
	// check if domain exists
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if domain has expired expired
	if iovns.SecondsToTime(domain.ValidUntil).Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "account transfer is not allowed for expired domains, expire date: %s", iovns.SecondsToTime(domain.ValidUntil))
	}
	// check if account exists
	account, exists := k.GetAccount(ctx, msg.Domain, msg.Name)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found in domain %s: %s", msg.Domain, msg.Name)
	}
	// check if account has expired
	if iovns.SecondsToTime(account.ValidUntil).Before(ctx.BlockTime()) {
		return nil, sdkerrors.Wrapf(types.ErrAccountExpired, "account %s has expired", msg.Name)
	}
	// check if domain has super user
	switch domain.HasSuperuser {
	// if it has a super user then only domain admin can transfer accounts
	case true:
		if !msg.Owner.Equals(domain.Admin) {
			return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only domain admin %s is allowed to transfer accounts", domain.Admin)
		}
	// if it has not a super user then only account owner can transfer the account
	case false:
		if !msg.Owner.Equals(account.Owner) {
			return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to transfer the account", account.Owner)
		}
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// transfer account
	k.TransferAccount(ctx, account, msg.NewOwner)
	// success, todo emit event?
	return &sdk.Result{}, nil
}

func handlerMsgDeleteDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteDomain) (*sdk.Result, error) {
	// check if domain exists
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if domain has super user
	if !domain.HasSuperuser {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "can not delete domain with no superuser")
	}
	// check if domain admin matches msg owner and if the domain has not expired (plus the grace period)
	gracePeriod := k.ConfigurationKeeper.GetDomainGracePeriod(ctx)

	// check if domain has expired and we are not over grace period
	if !ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil).Add(gracePeriod)) {
		if !domain.Admin.Equals(msg.Owner) {
			return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "address %s is not allowed to delete the domain owned by: %s", msg.Owner, domain.Admin)
		}
	}
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
	// check if domain exists
	if _, ok := k.GetDomain(ctx, msg.Name); ok {
		return nil, sdkerrors.Wrap(types.ErrDomainAlreadyExists, msg.Name)
	}
	// if domain does not exist then check if we can register it
	// check if name is valid based on the configuration saved in the state
	if !regexp.MustCompile(k.ConfigurationKeeper.GetValidDomainRegexp(ctx)).MatchString(msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrInvalidDomainName, msg.Name)
	}
	// if domain has not a super user then admin must be configuration owner
	if !msg.HasSuperuser && !k.ConfigurationKeeper.IsOwner(ctx, msg.Admin) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to register a domain without a superuser", msg.Admin)
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
	// check if domain exists
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found %s", msg.Domain)
	}
	// get configuration
	renewDuration := k.ConfigurationKeeper.GetDomainRenewDuration(ctx)
	// update domain valid until
	domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// update domain
	k.SetDomain(ctx, domain)
	// success TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgTransferDomain(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferDomain) (*sdk.Result, error) {
	// get domain
	domain, exists := k.GetDomain(ctx, msg.Domain)
	if !exists {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", msg.Domain)
	}
	// check if has superuser
	if !domain.HasSuperuser {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "domain %s without superuser cannot be transferred", msg.Domain)
	}
	// check if signer is domain owner
	if !msg.Owner.Equals(domain.Admin) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to transfer domain owned by %s", msg.Owner, domain.Admin)
	}
	// check if domain is valid
	if ctx.BlockTime().After(iovns.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrapf(types.ErrDomainExpired, "%s has expired", msg.Domain)
	}
	// collect fees
	err := k.CollectFees(ctx, msg, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// transfer account ownership
	k.TransferDomain(ctx, msg.NewAdmin, domain)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}
