package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
)

// contains all the functions to interact with the domain store

// GetDomain returns the domain based on its name, if domain is not found ok will be false
func (k Keeper) GetDomain(ctx sdk.Context, domainName string) (domain types.Domain, ok bool) {
	store := ctx.KVStore(k.domainStoreKey)
	// get domain in form of bytes
	domainBytes := store.Get([]byte(domainName))
	// if nothing is returned, return nil
	if domainBytes == nil {
		return
	}
	// if domain exists then unmarshal
	k.cdc.MustUnmarshalBinaryBare(domainBytes, &domain)
	// success
	return domain, true
}

// SetDomain saves the domain inside the KVStore with its name as key
func (k Keeper) SetDomain(ctx sdk.Context, domain types.Domain) {
	store := ctx.KVStore(k.domainStoreKey)
	store.Set([]byte(domain.Name), k.cdc.MustMarshalBinaryBare(domain))
}

// IterateAllDomains will return an iterator for all the domain keys
// present in the KVStore, it's callers duty to close the iterator.
func (k Keeper) IterateAllDomains(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.domainStoreKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}

// DeleteDomain deletes the domain and the accounts in it
// this operation can only fail in case the domain does not exist
func (k Keeper) DeleteDomain(ctx sdk.Context, domainName string) (exists bool) {
	_, exists = k.GetDomain(ctx, domainName)
	if !exists {
		return
	}
	// delete domain
	domainStore := ctx.KVStore(k.domainStoreKey)
	domainStore.Delete([]byte(domainName))
	// delete accounts,
	accountKeys := k.GetAccountsInDomain(ctx, domainName)
	// delete keys in domain
	for _, accountKey := range accountKeys {
		k.DeleteAccount(ctx, domainName, accountKeyToString(accountKey))
	}
	// done
	return true
}

// FlushDomain removes all accounts, except the empty one, from the domain.
// returns true in case the domain exists and the operation has been done.
// returns false only in case the domain does not exist.
func (k Keeper) FlushDomain(ctx sdk.Context, domainName string) (exists bool) {
	_, exists = k.GetDomain(ctx, domainName)
	if !exists {
		return
	}
	// iterate accounts
	domainAccountKeys := k.GetAccountsInDomain(ctx, domainName)
	// now delete accounts
	for _, accountKey := range domainAccountKeys {
		// account key is empty then skip
		if string(accountKey) == "" {
			continue
		}
		// otherwise delete
		k.DeleteAccount(ctx, domainName, accountKeyToString(accountKey))
	}
	// success
	return
}

// TransferDomain transfers a domain
func (k Keeper) TransferDomain(ctx sdk.Context, newOwner sdk.AccAddress, domain types.Domain) {
	// update domain owner
	domain.Admin = newOwner
	// set domain
	k.SetDomain(ctx, domain)
	// get account keys related to the domain
	accountKeys := k.GetAccountsInDomain(ctx, domain.Name)
	// iterate over accounts
	for _, accountKey := range accountKeys {
		// skip if account key is empty account name
		if string(accountKey) == iovns.EmptyAccountName {
			continue
		}
		// get account;
		account, _ := k.GetAccount(ctx, domain.Name, accountKeyToString(accountKey))
		// update account
		account.Certificates = nil // delete certs
		account.Targets = nil      // delete targets
		account.Owner = newOwner   // update admin
		// save to kvstore
		k.SetAccount(ctx, account)
	}
}
