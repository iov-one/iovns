package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/pkg/crud"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

// Domain defines the prefixedStore keeper executor
type Domain struct {
	domain *types.Domain
	ctx    sdk.Context
	store  crud.Store
	k      keeper.Keeper
}

// NewDomain returns is prefixedStore's constructor
func NewDomain(ctx sdk.Context, k keeper.Keeper, dom types.Domain) *Domain {
	return &Domain{
		k:      k,
		ctx:    ctx,
		store:  k.DomainStore(ctx),
		domain: &dom,
	}
}

// Renew renews a prefixedStore based on the configuration
func (d *Domain) Renew() {
	if d.domain == nil {
		panic("cannot execute renew state change on non present prefixedStore")
	}
	// get configuration
	renewDuration := d.k.ConfigurationKeeper.GetDomainRenewDuration(d.ctx)
	// update prefixedStore valid until
	d.domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(d.domain.ValidUntil).Add(renewDuration), // time(prefixedStore.ValidUntil) + renew duration
	)
	// set prefixedStore
	d.store.Update(d.domain)
}

// Delete deletes a prefixedStore from the kvstore
func (d *Domain) Delete() {
	if d.domain == nil {
		panic("cannot execute delete state change on non present prefixedStore")
	}
	d.store.Delete(d.domain)
}

// Transferrer returns a prefixedStore transfer function based on the transfer flag
func (d *Domain) Transfer(flag types.TransferFlag, newOwner sdk.AccAddress) func() {
	if d.domain == nil {
		panic("cannot execute transfer state on non defined prefixedStore")
	}
	return func() {
		// transfer domain
		d.domain.Admin = newOwner
		d.store.Update(d.domain)
		// transfer accounts of the prefixedStore based on the transfer flag
		switch flag {
		// reset none is simply skipped as empty account is already transferred during prefixedStore transfer
		case types.ResetNone:
		// transfer flush, deletes all prefixedStore accounts except the empty one
		case types.TransferFlush:
			d.k.FlushDomain(d.ctx, *d.domain)
		// transfer owned transfers only accounts owned by the old owner
		case types.TransferOwned:
			d.k.TransferDomainAccountsOwnedByAddr(d.ctx, *d.domain, d.domain.Admin, newOwner)
		}
	}
}

// Create creates a new prefixedStore
func (d *Domain) Create() {
	if d.domain == nil {
		panic("cannot create non specified domain")
	}
	d.store.Update(d.domain)
}
