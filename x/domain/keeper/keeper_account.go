package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/domain/types"
)

// contains all the functions to interact with the account store

// GetAccount finds an account based on its key name, if not found it will return
// a zeroed account and false.
func (k Keeper) GetAccount(ctx sdk.Context, domainName, accountName string) (account types.Account, exists bool) {
	// get domain prefix key
	domainKey := getDomainPrefixKey(domainName)
	store := prefix.NewStore(ctx.KVStore(k.accountStoreKey), domainKey)
	// get account key
	accountKey := getAccountKey(accountName)
	// get account
	accountBytes := store.Get(accountKey)
	if accountBytes == nil {
		return
	}
	// key exists
	exists = true
	k.cdc.MustUnmarshalBinaryBare(accountBytes, &account)
	return
}

// SetAccount inserts an account in the KVStore
func (k Keeper) SetAccount(ctx sdk.Context, account types.Account) {
	// get domain prefix key and account key
	domainKey, accountKey := getDomainPrefixKey(account.Domain), getAccountKey(account.Name)
	// get prefixed store
	store := prefix.NewStore(ctx.KVStore(k.accountStoreKey), domainKey)
	// set store
	store.Set(accountKey, k.cdc.MustMarshalBinaryBare(account))
}

// DeleteAccount deletes an account based on it full account name -> domain + iovns.Separator + account
func (k Keeper) DeleteAccount(ctx sdk.Context, domainName, accountName string) {
	domainKey := getDomainPrefixKey(domainName)
	accountKey := getAccountKey(accountName)
	store := prefix.NewStore(ctx.KVStore(k.accountStoreKey), domainKey)
	store.Delete(accountKey)
}

// GetAccountsInDomain provides all the account keys related to the given domain name
func (k Keeper) GetAccountsInDomain(ctx sdk.Context, domainName string) [][]byte {
	// get store
	accountStore := prefix.NewStore(ctx.KVStore(k.accountStoreKey), []byte(domainName))
	// create iterator
	iterator := accountStore.Iterator(nil, nil)
	defer iterator.Close()
	// create keys
	var domainAccountKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		// append
		domainAccountKeys = append(domainAccountKeys, iterator.Key())
	}
	// return keys
	return domainAccountKeys
}

// TransferAccount transfers the account to a new owner after resetting certificates and targets
func (k Keeper) TransferAccount(ctx sdk.Context, account types.Account, newOwner sdk.AccAddress) {
	// update account
	account.Owner = newOwner   // transfer owner
	account.Certificates = nil // remove certs
	account.Targets = nil      // remove targets
	// save account
	k.SetAccount(ctx, account)
}
