package executor

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/x/domain/types"
)

// IndexStorePrefix is the prefix used to defines the prefixed store containing indexing data
var IndexStorePrefix = []byte{0x02}

// DomainStorePrefix is the prefix used to define the prefixed store containing prefixedStore data
var DomainStorePrefix = []byte{0x00}

// AccountPrefixStore is the prefix used to define the prefixed store containing account data
var AccountStorePrefix = []byte{0x01}
var ownerToDomainPrefix = []byte{0x05}
var ownerToAccountPrefix = []byte{0x04}
var blockchainTargetsPrefix = []byte{0x06}

type indexStore struct {
	store sdk.KVStore
}

func newIndexStore(ctx sdk.Context, key sdk.StoreKey) indexStore {
	store := ctx.KVStore(key)
	store = prefix.NewStore(store, IndexStorePrefix)
	return indexStore{
		store: store,
	}
}

func (i indexStore) ownerToDomain(addr sdk.AccAddress) (index.Store, error) {
	if addr.Empty() {
		panic("indexing of empty addresses is not allowed")
	}
	return index.NewAddressIndex(i.store, ownerToDomainPrefix, addr)
}

func (i indexStore) ownerToAccount(addr sdk.AccAddress) (index.Store, error) {
	if addr.Empty() {
		panic("indexing of empty addresses is not allowed")
	}
	return index.NewAddressIndex(i.store, ownerToAccountPrefix, addr)
}

func (i indexStore) targetToAccount(target types.Resource) (index.Store, error) {
	return index.NewIndexedStore(i.store, blockchainTargetsPrefix, target)
}

type accountStore struct {
	cdc      *codec.Codec
	store    sdk.KVStore
	idxStore indexStore
}

func newAccountStore(ctx sdk.Context, key sdk.StoreKey, cdc *codec.Codec) accountStore {
	return accountStore{
		store:    prefix.NewStore(ctx.KVStore(key), AccountStorePrefix),
		idxStore: newIndexStore(ctx, key),
		cdc:      cdc,
	}
}

// create creates a new account in the account store
func (a accountStore) create(account types.Account) {
	store := a.prefixedStore(account.Domain)
	a.index(account)
	store.Set(a.kv(account))
}

func (a accountStore) read(domain, name string) (types.Account, bool) {
	store := a.prefixedStore(domain)
	v := store.Get(a.key(name))
	if v == nil {
		return types.Account{}, false
	}
	var acc types.Account
	a.cdc.MustUnmarshalBinaryBare(v, &acc)
	return acc, true
}

func (a accountStore) update(account types.Account) {
	// check if account exists
	oldAccount, ok := a.read(account.Domain, account.Name)
	if !ok {
		panic(fmt.Sprintf("update operation on a non existing account is not allowed %+v", account))
	}
	// remove old indexes
	a.unindex(oldAccount)
	// add new indexes
	a.index(account)
	// update account
	store := a.prefixedStore(account.Domain)
	store.Set(a.kv(account))
}

func (a accountStore) delete(account types.Account) {
	// remove associated indexes
	a.unindex(account)
	// delete
	store := a.prefixedStore(account.Domain)
	store.Delete(a.key(account.Name))
}

func (a accountStore) index(account types.Account) {
	// index targets
	for _, trg := range account.Resources {
		if trg.URI == "" || trg.Resource == "" {
			continue
		}
		trgStore, err := a.idxStore.targetToAccount(trg)
		if err != nil {
			panic(fmt.Sprintf("unable to create index store for target %+v: %s", trg, err))
		}
		err = trgStore.Set(account)
		if err != nil {
			panic(fmt.Sprintf("unable to set account %+v: %s", account, err))
		}
	}
	// index owner to account
	ownStore, err := a.idxStore.ownerToAccount(account.Owner)
	if err != nil {
		panic(fmt.Sprintf("unable to create owner to account store %+v: %s", account, err))
	}
	err = ownStore.Set(account)
	if err != nil {
		panic(fmt.Sprintf("unable to set account: %+v: %s", account, err))
	}
}

// unindex removes all the index keys associated with the account
func (a accountStore) unindex(account types.Account) {
	// index targets
	for _, trg := range account.Resources {
		if trg.URI == "" || trg.Resource == "" {
			continue
		}
		trgStore, err := a.idxStore.targetToAccount(trg)
		if err != nil {
			panic(fmt.Sprintf("unable to create index store for target %+v: %s", trg, err))
		}
		err = trgStore.Delete(account)
		if err != nil {
			panic(fmt.Sprintf("unable to remove account %+v: %s", account, err))
		}
	}
	// index owner to account
	ownStore, err := a.idxStore.ownerToAccount(account.Owner)
	if err != nil {
		panic(fmt.Sprintf("unable to create owner to account store %+v: %s", account, err))
	}
	err = ownStore.Delete(account)
	if err != nil {
		panic(fmt.Sprintf("unable to remove account: %+v: %s", account, err))
	}
}

// key returns the unique account key
func (a accountStore) kv(account types.Account) ([]byte, []byte) {
	k := a.key(account.Name)
	v := a.cdc.MustMarshalBinaryBare(account)
	return k, v
}

func (a accountStore) key(name string) []byte {
	return []byte(name)
}

// prefixedStore returns the prefixed of the account based on its domain
func (a accountStore) prefixedStore(domain string) sdk.KVStore {
	domBytes := []byte(domain)
	if bytes.Contains(domBytes, []byte{index.ReservedSeparator}) {
		panic(fmt.Sprintf("domains with reserved separator should not be allowed: %s", domain))
	}
	prfx := append(domBytes, index.ReservedSeparator)
	return prefix.NewStore(a.store, prfx)
}

type domainStore struct {
}
