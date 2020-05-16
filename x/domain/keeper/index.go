package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/x/domain/types"
)

// ownerToAccountPrefix is the prefix that matches owners to accounts
var ownerToAccountPrefix = []byte{0x04}

// ownerToDomainPrefix is the prefix that matches owners to domains
var ownerToDomainPrefix = []byte{0x05}

// blockchainTargetsPrefix is the prefix used to index targets to account
var blockchainTargetsPrefix = []byte{0x06}

// certificatesPrefix is the prefix used to index certificates to accounts
var certificatesPrefix = []byte{0x07}

// blockchainTargetIndexedStore returns the store used to index blockchain targets
func blockchainTargetIndexedStore(store sdk.KVStore, target types.BlockchainAddress) (index.Store, error) {
	prefixedIndexStore := indexStore(store)
	return index.NewIndexedStore(prefixedIndexStore, blockchainTargetsPrefix, target)
}

// certificatesIndexedStore returns the store used to index certificates
func certificatesIndexedStore(store sdk.KVStore, cert types.Certificate) (index.Store, error) {
	prefixedIndexStore := indexStore(store)
	return index.NewIndexedStore(prefixedIndexStore, certificatesPrefix, cert)
}

// mapCertificateToAccount maps given account to  a certificate
func (k Keeper) mapCertificateToAccount(ctx sdk.Context, account types.Account, certs ...types.Certificate) error {
	for _, cert := range certs {
		if len(cert) == 0 {
			continue
		}
		store, err := certificatesIndexedStore(ctx.KVStore(k.storeKey), cert)
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

// unmapCertificateToAccount removes an account associated to a certificate
func (k Keeper) unmapCertificateToAccount(ctx sdk.Context, account types.Account, certs ...types.Certificate) error {
	for _, cert := range certs {
		if len(cert) == 0 {
			continue
		}
		store, err := certificatesIndexedStore(ctx.KVStore(k.storeKey), cert)
		if err != nil {
			return err
		}
		if err = store.Delete(account); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) iterateCertificateAccounts(ctx sdk.Context, cert types.Certificate, do func(key []byte) bool) error {
	store, err := certificatesIndexedStore(ctx.KVStore(k.storeKey), cert)
	if err != nil {
		return err
	}
	store.IterateKeys(do)
	return nil
}

func (k Keeper) mapTargetToAccount(ctx sdk.Context, account types.Account, targets ...types.BlockchainAddress) error {
	for _, target := range targets {
		// if targets are empty ignore
		if target.Address == "" || target.ID == "" {
			continue
		}
		// otherwise map target to given account
		store, err := blockchainTargetIndexedStore(ctx.KVStore(k.storeKey), target)
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

func (k Keeper) unmapTargetToAccount(ctx sdk.Context, account types.Account, targets ...types.BlockchainAddress) error {
	for _, target := range targets {
		// if targets are empty then ignore the process
		if target.ID == "" || target.Address == "" {
			continue
		}
		store, err := blockchainTargetIndexedStore(ctx.KVStore(k.storeKey), target)
		if err != nil {
			return err
		}
		if err = store.Delete(account); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) iterateBlockchainTargetsAccounts(ctx sdk.Context, target types.BlockchainAddress, do func(key []byte) bool) error {
	store, err := blockchainTargetIndexedStore(ctx.KVStore(k.storeKey), target)
	if err != nil {
		return err
	}
	store.IterateKeys(do)
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
	store.IterateKeys(do)
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
	store.IterateKeys(do)
	return nil
}
