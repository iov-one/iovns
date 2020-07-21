package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/pkg/crud"
	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

// Domain defines the domain keeper executor
type Domain struct {
	domain   *types.Domain
	ctx      sdk.Context
	domains  crud.Store
	accounts crud.Store
	k        keeper.Keeper
}

// NewDomain returns is domain's constructor
func NewDomain(ctx sdk.Context, k keeper.Keeper, dom types.Domain) *Domain {
	return &Domain{
		k:        k,
		ctx:      ctx,
		domains:  k.DomainStore(ctx),
		accounts: k.AccountStore(ctx),
		domain:   &dom,
	}
}

// Renew renews a domain based on the configuration
func (d *Domain) Renew(accValidUntil ...int64) {
	if d.domain == nil {
		panic("cannot execute renew state change on non present domain")
	}
	// if account valid until is specified then the renew is coming from accounts
	if len(accValidUntil) != 0 {
		d.domain.ValidUntil = accValidUntil[0]
		d.domains.Update(d.domain.PrimaryKey(), d.domain)
		return
	}
	// get configuration
	renewDuration := d.k.ConfigurationKeeper.GetDomainRenewDuration(d.ctx)
	// update domain valid until
	d.domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(d.domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// set domain
	d.domains.Update(d.domain.PrimaryKey(), d.domain)
}

// Delete deletes a domain from the kvstore
func (d *Domain) Delete() {
	if d.domain == nil {
		panic("cannot execute delete state change on non present domain")
	}
	filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name})
	for ; filter.Valid(); filter.Next() {
		filter.Delete()
	}
	d.domains.Delete(d.domain.PrimaryKey(), d.domain)
}

// Transfer transfers a domain given a flag and an owner
func (d *Domain) Transfer(flag types.TransferFlag, newOwner sdk.AccAddress) {
	if d.domain == nil {
		panic("cannot execute transfer state on non defined domain")
	}

	// transfer domain
	var oldOwner = d.domain.Admin // cache it for future uses
	d.domain.Admin = newOwner
	d.domains.Update(d.domain.PrimaryKey(), d.domain)
	// transfer empty account
	filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name, Name: utils.StrPtr(types.EmptyAccountName)})
	emptyAccount := new(types.Account)
	filter.Read(emptyAccount)
	ac := NewAccount(d.ctx, d.k, *emptyAccount)
	ac.Transfer(newOwner, true)
	// transfer accounts of the domain based on the transfer flag
	switch flag {
	// reset none is simply skipped as empty account is already transferred during domain transfer
	case types.TransferResetNone:
		return
	// transfer flush, deletes all domains accounts except the empty one since it was transferred in the first step
	case types.TransferFlush:
		filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name})
		for ; filter.Valid(); filter.Next() {
			filter.Delete()
		}
	// transfer owned transfers only accounts owned by the old owner
	case types.TransferOwned:
		filter := d.accounts.Filter(&types.Account{Domain: d.domain.Name, Owner: oldOwner})
		for ; filter.Valid(); filter.Next() {
			acc := new(types.Account)
			filter.Read(acc)
			acc.Owner = newOwner
			// do account reset
			acc.Resources = nil
			acc.Certificates = nil
			acc.MetadataURI = ""
			// update account
			filter.Update(acc)
		}
	}
}

// Create creates a new domain
func (d *Domain) Create() {
	if d.domain == nil {
		panic("cannot create non specified domain")
	}
	d.domains.Create(d.domain)
	emptyAccount := &types.Account{
		Domain:       d.domain.Name,
		Name:         utils.StrPtr(types.EmptyAccountName),
		Owner:        d.domain.Admin,
		ValidUntil:   d.domain.ValidUntil, // is this right per spec?
		Resources:    nil,
		Certificates: nil,
		Broker:       nil,
		MetadataURI:  "",
	}
	d.accounts.Create(emptyAccount)
}
