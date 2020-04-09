package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/domain/types"
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
		// check if this belongs to the domain
		accountDomain, accountName := iovnsd.SplitAccountKey(key)
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
