package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/keeper"
	"github.com/iov-one/iovns/x/domain/types"
)

func NewAccount(ctx sdk.Context, k keeper.Keeper, account types.Account) *Account {
	return &Account{
		account: &account,
		ctx:     ctx,
		k:       k,
	}
}

// Account defines an account executor
type Account struct {
	account *types.Account
	ctx     sdk.Context
	k       keeper.Keeper
}

func (a *Account) Transfer(newOwner sdk.AccAddress, reset bool) {
	if a.account == nil {
		panic("cannot transfer non specified account")
	}
	a.k.TransferAccountWithReset(a.ctx, *a.account, newOwner, reset)
}

func (a *Account) UpdateMetadata(newMetadata string) {
	if a.account == nil {
		panic("cannot update metadata on non specified account")
	}
	a.k.UpdateMetadataAccount(a.ctx, *a.account, newMetadata)
}

func (a *Account) ReplaceTargets(newTargets []types.BlockchainAddress) {
	if a.account == nil {
		panic("cannot replace targets on non specified account")
	}
	a.k.ReplaceAccountTargets(a.ctx, *a.account, newTargets)
}

func (a *Account) Renew() {
	if a.account == nil {
		panic("cannot renew a non specified account")
	}
	renew := a.k.ConfigurationKeeper.GetConfiguration(a.ctx).AccountRenewalPeriod
	a.account.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(a.account.ValidUntil).Add(renew),
	)
	// update account in kv store
	a.k.SetAccount(a.ctx, *a.account)
}

func (a *Account) Create() {
	if a.account == nil {
		panic("cannot create a non specified account")
	}
	a.k.CreateAccount(a.ctx, *a.account)
}

func (a *Account) Delete() {
	if a.account == nil {
		panic("cannot delete a non specified account")
	}
	a.k.DeleteAccount(a.ctx, a.account.Domain, a.account.Name)
}

func (a *Account) DeleteCertificate(index int) {
	if a.account == nil {
		panic("cannot delete certificate on a non specified account")
	}
	a.k.DeleteAccountCertificate(a.ctx, *a.account, index)
}

func (a *Account) AddCertificate(cert []byte) {
	if a.account == nil {
		panic("cannot add certificate on a non specified account")
	}
	a.k.AddAccountCertificate(a.ctx, *a.account, cert)
}

// State returns the current state of the account
func (a *Account) State() types.Account {
	if a.account == nil {
		panic("cannot get state of a non specified account")
	}
	return *a.account
}
