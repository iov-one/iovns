package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/x/domain/types"
)

// ownerToAccountPrefix is the prefix that matches owners to accounts
var ownerToAccountPrefix = []byte{0x04}

// ownerToDomainPrefix is the prefix that matches owners to domains
var ownerToDomainPrefix = []byte{0x05}

// resourcesPrefix is the prefix used to index resources to account
var resourcesPrefix = []byte{0x06}

// ownerToDomainIndexStore returns the store that indexes all the domains owned by an sdk.AccAddress
func ownerToDomainIndexStore(kvstore sdk.KVStore, addr sdk.AccAddress) (index.Store, error) {
	// check if admin is provided
	if addr.Empty() {
		return index.Store{}, fmt.Errorf("cannot index empty address")
	}
	// get index store
	indexPrefixedStore := indexStore(kvstore)
	// get address to domain index
	return index.NewAddressIndex(indexPrefixedStore, ownerToDomainPrefix, addr)
}

// ownerToAccountIndexStore returns the indexed store that indexes all the accounts owned by an sdk.AccAddress
func ownerToAccountIndexStore(kvstore sdk.KVStore, addr sdk.AccAddress) (index.Store, error) {
	if addr.Empty() {
		return index.Store{}, fmt.Errorf("cannot index empty address")
	}
	// get index store
	indexPrefixedStore := indexStore(kvstore)
	// get address to account index
	return index.NewAddressIndex(indexPrefixedStore, ownerToAccountPrefix, addr)
}

// resourcesIndexStore returns the store used to index accounts resources
func resourcesIndexStore(store sdk.KVStore, resource types.Resource) (index.Store, error) {
	prefixedIndexStore := indexStore(store)
	return index.NewIndexedStore(prefixedIndexStore, resourcesPrefix, resource)
}

func (k Keeper) mapResourceToAccount(ctx sdk.Context, account types.Account, resources ...types.Resource) error {
	for _, resource := range resources {
		// if resources are empty ignore
		if resource.Resource == "" || resource.URI == "" {
			continue
		}
		// otherwise map resource to given account
		store, err := resourcesIndexStore(ctx.KVStore(k.storeKey), resource)
		if err != nil {
			return err
		}
		err = store.Set(account)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) unmapResourcesToAccount(ctx sdk.Context, account types.Account, resources ...types.Resource) error {
	for _, resource := range resources {
		// if resources are empty then ignore the process
		if resource.URI == "" || resource.Resource == "" {
			continue
		}
		store, err := resourcesIndexStore(ctx.KVStore(k.storeKey), resource)
		if err != nil {
			return err
		}
		if err = store.Delete(account); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) iterateResourceAccounts(ctx sdk.Context, resource types.Resource, do func(key []byte) bool) error {
	store, err := resourcesIndexStore(ctx.KVStore(k.storeKey), resource)
	if err != nil {
		return err
	}
	store.IterateAllKeys(do)
	return nil
}

func (k Keeper) unmapAccountToOwner(ctx sdk.Context, account types.Account) error {
	// get store
	store, err := ownerToAccountIndexStore(ctx.KVStore(k.storeKey), account.Owner)
	if err != nil {
		return err
	}
	// delete account
	err = store.Delete(account)
	if err != nil {
		return err
	}
	return nil
}

// mapAccountToOwner maps accounts to an owner
func (k Keeper) mapAccountToOwner(ctx sdk.Context, account types.Account) error {
	// get store
	store, err := ownerToAccountIndexStore(ctx.KVStore(k.storeKey), account.Owner)
	if err != nil {
		return err
	}
	// set key
	err = store.Set(account)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) iterAccountToOwner(ctx sdk.Context, address sdk.AccAddress, do func(key []byte) bool) error {
	// get store
	store, err := ownerToAccountIndexStore(ctx.KVStore(k.storeKey), address)
	if err != nil {
		return err
	}
	store.IterateAllKeys(do)
	return nil
}

func (k Keeper) mapDomainToOwner(ctx sdk.Context, domain types.Domain) error {
	// get index store
	store, err := ownerToDomainIndexStore(ctx.KVStore(k.storeKey), domain.Admin)
	if err != nil {
		return err
	}
	// set key
	err = store.Set(domain)
	if err != nil {
		return err
	}
	// success
	return nil
}

func (k Keeper) unmapDomainToOwner(ctx sdk.Context, domain types.Domain) error {
	// get store
	store, err := ownerToDomainIndexStore(ctx.KVStore(k.storeKey), domain.Admin)
	if err != nil {
		return err
	}
	// delete domain
	err = store.Delete(domain)
	if err != nil {
		return err
	}
	// success
	return nil
}

// iterDomainToOwner iterates over all the domains owned by address
// and returns the unique keys
func (k Keeper) iterDomainToOwner(ctx sdk.Context, address sdk.AccAddress, do func(key []byte) bool) error {
	// get store
	store, err := ownerToDomainIndexStore(ctx.KVStore(k.storeKey), address)
	if err != nil {
		return err
	}
	store.IterateAllKeys(do)
	return nil
}
