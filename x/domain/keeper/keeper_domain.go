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

// CreateDomain creates the domain inside the KVStore with its name as key
func (k Keeper) CreateDomain(ctx sdk.Context, domain types.Domain) {
	// map domain to owner
	k.mapDomainToOwner(ctx, domain)
	// set domain
	k.SetDomain(ctx, domain)
}

// SetDomain updates or creates a new domain in the store
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
	domain, exists := k.GetDomain(ctx, domainName)
	if !exists {
		return
	}
	// delete domain
	domainStore := ctx.KVStore(k.domainStoreKey)
	domainStore.Delete([]byte(domainName))
	// delete accounts
	var accountKeys [][]byte
	k.GetAccountsInDomain(ctx, domainName, func(key []byte) bool {
		accountKeys = append(accountKeys, key)
		return true
	})
	// delete keys in domain
	for _, accountKey := range accountKeys {
		k.DeleteAccount(ctx, domainName, accountKeyToString(accountKey))
	}
	// unmap domain to owner
	k.unmapDomainToOwner(ctx, domain)
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

	// delete accounts
	var domainAccountKeys [][]byte
	k.GetAccountsInDomain(ctx, domainName, func(key []byte) bool {
		domainAccountKeys = append(domainAccountKeys, key)
		return true
	})
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

// TransferDomain transfers aliceAddr domain
func (k Keeper) TransferDomain(ctx sdk.Context, newOwner sdk.AccAddress, domain types.Domain) {
	// unmap domain owner
	k.unmapDomainToOwner(ctx, domain)
	// update domain owner
	domain.Admin = newOwner
	// update domain in kvstore
	k.SetDomain(ctx, domain)
	// get account keys related to the domain

	// delete accounts
	var accountKeys [][]byte
	k.GetAccountsInDomain(ctx, domain.Name, func(key []byte) bool {
		accountKeys = append(accountKeys, key)
		return true
	})
	// iterate over accounts
	for _, accountKey := range accountKeys {
		// skip if account key is empty account name
		if string(accountKey) == iovns.EmptyAccountName {
			continue
		}
		// get account
		account, _ := k.GetAccount(ctx, domain.Name, accountKeyToString(accountKey))
		// transfer it
		k.TransferAccount(ctx, account, newOwner)
	}
	// map domain to new owner
	k.mapDomainToOwner(ctx, domain)
}
