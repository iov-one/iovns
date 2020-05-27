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
		Domain:     domain.Name,
		Name:       "",
		Owner:      domain.Admin,
		ValidUntil: types.MaxValidUntil,
	}
	// apply changes in case domain is open
	if domain.Type == types.OpenDomain {
		// empty account expiration becomes domain expiration
		acc.ValidUntil = domain.ValidUntil
	}
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

// TransferDomain transfers aliceAddr domain
func (k Keeper) TransferDomain(ctx sdk.Context, newOwner sdk.AccAddress, domain types.Domain, reset types.DomainResetFlag) {
	// unmap domain owner
	err := k.unmapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
	// update domain owner
	domain.Admin = newOwner
	// update domain in kvstore
	k.SetDomain(ctx, domain)
	// transfer empty account
	emptyAcc, _ := k.GetAccount(ctx, domain.Name, iovns.EmptyAccountName)
	k.TransferAccount(ctx, emptyAcc, newOwner)
	// operate based on reset flag

	// if the flag is reset none, then just skip everything
	if reset == types.ResetNone {
		return
	}
	// get account keys related to the domain; TODO an multikey index owner:domain might make the iteration less expensive
	var accountKeys [][]byte
	k.GetAccountsInDomain(ctx, domain.Name, func(key []byte) bool {
		accountKeys = append(accountKeys, key)
		return true
	})
	// iterate over accounts
	for _, accountKey := range accountKeys {
		// empty account key was already transferred so we are not expecting it here
		// get account
		account, _ := k.GetAccount(ctx, domain.Name, accountKeyToString(accountKey))
		// if reset flag is owned only and current owner is not domain (old) owner then skip
		if reset == types.ResetOwned && !account.Owner.Equals(domain.Admin) {
			continue
		}
		// transfer account
		k.TransferAccount(ctx, account, newOwner)
	}
	// map domain to new owner
	err = k.mapDomainToOwner(ctx, domain)
	if err != nil {
		panic(fmt.Errorf("indexing error: (%#v): %w", domain, err))
	}
}
