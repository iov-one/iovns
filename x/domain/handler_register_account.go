package domain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/configuration"
	"github.com/iov-one/iovnsd/x/domain/keeper"
	"github.com/iov-one/iovnsd/x/domain/types"
	"regexp"
	"time"
)

// handleMsgRegisterAccount registers the domain
func handleMsgRegisterAccount(ctx sdk.Context, k keeper.Keeper, msg types.MsgRegisterAccount) (*sdk.Result, error) {
	// verify request
	// get config
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	// validate blockchain targets
	if err := validateBlockchainTargets(msg.Targets, conf); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidBlockchainTarget, err.Error())
	}
	// validate account name
	if !regexp.MustCompile(conf.ValidName).MatchString(msg.Name) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidName, "account name %s is invalid", msg.Name)
	}
	// check if domain name exists
	domain, ok := k.GetDomain(ctx, msg.Domain)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "domain %s does not exist", msg.Domain)
	}
	// check if domain has super user that owner equals to the domain admin
	if domain.HasSuperuser && !domain.Admin.Equals(msg.Owner) {
		return nil, sdkerrors.Wrapf(types.ErrUnauthorized, "address %s is not authorized to register an account in a domain with superuser", msg.Owner)
	}
	// check if domain is still valid
	if ctx.BlockTime().After(iovnsd.SecondsToTime(domain.ValidUntil)) {
		return nil, sdkerrors.Wrap(types.ErrDomainExpired, "account registration is not allowed")
	}
	// check account does not exist already
	if _, ok := k.GetAccount(ctx, iovnsd.GetAccountKey(msg.Domain, msg.Name)); ok {
		return nil, sdkerrors.Wrapf(types.ErrAccountExists, "account: %s exists for domain %s", msg.Name, msg.Domain)
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
	k.SetAccount(ctx, account)
	// success; TODO can we emit events?
	return &sdk.Result{}, nil
}

// validateBlockchainTargets validates different blockchain targets address and ID
func validateBlockchainTargets(targets []iovnsd.BlockchainAddress, conf configuration.Config) error {
	validBlockchainID := regexp.MustCompile(conf.ValidBlockchainID)
	validBlockchainAddress := regexp.MustCompile(conf.ValidBlockchainAddress)
	// iterate over targets to check their validity
	for _, target := range targets {
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
