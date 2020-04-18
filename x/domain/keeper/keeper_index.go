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

// accountIndexStore returns the index that maps owners to accounts
func accountIndexStore(store sdk.KVStore) sdk.KVStore {
	return prefix.NewStore(store, ownerToAccountPrefix)
}

// getOwnerToAccountKey
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
	// delete key
	store.Set(key, []byte{})
}

func (k Keeper) iterAccountToOwner(ctx sdk.Context, address sdk.AccAddress) {
	// get store
	store := accountIndexStore(ctx.KVStore(k.indexStoreKey))
	// get iterator
	iterator := sdk.KVStorePrefixIterator(store, address.Bytes())
	defer iterator.Close()

	var accountKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		accountKeys = append(accountKeys, iterator.Key())
	}
}
