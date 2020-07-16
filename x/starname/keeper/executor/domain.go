package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/pkg/crud"
	"github.com/iov-one/iovns/tutils"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

// Domain defines the prefixedStore keeper executor
type Domain struct {
	domain   *types.Domain
	ctx      sdk.Context
	domains  crud.Store
	accounts crud.Store
	k        keeper.Keeper
}

// NewDomain returns is prefixedStore's constructor
func NewDomain(ctx sdk.Context, k keeper.Keeper, dom types.Domain) *Domain {
	return &Domain{
		k:        k,
		ctx:      ctx,
		domains:  k.DomainStore(ctx),
		accounts: k.AccountStore(ctx),
		domain:   &dom,
	}
}

// Renew renews a prefixedStore based on the configuration
func (d *Domain) Renew(accValidUntil ...int64) {
	if d.domain == nil {
		panic("cannot execute renew state change on non present prefixedStore")
	}
	// if account valid until is specified then the renew is coming from accounts
	if len(accValidUntil) != 0 {
		d.domain.ValidUntil = accValidUntil[0]
		d.domains.Update(d.domain.PrimaryKey(), d.domain)
		return
	}
	// get configuration
	renewDuration := d.k.ConfigurationKeeper.GetDomainRenewDuration(d.ctx)
	// update prefixedStore valid until
	d.domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(d.domain.ValidUntil).Add(renewDuration), // time(prefixedStore.ValidUntil) + renew duration
	)
	// set prefixedStore
	d.domains.Update(d.domain.PrimaryKey(), d.domain)
}

// Delete deletes a domain from the kvstore
func (d *Domain) Delete() {
	if d.domain == nil {
		panic("cannot execute delete state change on non present prefixedStore")
	}
	filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name})
	for filter.Next() {
		filter.Delete()
	}
	d.domains.Delete(d.domain.PrimaryKey(), d.domain)
}

// Transferrer returns a prefixedStore transfer function based on the transfer flag
func (d *Domain) Transfer(flag types.TransferFlag, newOwner sdk.AccAddress) func() {
	if d.domain == nil {
		panic("cannot execute transfer state on non defined prefixedStore")
	}
	return func() {
		// transfer domain
		var oldOwner = d.domain.Admin // cache it for future uses
		d.domain.Admin = newOwner
		d.domains.Update(d.domain.PrimaryKey(), d.domain)
		filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name, Name: tutils.StrPtr(types.EmptyAccountName)}) // delete empty account
		filter.Next()
		filter.Delete()
		// transfer accounts of the prefixedStore based on the transfer flag
		switch flag {
		// reset none is simply skipped as empty account is already transferred during prefixedStore transfer
		case types.ResetNone:
		// transfer flush, deletes all domains accounts except the empty one since it was transferred in the first step
		case types.TransferFlush:
			filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name})
			for filter.Next() {
				filter.Delete()
			}
		// transfer owned transfers only accounts owned by the old owner
		case types.TransferOwned:
			filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name, Owner: oldOwner})
			for filter.Next() {
				acc := new(types.Account)
				filter.Read(acc)
				acc.Owner = newOwner
				filter.Update(acc)
			}
		}
	}
}

// Create creates a new prefixedStore
func (d *Domain) Create() {
	if d.domain == nil {
		panic("cannot create non specified domain")
	}
	d.domains.Create(d.domain)
	emptyAccount := &types.Account{
		Domain:       d.domain.Name,
		Name:         tutils.StrPtr(types.EmptyAccountName),
		Owner:        d.domain.Admin,
		ValidUntil:   0,
		Resources:    nil,
		Certificates: nil,
		Broker:       nil,
		MetadataURI:  "",
	}
	d.accounts.Create(emptyAccount)
}
