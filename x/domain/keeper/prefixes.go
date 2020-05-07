package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types"
)

var (
	// DomainStorePrefix is the prefix used to define the prefixed store containing domain data
	DomainStorePrefix = []byte{0x00}
	// AccountPrefixStore is the prefix used to define the prefixed store containing account data
	AccountStorePrefix = []byte{0x01}
	// IndexStorePrefix is the prefix used to defines the prefixed store containing indexing data
	IndexStorePrefix = []byte{0x02}
)

// domainStore returns the domain store from the module's kvstore
func domainStore(store types.KVStore) types.KVStore {
	return prefix.NewStore(store, DomainStorePrefix)
}

// accountStore returns the account store from the module's kvstore
func accountStore(store types.KVStore) types.KVStore {
	return prefix.NewStore(store, AccountStorePrefix)
}

// indexStore returns the indexing store from the module's kvstore
func indexStore(store types.KVStore) types.KVStore {
	return prefix.NewStore(store, IndexStorePrefix)
}

// accountInDomainsStore returns the prefixed store containing
// all the account keys contained in a domain
func accountsInDomainStore(store types.KVStore, domain string) types.KVStore {
	// get account store
	accountStore := accountStore(store)
	// get accounts in domain store
	return prefix.NewStore(accountStore, getDomainPrefixKey(domain))
}

// getDomainPrefixKey returns the domain prefix byte key
func getDomainPrefixKey(domainName string) []byte {
	return []byte(domainName)
}

// getAccountKey returns the account byte key by its name
func getAccountKey(accountName string) []byte {
	return []byte(accountName)
}

// accountKeyToString converts account key bytes to string
func accountKeyToString(accountKeyBytes []byte) string {
	return string(accountKeyBytes)
}
