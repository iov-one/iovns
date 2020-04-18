package keeper

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
)

// ownerToAccountPrefix is the prefix that matches owners to accounts
var ownerToAccountPrefix = []byte("owneracc")
var ownerToAccountIndexSeparator = []byte(":")

var ownerToDomainPrefix = []byte("ownerdom")
var ownerToDomainIndexSeparator = []byte(":")

// domainIndexStore returns the kvstore space that maps
// owner to domain keys
func domainIndexStore(store sdk.KVStore) sdk.KVStore {
	return prefix.NewStore(store, ownerToDomainPrefix)
}

// accountIndexStore returns the kvstore space that maps
// owner to accounts
func accountIndexStore(store sdk.KVStore) sdk.KVStore {
	return prefix.NewStore(store, ownerToAccountPrefix)
}

// getOwnerToDomainKey generates the unique key that maps owner to domain
func getOwnerToDomainKey(owner sdk.AccAddress, domain string) []byte {
	return bytes.Join([][]byte{owner.Bytes(), []byte(domain)}, ownerToDomainPrefix)
}

// getOwnerToAccountKey generates the unique key that maps owner to account
func getOwnerToAccountKey(owner sdk.AccAddress, domain string, account string) []byte {
	return bytes.Join([][]byte{owner.Bytes(), []byte(domain), []byte(account)}, ownerToAccountIndexSeparator)
}

// splitOwnerToAccountKey takes an indexed owner to account key and splits it
// into owner address, domain name and account name
func splitOwnerToAccountKey(key []byte) (addr sdk.AccAddress, domain string, account string) {
	splitBytes := bytes.SplitN(key, ownerToAccountIndexSeparator, 3)
	if len(splitBytes) != 3 {
		panic(fmt.Sprintf("unexpected split length: %d", len(splitBytes)))
	}
	// convert back to their original types
	addr, domain, account = splitBytes[0], string(splitBytes[1]), string(splitBytes[2])
	return
}

// splitOwnerToDomainKey takes an indexed owner to domain key
// and splits it into owner address and domain name
func splitOwnerToDomainKey(key []byte) (addr sdk.AccAddress, domain string) {
	splitBytes := bytes.SplitN(key, ownerToDomainIndexSeparator, 2)
	if len(splitBytes) != 2 {
		panic(fmt.Sprintf("expected split lenght: %d", len(splitBytes)))
	}
	// convert back to their original types
	addr, domain = splitBytes[0], string(splitBytes[1])
	return
}

func (k Keeper) unmapAccountToOwner(ctx sdk.Context, account types.Account) {
	// get store
	store := accountIndexStore(ctx.KVStore(k.indexStoreKey))

	// check if key exists TODO remove panic
	key := getOwnerToAccountKey(account.Owner, account.Domain, account.Name)
	if !store.Has(key) {
		panic(fmt.Sprintf("missing store key: %s", key))
	}
	// delete key
	store.Delete(key)
}

// mapAccountToOwner maps accounts to an owner
func (k Keeper) mapAccountToOwner(ctx sdk.Context, account types.Account) {
	// get store
	store := accountIndexStore(ctx.KVStore(k.indexStoreKey))
	key := getOwnerToAccountKey(account.Owner, account.Domain, account.Name)
	// check if key exists TODO remove panic
	if store.Has(key) {
		panic(fmt.Sprintf("existing store key: %s", key))
	}
	// set key
	store.Set(key, []byte{})
}

func (k Keeper) iterAccountToOwner(ctx sdk.Context, address sdk.AccAddress) [][]byte {
	// get store
	store := accountIndexStore(ctx.KVStore(k.indexStoreKey))
	// get iterator
	iterator := sdk.KVStorePrefixIterator(store, address.Bytes())
	defer iterator.Close()

	var accountKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		accountKeys = append(accountKeys, iterator.Key())
	}
	return accountKeys
}

func (k Keeper) mapDomainToOwner(ctx sdk.Context, domain types.Domain) {
	// get store
	store := domainIndexStore(ctx.KVStore(k.indexStoreKey))
	// get unique key
	key := getOwnerToDomainKey(domain.Admin, domain.Name)
	// check if key exists TODO remove panic
	if store.Has(key) {
		panic(fmt.Sprintf("existing store key: %s", key))
	}
	// set key
	store.Set(key, []byte{})
}

func (k Keeper) unmapDomainToOwner(ctx sdk.Context, domain types.Domain) {
	// get store
	store := domainIndexStore(ctx.KVStore(k.indexStoreKey))
	// check if key exists TODO remove panic
	key := getOwnerToDomainKey(domain.Admin, domain.Name)
	if !store.Has(key) {
		panic(fmt.Sprintf("missing store key: %s", key))
	}
	// delete key
	store.Delete(key)
}

// iterDomainToOwner iterates over all the domains owned by address
// and returns the unique keys
func (k Keeper) iterDomainToOwner(ctx sdk.Context, address sdk.AccAddress) [][]byte {
	// get store
	store := domainIndexStore(ctx.KVStore(k.indexStoreKey))
	// get iterator
	iterator := sdk.KVStorePrefixIterator(store, address.Bytes())
	defer iterator.Close()

	var domainKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		domainKeys = append(domainKeys, iterator.Key())
	}
	return domainKeys
}
