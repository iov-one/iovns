package account

import (
	"bytes"
	"regexp"
	"time"

	"github.com/iov-one/iovns/x/domain/controllers/domain"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// ControllerFunc is the function signature used by account controllers
type ControllerFunc func(ctrl *Account) error

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

	domainCtrl *domain.Domain
}

// NewController is Account constructor
func NewController(ctx sdk.Context, k keeper.Keeper, domain, name string) *Account {
	return &Account{
		name:   name,
		domain: domain,
		ctx:    ctx,
		k:      k,
	}
}

// WithDomainController allows to specify a cached domain controller
func (a *Account) WithDomainController(dom *domain.Domain) *Account {
	a.domainCtrl = dom
	return a
}

// WithConfiguration allows to specify a cached config
func (a *Account) WithConfiguration(cfg configuration.Config) *Account {
	a.conf = &cfg
	return a
}

// WithAccount allows to specify a cached account
func (a *Account) WithAccount(acc types.Account) *Account {
	a.account = &acc
	a.domain = acc.Domain
	a.name = acc.Name
	return a
}

// requireDomain builds the domain controller after asserting domain existence
func (a *Account) requireDomain() error {
	if a.domainCtrl != nil {
		return nil
	}
	a.domainCtrl = domain.NewController(a.ctx, a.k, a.domain)
	return a.domainCtrl.Validate(domain.MustExist)
}

// MustExist asserts if an account exists in the state,
// returns an error if it does not.
func MustExist(ctrl *Account) error {
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

// MustNotExist asserts that an account does not exist
func MustNotExist(ctrl *Account) error {
	return ctrl.mustNotExist()
}

// mustNotExist is the unexported function executed by MustNotExist
func (a *Account) mustNotExist() error {
	err := a.requireAccount()
	if err != nil {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrAccountExists, "account %s already exists in domain %s", a.name, a.domain)
}

// ValidAccountName asserts that an account has a valid name based
// on the account regexp  saved on the configuration module
func ValidName(ctrl *Account) error {
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

// validName is the unexported function used by ValidAccountName
func (a *Account) validName() error {
	a.requireConfiguration()
	if !regexp.MustCompile(a.conf.ValidAccountName).MatchString(a.name) {
		return sdkerrors.Wrapf(types.ErrInvalidAccountName, "invalid name: %s", a.name)
	}
	return nil
}

// NotExpired asserts that the account has
// not expired compared to the current block time
func NotExpired(ctrl *Account) error {
	return ctrl.notExpired()
}

// notExpired is the unexported function used by NotExpired
func (a *Account) notExpired() error {
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	if err := a.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	switch a.domainCtrl.Domain().Type {
	// if domain is closed type then skip the expiration validation checks
	case types.ClosedDomain:
		return nil
	}
	// check if account has expired
	expireTime := iovns.SecondsToTime(a.account.ValidUntil)
	if !expireTime.Before(a.ctx.BlockTime()) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrAccountExpired, "account %s in domain %s has expired", a.name, a.domain)
}

func MaxRenewNotExceeded(ctrl *Account) error {
	return ctrl.maxRenewNotExceeded()
}

func (a *Account) maxRenewNotExceeded() error {
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()

	// do calculations
	newValidUntil := iovns.SecondsToTime(a.account.ValidUntil).Add(a.conf.DomainRenewalPeriod) // set new expected valid until
	maximumValidUntil := a.ctx.BlockTime().Add(a.conf.AccountRenewalPeriod * time.Duration(a.conf.AccountRenewalCountMax))
	// check if new valid until is after maximum allowed
	if newValidUntil.After(maximumValidUntil) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "unable to renew account %s in domain %s, maximum domain renewal has exceeded: %s", a.account.Name, a.domain, maximumValidUntil)
	}

	// if it has expired return error
	return nil
}

// Owner asserts the account is owned by the provided address
func Owner(addr sdk.AccAddress) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.ownedBy(addr)
	}
}

// ownedBy is the unexported function used by Owner
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

