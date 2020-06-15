package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// Domain defines the domain keeper executor
type Domain struct {
	domain *types.Domain
	ctx    sdk.Context
	k      keeper.Keeper
}

// NewDomain returns is domain's constructor
func NewDomain(ctx sdk.Context, k keeper.Keeper, dom types.Domain) *Domain {
	return &Domain{
		ctx:    ctx,
		k:      k,
		domain: &dom,
	}
}

// Renew renews a domain based on the configuration
func (d *Domain) Renew() {
	if d.domain == nil {
		panic("cannot execute renew state change on non present domain")
	}
	// get configuration
	renewDuration := d.k.ConfigurationKeeper.GetDomainRenewDuration(d.ctx)
	// update domain valid until
	d.domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(d.domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// set domain
	d.k.SetDomain(d.ctx, *d.domain)
}

// Delete deletes a domain from the kvstore
func (d *Domain) Delete() {
	if d.domain == nil {
		panic("cannot execute delete state change on non present domain")
	}
	d.k.DeleteDomain(d.ctx, d.domain.Name)
}

// Transferrer returns a domain transfer function based on the transfer flag
func (d *Domain) Transfer(flag types.TransferFlag, newOwner sdk.AccAddress) func() {
	if d.domain == nil {
		panic("cannot execute transfer state on non defined domain")
	}
	return func() {
		// transfer domain
		d.k.TransferDomainOwnership(d.ctx, *d.domain, newOwner)
		// transfer accounts of the domain based on the transfer flag
		switch flag {
		// reset none is simply skipped as empty account is already transferred during domain transfer
		case types.ResetNone:
		// transfer flush, deletes all domain accounts except the empty one
		case types.TransferFlush:
			d.k.FlushDomain(d.ctx, *d.domain)
		// transfer owned transfers only accounts owned by the old owner
		case types.TransferOwned:
			d.k.TransferDomainAccountsOwnedByAddr(d.ctx, *d.domain, d.domain.Admin, newOwner)
		}
	}
}

// Create creates a new domain
func (d *Domain) Create() {
	if d.domain == nil {
		panic("cannot create non specified domain")
	}
	d.k.CreateDomain(d.ctx, *d.domain)
}
