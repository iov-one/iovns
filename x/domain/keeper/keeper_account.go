package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd"
	"github.com/iov-one/iovnsd/x/domain/types"
)

// contains all the functions to interact with the account store

// GetAccount finds an account based on its key name, if not found it will return
// a zeroed account and false.
func (k Keeper) GetAccount(ctx sdk.Context, accountName string) (account types.Account, exists bool) {
	store := ctx.KVStore(k.accountKey)
	accountBytes := store.Get([]byte(accountName))
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
	store := ctx.KVStore(k.accountKey)
	accountKey := iovnsd.GetAccountKey(account.Domain, account.Name)
	store.Set([]byte(accountKey), k.cdc.MustMarshalBinaryBare(account))
}