// CertificateExists asserts that the provided certificate
// exists and if it does the index is saved in the provided pointer
// if certIndex pointer is nil the certificate index will not be saved
func CertificateExists(cert []byte, certIndex *int) ControllerFunc {
	return func(ctrl *Account) error {
		err := ctrl.certNotExist(cert, certIndex)
		if err == nil {
			return sdkerrors.Wrapf(types.ErrCertificateDoesNotExist, "%x", cert)
		}
		return nil
	}
}

// CertificateNotExist asserts the provided certificate
// does not exist in the account already
func CertificateNotExist(cert []byte) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.certNotExist(cert, nil)
	}
}

// certNotExist is the unexported function used by CertificateNotExist
// and CertificateExists, it saves the index of the found certificate
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

func CertificateSizeNotExceeded(cert []byte) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.certSizeNotExceeded(cert)
	}
}

func (a *Account) certSizeNotExceeded(newCert []byte) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint64(len(newCert)) > a.conf.CertificateSizeMax {
		return sdkerrors.Wrapf(types.ErrCertificateSizeExceeded, "max certificate size %d exceeded", a.conf.CertificateSizeMax)
	}
	return nil
}

func CertificateLimitNotExceeded(ctrl *Account) error {
	return ctrl.certLimitNotExceeded()
}

func (a *Account) certLimitNotExceeded() error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint32(len(a.account.Certificates)) >= a.conf.CertificateCountMax {
		return sdkerrors.Wrapf(types.ErrCertificateLimitReached, "max certificate limit %d reached, cannot add more", a.conf.CertificateCountMax)
	}
	return nil
}

// Validate verifies the account against the order of provided controllers
func (a *Account) Validate(checks ...ControllerFunc) error {
	for _, check := range checks {
		if err := check(a); err != nil {
			return err
		}
	}
	return nil
}

// DeletableBy checks if the account can be deleted by the provided address
func DeletableBy(addr sdk.AccAddress) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.deletableBy(addr)
	}
}

func (a *Account) deletableBy(addr sdk.AccAddress) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check on a non existing domain is not allowed")
	}
	// get cached domain
	d := a.domainCtrl.Domain()
	if err := a.requireAccount(); err != nil {
		panic("validation check on a non existing account is not allowed")
	}
	switch d.Type {
	case types.ClosedDomain:
		if err := a.domainCtrl.Validate(domain.Admin(addr), domain.NotExpired); err != nil {
			return err
		}
	case types.OpenDomain:
		if a.gracePeriodFinished() != nil {
			if a.ownedBy(addr) != nil {
				return sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to delete the account before grace period", a.account.Owner)
			}
		}
	}
	return nil
}

// ValidTargets verifies that the provided targets are valid for the account
func ValidTargets(targets []types.BlockchainAddress) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.validTargets(targets)
	}
}

// validTargets validates different blockchain targets address and ID
func (a *Account) validTargets(targets []types.BlockchainAddress) error {
	a.requireConfiguration()
	validBlockchainID := regexp.MustCompile(a.conf.ValidBlockchainID)
	validBlockchainAddress := regexp.MustCompile(a.conf.ValidBlockchainAddress)
	// create blockchain targets set to identify duplicates
	sets := make(map[string]struct{}, len(targets))
	// iterate over targets to check their validity
	for _, target := range targets {
		// check if blockchain ID was already specified
		if _, ok := sets[target.ID]; ok {
			return sdkerrors.Wrapf(types.ErrInvalidBlockchainTarget, "duplicate blockchain ID %s", target.ID)
		}
		sets[target.ID] = struct{}{}
		// is blockchain id valid?
		if !validBlockchainID.MatchString(target.ID) {
			return sdkerrors.Wrapf(types.ErrInvalidBlockchainTarget, "%s is not a valid blockchain ID", target.ID)
		}
		// is blockchain address valid?
		if !validBlockchainAddress.MatchString(target.Address) {
			return sdkerrors.Wrapf(types.ErrInvalidBlockchainTarget, "%s is not a valid blockchain address", target.Address)
		}
	}
	// success
	return nil
}

