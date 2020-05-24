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

// AccountControllerFunc is the function signature used by account controllers
type AccountControllerFunc func(ctrl *Account) error

// AccountControllerCond is the function signature used by account condition controllers
type AccountControllerCond func(ctrl *Account) bool

// Account is an account controller, it caches information
// in order to avoid useless query to state to get the same
// information. Order of execution of controllers matters
// if the correct order is not followed the controller will
// panic because of bad operation flow.
// Errors returned are wrapped sdk.Error types.
type Account struct {
	name, domain string
	account      *types.Account
	conf         *configuration.Config

	ctx sdk.Context
	k   keeper.Keeper
}

// NewAccountController is Account constructor
func NewAccountController(ctx sdk.Context, k keeper.Keeper, domain, name string) *Account {
	return &Account{
		name:   name,
		domain: domain,
		ctx:    ctx,
		k:      k,
	}
}

// AccountMustExist asserts if an account exists in the state,
// returns an error if it does not.
func AccountMustExist(ctrl *Account) error {
	return ctrl.mustExist()
}

// requireAccount finds the accounts and caches it, so future
// queries will always use the same account first found account
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

// mustExist makes sure an account exist
func (a *Account) mustExist() error {
	return a.requireAccount()
}

// AccountMustNotExist asserts that an account does not exist
func AccountMustNotExist(ctrl *Account) error {
	return ctrl.mustNotExist()
}

// mustNotExist is the unexported function executed by AccountMustNotExist
func (a *Account) mustNotExist() error {
	err := a.requireAccount()
	if err != nil {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrAccountExists, "account %s already exists in domain %s", a.name, a.domain)
}

// AccountValidName asserts that an account has a vaid name based
// on the account regexp  saved on the configuration module
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

// validName is the unexported function used by AccountValidName
func (a *Account) validName() error {
	a.requireConfiguration()
	if !regexp.MustCompile(a.conf.ValidName).MatchString(a.name) {
		return sdkerrors.Wrapf(types.ErrInvalidAccountName, "invalid name: %s", a.name)
	}
	return nil
}

// AccountNotExpired asserts that the account has
// not expired compared to the current block time
func AccountNotExpired(ctrl *Account) error {
	return ctrl.notExpired()
}

// notExpired is the unexported function used by AccountNotExpired
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

// AccountOwner asserts the account is owned by the provided address
func AccountOwner(addr sdk.AccAddress) AccountControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.ownedBy(addr)
	}
}

// ownedBy is the unexported function used by AccountOwner
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

// AccountCertificateExists asserts that the provided certificate
// exists and if it does the index is saved in the provided pointer
// if certIndex pointer is nil the certificate index will not be saved
func AccountCertificateExists(cert []byte, certIndex *int) AccountControllerFunc {
	return func(ctrl *Account) error {
		err := ctrl.certNotExist(cert, certIndex)
		if err == nil {
			return sdkerrors.Wrapf(types.ErrCertificateDoesNotExist, "%x", cert)
		}
		return nil
	}
}

// AccountCertificateNotExist asserts the provided certificate
// does not exist in the account already
func AccountCertificateNotExist(cert []byte) AccountControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.certNotExist(cert, nil)
	}
}

// certNotExist is the unexported function used by AccountCertificateNotExist
// and AccountCertificateExists, it saves the index of the found certificate
// in indexPointer if it is not nil
func (a *Account) certNotExist(newCert []byte, indexPointer *int) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	// check if certificate is already present in account
	for i, cert := range a.account.Certificates {
		if bytes.Equal(cert, newCert) {
			if indexPointer != nil {
				*indexPointer = i
			}
			return sdkerrors.Wrapf(types.ErrCertificateExists, "certificate is already present")
		}
	}
	return nil
}

// Validate verifies the account against the order of provided controllers
func (a *Account) Validate(checks ...AccountControllerFunc) error {
	for _, check := range checks {
		if err := check(a); err != nil {
			return err
		}
	}
	return nil
}

// Account returns the cached account, if the account existence
// was not asserted before, it panics.
func (a *Account) Account() types.Account {
	if a.account == nil {
		panic("getting an account is not allowed before existence checks")
	}
	return *a.account
}
