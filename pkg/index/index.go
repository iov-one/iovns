package index

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types"
	"log"
)

const ReservedSeparator byte = 0xFF

// Index defines the behaviour
// of a type that can index itself
// into an unique byte key
type Indexer interface {
	Index() ([]byte, error)
}

type Store struct {
	kv types.KVStore
}

func encode(src []byte) []byte {
	dst := make([]byte, base64.RawStdEncoding.EncodedLen(len(src)))
	base64.RawStdEncoding.Encode(dst, src)
	return dst
}

func decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.RawStdEncoding.DecodedLen(len(src)))
	_, err := base64.RawStdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func index(i Indexer) ([]byte, error) {
	indexKey, err := i.Index()
	if err != nil {
		return nil, err
	}
	if bytes.Contains(indexKey, []byte{ReservedSeparator}) {
		// TODO print a warning, receiving an index with the separator inside should not happen, my dear.
		log.Printf("Key %x, containing separator was encoded.", indexKey)
		indexKey = encode(indexKey)
	}
	indexKey = append(indexKey, ReservedSeparator)
	return indexKey, nil
}

// NewIndexedStore returns a prefixed indexed Store
// with the provided indexer key
func NewIndexedStore(kv types.KVStore, pref []byte, indexer Indexer) (Store, error) {
	// get indexing key
	indexingKey, err := index(indexer)
	if err != nil {
		return Store{}, err
	}
	// get prefixed store matching a certain index type
	prefixedStore := prefix.NewStore(kv, pref)
	// get prefixed store of the values matched by the index
	indexedStore := prefix.NewStore(prefixedStore, indexingKey)
	// return the indexed store
	return Store{indexedStore}, nil
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
func (s Store) Set(indexed Indexed) error {
	key, err := indexed.Pack()
	if err != nil {
		return err
	}
	s.kv.Set(key, []byte{})
	return nil
}

// Delete deletes a key from the indexed Store
func (s Store) Delete(indexed Indexed) error {
	key, err := indexed.Pack()
	if err != nil {
		return err
	}
	if !s.kv.Has(key) {
		return fmt.Errorf("key not found: %x", key)
	}
	s.kv.Delete(key)
	return nil
}
