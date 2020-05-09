package keeper

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/x/domain/types"
)

// ownerToAccountPrefix is the prefix that matches owners to accounts
var ownerToAccountPrefix = []byte{0x04}

// ownerToAccountIndexSeparator is the separator used to map owner address + domain + account name
var ownerToAccountIndexSeparator = []byte(":")

// ownerToDomainPrefix is the prefix that matches owners to domains
var ownerToDomainPrefix = []byte{0x05}

// ownerToDomainIndexSeparator is the separator used to map owner address + domain
var ownerToDomainIndexSeparator = []byte(":")

var blockchainTargetsPrefix = []byte{0x06}

func blockchainTargetIndexedStore(store sdk.KVStore, target types.BlockchainAddress) index.Store {
	return index.NewIndexedStore(store, blockchainTargetsPrefix, target)
}

func (k Keeper) mapTargetToAccount(ctx sdk.Context, account types.Account, targets ...types.BlockchainAddress) {
	for _, target := range targets {
		// if targets are empty ignore
		if target.Address == "" && target.ID == "" {
			return
		}
		// otherwise map target to given account
		store := blockchainTargetIndexedStore(ctx.KVStore(k.storeKey), target)
		store.Set(account.Index())
	}
}

func (k Keeper) unmapTargetToAccount(ctx sdk.Context, account types.Account, targets ...types.BlockchainAddress) {
	for _, target := range targets {
		// if targets are empty then ignore the process
		if target.ID == "" && target.Address == "" {
			return
		}
		store := blockchainTargetIndexedStore(ctx.KVStore(k.storeKey), target)
		if ok := store.Delete(account.Index()); !ok {
			panic(fmt.Sprintf("index not found: target %#v, account: %#v", target, account))
		}
	}
}

func (k Keeper) iterateBlockchainTargetsAccounts(ctx sdk.Context, target types.BlockchainAddress, do func(key []byte) bool) {
	store := blockchainTargetIndexedStore(ctx.KVStore(k.storeKey), target)
	store.IterateKeys(do)
}

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

// getOwnerToAccountKey generates the unique key that maps owner to account
func getOwnerToAccountKey(owner sdk.AccAddress, domain string, account string) []byte {
	// get index bytes of addr
	addr := indexAddr(owner)
	// generate unique key
	return bytes.Join([][]byte{addr, []byte(domain), []byte(account)}, ownerToAccountIndexSeparator)
}

func indexAddr(addr sdk.AccAddress) []byte {
	x := addr.String()
	return []byte(x)
}

func accAddrFromIndex(indexedAddr []byte) sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(string(indexedAddr))
	if err != nil {
		panic(err)
	}
	return accAddr
}

// getOwnerToDomainKey generates the unique key that maps owner to domain
func getOwnerToDomainKey(owner sdk.AccAddress, domain string) []byte {
	addrBytes := indexAddr(owner)
	return bytes.Join([][]byte{addrBytes, []byte(domain)}, ownerToDomainIndexSeparator)
}

// splitOwnerToAccountKey takes an indexed owner to account key and splits it
// into owner address, domain name and account name
func splitOwnerToAccountKey(key []byte) (addr sdk.AccAddress, domain string, account string) {
	splitBytes := bytes.SplitN(key, ownerToAccountIndexSeparator, 3)
	if len(splitBytes) != 3 {
		panic(fmt.Sprintf("unexpected split length: %d", len(splitBytes)))
	}
	// convert back to their original types
	addr, domain, account = accAddrFromIndex(splitBytes[0]), string(splitBytes[1]), string(splitBytes[2])
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
	addr, domain = accAddrFromIndex(splitBytes[0]), string(splitBytes[1])
	return
}

func (k Keeper) unmapAccountToOwner(ctx sdk.Context, account types.Account) {
	// get store
	store := accountIndexStore(indexStore(ctx.KVStore(k.storeKey)))

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
	store := accountIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	key := getOwnerToAccountKey(account.Owner, account.Domain, account.Name)
	// check if key exists TODO remove panic
	if store.Has(key) {
		panic(fmt.Sprintf("existing store key: %s", key))
	}
	// set key
	store.Set(key, []byte{})
}

func (k Keeper) iterAccountToOwner(ctx sdk.Context, address sdk.AccAddress, do func(key []byte) bool) {
	// get store
	store := accountIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	// get iterator
	iterator := sdk.KVStorePrefixIterator(store, indexAddr(address))
	defer iterator.Close()
	// iterate keys
	for ; iterator.Valid(); iterator.Next() {
		// do action
		keepGoing := do(iterator.Key())
		// keep going?
		if !keepGoing {
			return
		}
	}
}

func (k Keeper) mapDomainToOwner(ctx sdk.Context, domain types.Domain) {
	// get store
	store := domainIndexStore(indexStore(ctx.KVStore(k.storeKey)))
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
	store := domainIndexStore(indexStore(ctx.KVStore(k.storeKey)))
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
func (k Keeper) iterDomainToOwner(ctx sdk.Context, address sdk.AccAddress, do func(key []byte) bool) {
	// get store
	store := domainIndexStore(indexStore(ctx.KVStore(k.storeKey)))
	// get iterator
	iterator := sdk.KVStorePrefixIterator(store, indexAddr(address))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if !do(iterator.Key()) {
			return
		}
	}
}
