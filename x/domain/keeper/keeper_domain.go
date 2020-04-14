package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
)

// contains all the functions to interact with the domain store

// GetDomain returns the domain based on its name, if domain is not found ok will be false
func (k Keeper) GetDomain(ctx sdk.Context, domainName string) (domain types.Domain, ok bool) {
	store := ctx.KVStore(k.domainKey)
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
	store := ctx.KVStore(k.domainKey)
	store.Set([]byte(domain.Name), k.cdc.MustMarshalBinaryBare(domain))
}

// IterateAllDomains will return an iterator for all the domain keys
// present in the KVStore, it's callers duty to close the iterator.
func (k Keeper) IterateAllDomains(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.domainKey)
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
	domainStore := ctx.KVStore(k.domainKey)
	domainStore.Delete([]byte(domainName))
	// delete accounts, TODO do it efficiently with KVUtils
	accountStore := ctx.KVStore(k.accountKey)
	iterator := accountStore.Iterator(nil, nil)
	var accountKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		accountKeys = append(accountKeys, iterator.Key())
	}
	iterator.Close()
	// delete account keys
	for _, key := range accountKeys {
		// check if this account key belongs to the domain
		accountDomain, _ := iovns.SplitAccountKey(key)
		// if account domain does not match domain name
		// we want to delete then continue
		if accountDomain != domainName {
			continue
		}
		accountStore.Delete(key)
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
	accountStore := ctx.KVStore(k.accountKey)
	// get all account keys
	iterator := accountStore.Iterator(nil, nil)
	var domainAccountKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		// check if account key matches the domain
		key := iterator.Key()
		accountDomain, accountName := iovns.SplitAccountKey(key)
		// if key does not belong to domain skip
		if accountDomain != domainName {
			continue
		}
		// if accountName is empty account name then skip
		if accountName == "" {
			continue
		}
		// append
		domainAccountKeys = append(domainAccountKeys, iterator.Key())
	}
	iterator.Close()
	// now delete accounts
	for _, accountKey := range domainAccountKeys {
		accountStore.Delete(accountKey)
	}
	// success
	return
}
