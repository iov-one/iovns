package index

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types"
)

// Index defines the behaviour
// of a type that can index itself
// into an unique byte key
type Indexer interface {
	Index() []byte
}

type Store struct {
	kv types.KVStore
}

// NewIndexedStore returns a prefixed indexed Store
// with the provided indexer key
func NewIndexedStore(kv types.KVStore, pref []byte, indexer Indexer) Store {
	// get prefixed Store
	prefixedStore := prefix.NewStore(kv, pref)
	indexedStore := prefix.NewStore(prefixedStore, indexer.Index())
	// get Store from indexKey
	return Store{indexedStore}
}

// IterateKeys iterates over keys given an Indexer
// performing the do function on those keys, if 'do'
// returns false then the iteration stops
// CONTRACT: while IterateKeys is running no operations
// can be performed on the kv Store associated with Store
func (s Store) IterateKeys(do func(b []byte) bool) {
	iterator := types.KVStorePrefixIterator(s.kv, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if key := iterator.Key(); !do(key) {
			return
		}
	}
}

// Set sets a key in the index returned by indexer
func (s Store) Set(key []byte) {
	s.kv.Set(key, []byte{})
}

// Delete deletes a key from the indexed Store
func (s Store) Delete(key []byte) bool {
	if !s.kv.Has(key) {
		return false
	}
	s.kv.Delete(key)
	return true
}
