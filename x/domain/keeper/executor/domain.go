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
		d.domains.Update(d)
		return
	}
	// get configuration
	renewDuration := d.k.ConfigurationKeeper.GetDomainRenewDuration(d.ctx)
	// update prefixedStore valid until
	d.domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(d.domain.ValidUntil).Add(renewDuration), // time(prefixedStore.ValidUntil) + renew duration
	)
	// set prefixedStore
	d.domains.Update(d.domain)
}

// Delete deletes a prefixedStore from the kvstore
func (d *Domain) Delete() {
	if d.domain == nil {
		panic("cannot execute delete state change on non present prefixedStore")
	}
	d.domains.Delete(d.domain)
	d.accounts.Delete(&types.Account{Domain: d.domain.Name, Name: ""})
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
		d.domains.Update(d.domain)
		emptyAccount := new(types.Account)
		ok := d.accounts.Read((&types.Account{Domain: d.domain.Name, Name: ""}).PrimaryKey(), emptyAccount)
		if !ok {
			panic("empty account not found")
		}

		// transfer accounts of the prefixedStore based on the transfer flag
		switch flag {
		// reset none is simply skipped as empty account is already transferred during prefixedStore transfer
		case types.ResetNone:
		// transfer flush, deletes all domains accounts except the empty one since it was transferred in the first step
		case types.TransferFlush:
			var accountKeys []crud.PrimaryKey
			d.accounts.IterateIndex(crud.SecondaryKey{
				Key:         []byte(d.domain.Name),
				StorePrefix: nil,
			}, func(key crud.PrimaryKey) bool {
				accountKeys = append(accountKeys, key)
				return true
			})
			for _, key := range accountKeys {
				d.accounts.Delete(key)
			}
		// transfer owned transfers only accounts owned by the old owner
		case types.TransferOwned:
			// TODO change when crud supports multiple indexes
			var accountsInDomain []crud.PrimaryKey
			d.accounts.IterateIndex(crud.SecondaryKey{
				Key:         []byte(d.domain.Name),
				StorePrefix: nil,
			}, func(key crud.PrimaryKey) bool {
				accountsInDomain = append(accountsInDomain, key)
				return true
			})
			// iterate over accounts
			for _, key := range accountsInDomain {
				acc := new(types.Account)
				ok := d.accounts.Read(key, acc)
				if !ok {
					panic("missing account data which was indexed")
				}
				// skip accounts which are not owned by old owner
				if !acc.Owner.Equals(oldOwner) {
					continue
				}
				// change owner and update
				acc.Owner = newOwner
				d.accounts.Update(acc)
			}
		}
	}
}

// Create creates a new prefixedStore
func (d *Domain) Create() {
	if d.domain == nil {
		panic("cannot create non specified domain")
	}
	d.domains.Update(d.domain)
}