// TransferableBy checks if the account can be transferred by the provided address
func TransferableBy(addr sdk.AccAddress) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.transferableBy(addr)
	}
}

func (a *Account) transferableBy(addr sdk.AccAddress) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check not allowed on a non existing domain")
	}
	// check if domain has super user
	switch a.domainCtrl.Domain().Type {
	// if it has a super user then only domain admin can transfer accounts
	case types.ClosedDomain:
		if a.domainCtrl.Validate(domain.Admin(addr)) != nil {
			return sdkerrors.Wrapf(types.ErrUnauthorized, "only domain admin %s is allowed to transfer accounts", a.domainCtrl.Domain().Admin)
		}
	// if it has not a super user then only account owner can transfer the account
	case types.OpenDomain:
		if a.ownedBy(addr) != nil {
			return sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to transfer the account", a.account.Owner)
		}
	}
	return nil
}

// ResettableBy checks if the account attributes resettable by the provided address
func ResettableBy(addr sdk.AccAddress, reset bool) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.resettableBy(addr, reset)
	}
}

func (a *Account) resettableBy(addr sdk.AccAddress, reset bool) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check not allowed on a non existing domain")
	}
	d := a.domainCtrl.Domain()
	switch d.Type {
	case types.OpenDomain:
		if reset {
			if d.Admin.Equals(addr) {
				return sdkerrors.Wrapf(types.ErrUnauthorized, "domain admin is not authorized to reset account contents on open domains")
			}
		}
	case types.ClosedDomain:
	}
	return nil
}

func GracePeriodFinished(controller *Account) error {
	return controller.gracePeriodFinished()
}

// gracePeriodFinished is the condition that checks if given account's grace period has finished
func (a *Account) gracePeriodFinished() error {
	// require configuration
	a.requireConfiguration()
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("condition check not allowed on non existing account ")
	}
	// get grace period and expiration time
	gracePeriod := a.conf.AccountGracePeriod
	expireTime := iovns.SecondsToTime(a.account.ValidUntil)
	if a.ctx.BlockTime().After(expireTime.Add(gracePeriod)) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrAccountGracePeriodNotFinished, "account %s grace period has not finished", a.account.Name)
}

// ResettableBy checks if the account attributes resettable by the provided address
func BlockchainTargetLimitNotExceeded(targets []types.BlockchainAddress) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.blockchainTargetLimitNotExceeded(targets)
	}
}

func (a *Account) blockchainTargetLimitNotExceeded(targets []types.BlockchainAddress) error {
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint32(len(targets)) > a.conf.BlockchainTargetMax {
		return sdkerrors.Wrapf(types.ErrBlockchainTargetLimitExceeded, "blockchain target limit: %d", a.conf.BlockchainTargetMax)
	}
	return nil
}

func MetadataSizeNotExceeded(metadata string) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.metadataSizeNotExceeded(metadata)
	}
}

func (a *Account) metadataSizeNotExceeded(metadata string) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint64(len(metadata)) > a.conf.MetadataSizeMax {
		return sdkerrors.Wrapf(types.ErrMetadataSizeExceeded, "max metadata size %d exceeded", a.conf.MetadataSizeMax)
	}
	return nil
}

// RegistrableBy asserts that an account can be registered by the provided address
func RegistrableBy(addr sdk.AccAddress) ControllerFunc {
	return func(ctrl *Account) error {
		return ctrl.registrableBy(addr)
	}
}

func (a *Account) registrableBy(addr sdk.AccAddress) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// check domain type
	switch a.domainCtrl.Domain().Type {
	// if domain is closed then the registerer must be domain owner
	case types.ClosedDomain:
		return a.domainCtrl.Validate(domain.Admin(addr))
	default:
		return nil
	}
}

// Account returns the cached account, if the account existence
// was not asserted before, it panics.
func (a *Account) Account() types.Account {
	if err := a.requireAccount(); err != nil {
		panic("getting an account is not allowed before existence checks")
	}
	return *a.account
}
