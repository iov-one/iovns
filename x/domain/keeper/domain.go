package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns"
	"github.com/iov-one/iovns/x/domain/types"
)

// contains all the functions to interact with the domain store

// GetDomain returns the domain based on its name, if domain is not found ok will be false
func (k Keeper) GetDomain(ctx sdk.Context, domainName string) (domain types.Domain, ok bool) {
	store := domainStore(ctx.KVStore(k.storeKey))
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
	err := k.mapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
	// set domain
	k.SetDomain(ctx, domain)
	// generate empty name account
	acc := types.Account{
		Domain:       domain.Name,
		Name:         "",
		Owner:        domain.Admin, // TODO this is not clear, why the domain admin is zero address while this is msg.Admin
		ValidUntil:   types.MaxValidUntil,
		Targets:      nil,
		Certificates: nil,
		Broker:       nil, // TODO ??
	}
	// if domain type is open then account valid until needs to be updated
	if domain.Type == types.OpenDomain {
		acc.ValidUntil = domain.ValidUntil
	}
	// save account
	k.CreateAccount(ctx, acc)
}

// SetDomain updates or creates a new domain in the store
func (k Keeper) SetDomain(ctx sdk.Context, domain types.Domain) {
	store := domainStore(ctx.KVStore(k.storeKey))
	store.Set([]byte(domain.Name), k.cdc.MustMarshalBinaryBare(domain))
}

// IterateAllDomains will return an iterator for all the domain keys
// present in the KVStore, it's callers duty to close the iterator.
func (k Keeper) IterateAllDomains(ctx sdk.Context) []types.Domain {
	store := domainStore(ctx.KVStore(k.storeKey))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	var domains []types.Domain
	for ; iterator.Valid(); iterator.Next() {
		var d types.Domain
		domainBytes := store.Get(iterator.Key())
		k.cdc.MustUnmarshalBinaryBare(domainBytes, &d)
		domains = append(domains, d)
	}
	return domains
}

// DeleteDomain deletes the domain and the accounts in it
// this operation can only fail in case the domain does not exist
func (k Keeper) DeleteDomain(ctx sdk.Context, domainName string) (exists bool) {
	domain, exists := k.GetDomain(ctx, domainName)
	if !exists {
		return
	}
	// delete domain
	domainStore := domainStore(ctx.KVStore(k.storeKey))
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
	err := k.unmapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
	// done
	return true
}

// TransferDomainOwnership transfers the domain owner to newOwner
func (k Keeper) TransferDomainOwnership(ctx sdk.Context, domain types.Domain, newOwner sdk.AccAddress) {
	// unmap domain owner
	err := k.unmapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
	// update domain owner
	domain.Admin = newOwner
	// update domain in kvstore
	k.SetDomain(ctx, domain)
	// transfer empty domain account
	acc, _ := k.GetAccount(ctx, domain.Name, "")
	k.TransferAccount(ctx, acc, newOwner)
	// map domain to new owner
	err = k.mapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
}

// FlushDomain clears all the accounts in a domain, empty account excluded
func (k Keeper) FlushDomain(ctx sdk.Context, domain types.Domain) {
	// iterate accounts
	var accountKeys [][]byte
	k.GetAccountsInDomain(ctx, domain.Name, func(key []byte) bool {
		accountKeys = append(accountKeys, key)
		return true
	})
	// delete each account
	for _, accountKey := range accountKeys {
		accountName := accountKeyToString(accountKey)
		if accountName == types.EmptyAccountName {
			continue
		}
		k.DeleteAccount(ctx, domain.Name, accountName)
	}
}

// TransferDomainAccountsOwnedByAddr transfers all the accounts in the domain owned by an address to a new address
func (k Keeper) TransferDomainAccountsOwnedByAddr(ctx sdk.Context, domain types.Domain, currentOwner, newOwner sdk.AccAddress) {
	var accountKeys [][]byte
	k.GetAccountsInDomain(ctx, domain.Name, func(key []byte) bool {
		accountKeys = append(accountKeys, key)
		return true
	})
	for _, accountKey := range accountKeys {
		accountName := accountKeyToString(accountKey)
		acc, _ := k.GetAccount(ctx, domain.Name, accountName)
		// skip accounts not owned by the provided address
		if !acc.Owner.Equals(currentOwner) {
			continue
		}
		k.TransferAccount(ctx, acc, newOwner)
	}
}

// TransferDomainAll transfers the domain and the related accounts TODO deprecate
func (k Keeper) TransferDomainAll(ctx sdk.Context, newOwner sdk.AccAddress, domain types.Domain) {
	// unmap domain owner
	err := k.unmapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
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
		if string(accountKey) == types.EmptyAccountName {
			continue
		}
		// get account
		account, _ := k.GetAccount(ctx, domain.Name, accountKeyToString(accountKey))
		// transfer it
		k.TransferAccount(ctx, account, newOwner)
	}
	// map domain to new owner
	err = k.mapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
}

// RenewDomain takes care of renewing the domain expiration time based on the configuration
func (k *Keeper) RenewDomain(ctx sdk.Context, domain types.Domain) {
	// get configuration
	renewDuration := k.ConfigurationKeeper.GetDomainRenewDuration(ctx)
	// update domain valid until
	domain.ValidUntil = iovns.TimeToSeconds(
		iovns.SecondsToTime(domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// set domain
	k.SetDomain(ctx, domain)
}
