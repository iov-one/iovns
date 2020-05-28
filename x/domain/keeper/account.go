package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
	"time"
)

// contains all the functions to interact with the account store

// GetAccount finds an account based on its key name, if not found it will return
// a zeroed account and false.
func (k Keeper) GetAccount(ctx sdk.Context, domainName, accountName string) (account types.Account, exists bool) {
	// get domain prefix key
	store := accountsInDomainStore(ctx.KVStore(k.storeKey), domainName)
	// get account key
	accountKey := getAccountKey(accountName)
	// get account
	accountBytes := store.Get(accountKey)
	if accountBytes == nil {
		return
	}
	// key exists
	exists = true
	k.cdc.MustUnmarshalBinaryBare(accountBytes, &account)
	return
}

// CreateAccount creates an account
func (k Keeper) CreateAccount(ctx sdk.Context, account types.Account) {
	// create account
	k.SetAccount(ctx, account)
	// map account to owner
	err := k.mapAccountToOwner(ctx, account)
	if err != nil {
		panic(fmt.Errorf("indexing error (%#v): %w", account, err))
	}
	// map targets to account
	err = k.mapTargetToAccount(ctx, account, account.Targets...)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", account, err))
	}
}

// SetAccount upserts account data
func (k Keeper) SetAccount(ctx sdk.Context, account types.Account) {
	// get prefixed store
	store := accountsInDomainStore(ctx.KVStore(k.storeKey), account.Domain)
	// get account key
	accountKey := getAccountKey(account.Name)
	// set store
	store.Set(accountKey, k.cdc.MustMarshalBinaryBare(account))
}

// DeleteAccount deletes an account based on it full account name -> domain + iovns.Separator + account
func (k Keeper) DeleteAccount(ctx sdk.Context, domainName, accountName string) {
	// we need to retrieve account in order to unmap the account from the index; TODO can we avoid this?
	account, _ := k.GetAccount(ctx, domainName, accountName)
	store := accountsInDomainStore(ctx.KVStore(k.storeKey), domainName)
	// get account key
	accountKey := getAccountKey(account.Name)
	store.Delete(accountKey)
	// unmap account to owner
	err := k.unmapAccountToOwner(ctx, account)
	if err != nil {
		panic(fmt.Errorf("indexing error (%#v): %w", account, err))
	}
	// unmap targets to account
	err = k.unmapTargetToAccount(ctx, account, account.Targets...)
	if err != nil {
		panic(fmt.Errorf("indexing error (%#v): %w", account, err))
	}
}

// GetAccountsInDomain provides all the account keys related to the given domain name
func (k Keeper) GetAccountsInDomain(ctx sdk.Context, domainName string, do func(key []byte) bool) {
	// get store
	store := accountsInDomainStore(ctx.KVStore(k.storeKey), domainName)
	// create iterator
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		continueIterating := do(iterator.Key())
		if !continueIterating {
			return
		}
	}
	// return keys
	return
}

// TransferAccount transfers the account to aliceAddr new owner after resetting certificates and targets
func (k Keeper) TransferAccount(ctx sdk.Context, account types.Account, newOwner sdk.AccAddress) {
	// unmap account to owner
	err := k.unmapAccountToOwner(ctx, account)
	if err != nil {
		panic(fmt.Errorf("indexing error (%#v): %w", account, err))
	}
	// unmap account targets
	err = k.unmapTargetToAccount(ctx, account, account.Targets...)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", account, err))
	}
	// update account
	account.Owner = newOwner   // transfer owner
	account.Certificates = nil // remove certs
	account.Targets = nil      // remove targets
	// save account
	k.SetAccount(ctx, account)
	// map account to new owner
	err = k.mapAccountToOwner(ctx, account)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", account, err))
	}
	// map accounts new targets
	err = k.mapTargetToAccount(ctx, account, account.Targets...)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", account, err))
	}

}

// AddAccountCertificate adds aliceAddr new certificate to the account
func (k Keeper) AddAccountCertificate(ctx sdk.Context, account types.Account, newCert []byte) {
	// if not add it to accounts certs
	account.Certificates = append(account.Certificates, newCert)
	// update account
	k.SetAccount(ctx, account)
}

// DeleteAccountCertificate deletes aliceAddr certificate at given index, it will panic if the index is wrong
func (k Keeper) DeleteAccountCertificate(ctx sdk.Context, account types.Account, certificateIndex int) {
	// remove cert
	account.Certificates = append(account.Certificates[:certificateIndex], account.Certificates[certificateIndex+1:]...)
	// update account
	k.SetAccount(ctx, account)
}

// RenewAccount updates an account expiration time
func (k Keeper) RenewAccount(ctx sdk.Context, account types.Account, accountRenew time.Duration) {
	// update account time
	account.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(account.ValidUntil).Add(accountRenew),
	)
	// update account in kv store
	k.SetAccount(ctx, account)
}

// ReplaceAccountTargets updates an account targets
func (k Keeper) ReplaceAccountTargets(ctx sdk.Context, account types.Account, targets []types.BlockchainAddress) {
	// unmap old targets
	err := k.unmapTargetToAccount(ctx, account, account.Targets...)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", account, err))
	}
	// replace targets
	account.Targets = targets
	// update account
	k.SetAccount(ctx, account)
	// map new targets
	err = k.mapTargetToAccount(ctx, account, targets...)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", account, err))
	}
}

// IterateAllAccounts returns all the accounts inside the store
func (k Keeper) IterateAllAccounts(ctx sdk.Context) []types.Account {
	store := accountStore(ctx.KVStore(k.storeKey))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	var accounts []types.Account
	for ; iterator.Valid(); iterator.Next() {
		var a types.Account
		accountBytes := store.Get(iterator.Key())
		k.cdc.MustUnmarshalBinaryBare(accountBytes, &a)
		accounts = append(accounts, a)
	}
	return accounts
}
