package index

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types"
	"log"
)

// ReservedSeparator is the uint8 used to separate
// index key identifiers from the value of the key itself
const ReservedSeparator byte = 0xFF

// errKeyNotFound is for testing purposes and signals when a key was not found in the KVStore
var errKeyNotFound = errors.New("provided key is not inside the KVStore")

// Index defines the behaviour
// of a type that can index itself
// into an unique byte key
type Indexer interface {
	// Index should return an unique byte key for the object
	Index() ([]byte, error)
}

// Store contains the prefixed
// KVStore of an indexed entity
type Store struct {
	kv types.KVStore
}

// encode encodes bytes in base64 format
// it's used when index keys contain the
// reserved separator
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

// index takes an indexer and builds the unique
// defining key of it, returns an error only
// if the key can not index itself. It uses
// the reserved separator to signal the end of the
// index key, if the index contains the key then it is
// base64 encoded.
func index(i Indexer) ([]byte, error) {
	indexKey, err := i.Index()
	if err != nil {
		return nil, err
	}
	if bytes.Contains(indexKey, []byte{ReservedSeparator}) {
		// TODO print a warning, receiving an index with the separator inside should not happen, my dear.
		log.Printf("key %T:%x, containing reserved separator was encoded.", i, indexKey)
		indexKey = encode(indexKey)
	}
	indexKey = append(indexKey, ReservedSeparator)
	return indexKey, nil
}

// NewIndexedStore returns a prefixed indexed Store with the provided prefix + Indexer, it returns
// an error only if the indexer cannot marshal itself into a byte key
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

// IterateAllKeys iterates over all keys in the index
// performing the do function on those keys, if 'do'
// returns false then the iteration stops
// CONTRACT: while IterateAllKeys is running no operations
// can be performed on the kv Store associated with Store
func (s Store) IterateAllKeys(do func(b []byte) bool) {
	iterator := types.KVStorePrefixIterator(s.kv, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if key := iterator.Key(); !do(key) {
			return
		}
	}
}

// Set sets a key in the index, using an Indexed
// type that can marshal itself into bytes
// returns an error only if the key can not
// index itself into bytes
func (s Store) Set(indexed Indexed) error {
	key, err := indexed.Pack()
	if err != nil {
		return err
	}
	s.kv.Set(key, []byte{})
	return nil
}

// Delete deletes an Indexed item from the Index Store
// returns an error only if the item can not marshal
// itself into bytes, or if the key does not exist
func (s Store) Delete(indexed Indexed) error {
	key, err := indexed.Pack()
	if err != nil {
		return err
	}
	if !s.kv.Has(key) {
		return fmt.Errorf("%w: %x", errKeyNotFound, key)
	}
	s.kv.Delete(key)
	return nil
}
