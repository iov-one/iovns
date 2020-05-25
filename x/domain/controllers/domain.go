package controllers

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// DomainControllerFunc is the function signature for domain validation functions
type DomainControllerFunc func(controller *Domain) error

// DomainControllerCond is the function signature for domain condition functions
type DomainControllerCond func(controller *Domain) bool

// Domain is the domain controller
type Domain struct {
	domainName string
	ctx        sdk.Context
	domain     *types.Domain
	conf       *configuration.Config
	k          keeper.Keeper
}

// NewDomainController is the constructor for Domain
// everything is processed sequentially, a wrong order of the sequence
// is forbidden, example: asserting domain expiration on a non existing
// domain causes a panic as it violates the condition scope of action.
func NewDomainController(ctx sdk.Context, k keeper.Keeper, domain string) *Domain {
	return &Domain{
		domainName: domain,
		ctx:        ctx,
		k:          k,
	}
}

// ---------------------- VALIDATION -----------------------------

// Validate validates a domain based on the provided checks
func (c *Domain) Validate(checks ...DomainControllerFunc) error {
	for _, check := range checks {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

// Condition asserts if the given condition is true
func (c *Domain) Condition(cond DomainControllerCond) bool {
	return cond(c)
}

// DomainExpired checks if the provided domain has expired or not
func DomainExpired(controller *Domain) bool {
	return controller.domainExpired()
}

// domainExpired is the condition that checks if a domain has expired or not
func (c *Domain) domainExpired() bool {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("conditions check not allowed on non existing domain")
	}
	expireTime := iovns.SecondsToTime(c.domain.ValidUntil)
	// if expire time is before block time means domain expired
	return expireTime.Before(c.ctx.BlockTime())
}

func DomainGracePeriodFinished(controller *Domain) bool {
	return controller.gracePeriodFinished()
}

// gracePeriodFinished is the condition that checks if given domain is above grace period or not
func (c *Domain) gracePeriodFinished() bool {
	// require configuration
	c.requireConfiguration()
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("condition check not allowed on non existing domain")
	}
	// get grace period and expiration time
	gracePeriod := c.conf.DomainGracePeriod
	expireTime := iovns.SecondsToTime(c.domain.ValidUntil)
	// check if expiration time + grace period duration is before current block time
	return expireTime.Add(gracePeriod).Before(c.ctx.BlockTime())
}

func DomainOwner(addr sdk.AccAddress) DomainControllerFunc {
	return func(controller *Domain) error {
		return controller.ownedBy(addr)
	}
}

// ownedBy makes sure the domain is owned by the provided address
func (c *Domain) ownedBy(addr sdk.AccAddress) error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// check if admin matches addr
	if c.domain.Admin.Equals(addr) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to perform an operation in a domain owned by %s", addr, c.domain.Admin)
}

func DomainNotExpired(controller *Domain) error {
	return controller.notExpired()
}

func (c *Domain) notExpired() error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// check if domain has expired
	expireTime := iovns.SecondsToTime(c.domain.ValidUntil)
	// if block time is before expiration, return nil
	if c.ctx.BlockTime().Before(expireTime) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrDomainExpired, "%s has expired", c.domainName)
}

// DomainSuperuser makes sure the domain superuser is set to the provided condition
func DomainSuperuser(condition bool) DomainControllerFunc {
	return func(controller *Domain) error {
		return controller.superuser(condition)
	}
}

// superuser checks if the domain matches the superuser condition
func (c *Domain) superuser(condition bool) error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// check if superuser matches condition
	if c.domain.HasSuperuser == condition {
		return nil
	}
	switch condition {
	case true:
		return sdkerrors.Wrap(types.ErrUnauthorized, "operation is not allowed in domains with a superuser")
	default:
		return sdkerrors.Wrap(types.ErrUnauthorized, "operation is not allowed in domains without a superuser")
	}
}

// DomainMustExist checks if the provided domain exists
func DomainMustExist(controller *Domain) error {
	return controller.mustExist()
}

// requireDomain tries to find the domain by name
// if it is not found then an error is returned
func (c *Domain) requireDomain() error {
	if c.domain != nil {
		return nil
	}
	domain, ok := c.k.GetDomain(c.ctx, c.domainName)
	if !ok {
		return sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", c.domainName)
	}
	c.domain = &domain
	return nil
}

// mustExist checks if a domain exists
func (c *Domain) mustExist() error {
	return c.requireDomain()
}

// DomainMustNotExist checks if the provided domain does not exist
func DomainMustNotExist(controller *Domain) error {
	return controller.mustNotExist()
}

// mustNotExist asserts that a domain does not exist
func (c *Domain) mustNotExist() error {
	err := c.requireDomain()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrDomainAlreadyExists, c.domainName)
	}
	return nil
}

// DomainValidName checks if the name of the domain is valid
func DomainValidName(controller *Domain) error {
	return controller.validName()
}

// validName checks if the name of the domain is valid
func (c *Domain) validName() error {
	// require configuration
	c.requireConfiguration()
	// get valid domain regexp
	validator := regexp.MustCompile(c.conf.ValidDomain)
	// assert domain name validity
	if !validator.MatchString(c.domainName) {
		return sdkerrors.Wrap(types.ErrInvalidDomainName, c.domainName)
	}
	// success
	return nil
}

// requireConfiguration updates the configuration
// if it is not already set, and caches it after
func (c *Domain) requireConfiguration() {
	if c.conf != nil {
		return
	}
	conf := c.k.ConfigurationKeeper.GetConfiguration(c.ctx)
	c.conf = &conf
}

// Domain returns a copy the domain, panics if the operation is done without
// doing validity checks on domain existence as it is not an allowed op
func (c *Domain) Domain() types.Domain {
	if c.domain == nil {
		panic("get domain without running existence checks is not allowed")
	}
	return *c.domain
}

// ---------------------- MUTATION -----------------------------
