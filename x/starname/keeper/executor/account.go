package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud/pkg/crud"
	"github.com/iov-one/iovns/pkg/utils"
	"github.com/iov-one/iovns/x/starname/keeper"
	"github.com/iov-one/iovns/x/starname/types"
)

func NewAccount(ctx sdk.Context, k keeper.Keeper, account types.Account) *Account {
	return &Account{
		store:   k.AccountStore(ctx),
		account: &account,
		ctx:     ctx,
		k:       k,
	}
}

// Account defines an account executor
type Account struct {
	store   crud.Store
	account *types.Account
	ctx     sdk.Context
	k       keeper.Keeper
}

func (a *Account) Transfer(newOwner sdk.AccAddress, reset bool) {
	if a.account == nil {
		panic("cannot transfer non specified account")
	}
	// apply account changes
	// update owner
	a.account.Owner = newOwner
	// if reset is required then clear the account
	if reset {
		a.account.Certificates = nil
		a.account.Resources = nil
		a.account.MetadataURI = ""
	}
	// apply changes
	a.store.Update(a.account)
}

func (a *Account) UpdateMetadata(newMetadata string) {
	if a.account == nil {
		panic("cannot update metadata on non specified account")
	}
	a.account.MetadataURI = newMetadata
	a.store.Update(a.account)
}

func (a *Account) ReplaceResources(newTargets []types.Resource) {
	if a.account == nil {
		panic("cannot replace targets on non specified account")
	}
	a.account.Resources = newTargets
	a.store.Update(a.account)
}

func (a *Account) Renew() {
	if a.account == nil {
		panic("cannot renew a non specified account")
	}
	renew := a.k.ConfigurationKeeper.GetConfiguration(a.ctx).AccountRenewalPeriod
	a.account.ValidUntil = utils.TimeToSeconds(
		utils.SecondsToTime(a.account.ValidUntil).Add(renew),
	)
	// update account in kv store
	a.store.Update(a.account)
}

func (a *Account) Create() {
	if a.account == nil {
		panic("cannot create a non specified account")
	}
	a.store.Create(a.account)
}

func (a *Account) Delete() {
	if a.account == nil {
		panic("cannot delete a non specified account")
	}
	a.store.Delete(a.account.PrimaryKey())
}

func (a *Account) DeleteCertificate(index int) {
	if a.account == nil {
		panic("cannot delete certificate on a non specified account")
	}
	a.account.Certificates = append(a.account.Certificates[:index], a.account.Certificates[index+1:]...)
	a.store.Update(a.account)
}

func (a *Account) AddCertificate(cert []byte) {
	if a.account == nil {
		panic("cannot add certificate on a non specified account")
	}
	a.account.Certificates = append(a.account.Certificates, cert)
	a.store.Update(a.account)
}

// State returns the current state of the account
func (a *Account) State() types.Account {
	if a.account == nil {
		panic("cannot get state of a non specified account")
	}
	return *a.account
}
