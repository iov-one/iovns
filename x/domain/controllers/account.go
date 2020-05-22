package controllers

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
	"regexp"
)

type AccountControllerFunc func(ctrl *Account) error
type AccountControllerCond func(ctrl *Account) bool

type Account struct {
	name, domain string
	account      *types.Account
	conf         *configuration.Config

	ctx sdk.Context
	k   keeper.Keeper
}

func NewAccountController(ctx sdk.Context, k keeper.Keeper, domain, name string) *Account {
	return &Account{
		name:   name,
		domain: domain,
		ctx:    ctx,
		k:      k,
	}
}

func AccountMustExist(ctrl *Account) error {
	return ctrl.mustExist()
}

func (a *Account) requireAccount() error {
	if a.account != nil {
		return nil
	}
	account, ok := a.k.GetAccount(a.ctx, a.domain, a.name)
	if !ok {
		return sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "%s was not found in domain %s", a.name, a.domain)
	}
	a.account = &account
	return nil
}

func (a *Account) mustExist() error {
	return a.requireAccount()
}

func AccountMustNotExist(ctrl *Account) error {
	return ctrl.mustNotExist()
}

func (a *Account) mustNotExist() error {
	err := a.requireAccount()
	if err != nil {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrAccountExists, "account %s already exists in domain %s", a.name, a.domain)
}

func AccountValidName(ctrl *Account) error {
	return ctrl.validName()
}

// requireConfiguration updates the configuration
// if it is not already set, and caches it after
func (a *Account) requireConfiguration() {
	if a.conf != nil {
		return
	}
	conf := a.k.ConfigurationKeeper.GetConfiguration(a.ctx)
	a.conf = &conf
}

func (a *Account) validName() error {
	a.requireConfiguration()
	if !regexp.MustCompile(a.conf.ValidName).MatchString(a.name) {
		return sdkerrors.Wrapf(types.ErrInvalidAccountName, "invalid name: %s", a.name)
	}
	return nil
}

func AccountNotExpired(ctrl *Account) error {
	return ctrl.notExpired()
}

func (a *Account) notExpired() error {
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	// check if account has expired
	expireTime := iovns.SecondsToTime(a.account.ValidUntil)
	if !expireTime.Before(a.ctx.BlockTime()) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrAccountExpired, "account %s in domain %s has expired", a.name, a.domain)
}

func AccountOwner(addr sdk.AccAddress) AccountControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.ownedBy(addr)
	}
}

func (a *Account) ownedBy(addr sdk.AccAddress) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	// check if admin matches at least one address
	if a.account.Owner.Equals(addr) {
		return nil
	}

	return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to perform operation in the account owned by %s", addr, a.account.Owner)
}

func AccountCertificateNotExist(cert []byte) AccountControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.certNotExist(cert)
	}
}

func (a *Account) certNotExist(newCert []byte) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	// check if certificate is already present in account
	for _, cert := range a.account.Certificates {
		if bytes.Equal(cert, newCert) {
			return sdkerrors.Wrapf(types.ErrCertificateExists, "certificate is already present")
		}
	}
	return nil
}

func (a *Account) Validate(checks ...AccountControllerFunc) error {
	for _, check := range checks {
		if err := check(a); err != nil {
			return err
		}
	}
	return nil
}

func (a *Account) GetAccount() types.Account {
	if a.account == nil {
		panic("getting an account is not allowed before existence checks")
	}
	return *a.account
}
