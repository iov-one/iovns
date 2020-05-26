package domain

import (
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// ControllerFunc is the function signature for domain validation functions
type ControllerFunc func(controller *Domain) error

// ControllerCond is the function signature for domain condition functions
type ControllerCond func(controller *Domain) bool

// Domain is the domain controller
type Domain struct {
	domainName string
	ctx        sdk.Context
	domain     *types.Domain
	conf       *configuration.Config
	k          keeper.Keeper
}

// NewController is the constructor for Domain
// everything is processed sequentially, a wrong order of the sequence
// is forbidden, example: asserting domain expiration on a non existing
// domain causes a panic as it violates the condition scope of action.
func NewController(ctx sdk.Context, k keeper.Keeper, domain string) *Domain {
	return &Domain{
		domainName: domain,
		ctx:        ctx,
		k:          k,
	}
}

// ---------------------- VALIDATION -----------------------------

// Validate validates a domain based on the provided checks
func (c *Domain) Validate(checks ...ControllerFunc) error {
	for _, check := range checks {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

// Condition asserts if the given condition is true
func (c *Domain) Condition(cond ControllerFunc) bool {
	return cond(c) == nil
}

// Expired checks if the provided domain has expired or not
func Expired(controller *Domain) error {
	return controller.expired()
}

// expired returns nil if domain expired, otherwise ErrDomainNotExpired
func (c *Domain) expired() error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("conditions check not allowed on non existing domain")
	}
	expireTime := iovns.SecondsToTime(c.domain.ValidUntil)
	// if expire time is before block time means domain expired
	if expireTime.Before(c.ctx.BlockTime()) {
		return nil
	}

	return sdkerrors.Wrapf(types.ErrDomainNotExpired, "domain %s has not expired", c.domain.Name)
}

func GracePeriodFinished(controller *Domain) error {
	return controller.gracePeriodFinished()
}

// gracePeriodFinished is the condition that checks if given domain's grace period has finished
func (c *Domain) gracePeriodFinished() error {
	// require configuration
	c.requireConfiguration()
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("condition check not allowed on non existing domain")
	}
	// get grace period and expiration time
	gracePeriod := c.conf.DomainGracePeriod
	expireTime := iovns.SecondsToTime(c.domain.ValidUntil)
	if c.ctx.BlockTime().After(expireTime.Add(gracePeriod)) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrGracePeriodNotFinished, "domain %s grace period has not finished", c.domain.Name)
}

func Owner(addr sdk.AccAddress) ControllerFunc {
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

func NotExpired(controller *Domain) error {
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

// Superuser makes sure the domain superuser is set to the provided condition
func Type(Type types.DomainType) ControllerFunc {
	return func(controller *Domain) error {
		return controller.dType(Type)
	}
}

func (c *Domain) dType(Type types.DomainType) error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	if c.domain.Type != Type {
		return sdkerrors.Wrapf(types.ErrInvalidDomainType, "invalid domain type %s, expected %s", c.domain.Type, Type)
	}
	return nil
}

// MustExist checks if the provided domain exists
func MustExist(controller *Domain) error {
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

// MustNotExist checks if the provided domain does not exist
func MustNotExist(controller *Domain) error {
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

// ValidName checks if the name of the domain is valid
func ValidName(controller *Domain) error {
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
