package store

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/crud/types"
	"sort"
)

// _indexes is the indexes store, it simply uses a store in which it saves
// pointers to primary keys given an index with a fixed prefix
// and a pointer to the whole list of indexes in order to delete them when
// required, during object updates and deletions
type _indexes struct {
	pointers sdk.KVStore // pointers is the store which contains the pointers to primary keys
	list sdk.KVStore // list contains the list of indexes
	cdc *codec.Codec
}

// newIndexes is the _indexes constructor
func newIndexes(cdc *codec.Codec, store sdk.KVStore) indexes {
	return _indexes{
		pointers: prefix.NewStore(store, indexPrefix),
		list:     prefix.NewStore(store, indexListPrefix),
		cdc:      cdc,
	}
}

// storeFromSecondaryKey returns the prefixed key value store for the given secondary key
func (s _indexes) storeFromSecondaryKey(sk types.SecondaryKey) sdk.KVStore {
	indexPrefix := prefix.NewStore(s.pointers, []byte{sk.Prefix()})
	return prefix.NewStore(indexPrefix, sk.Key())
}

// create creates indexes for the given object
func (s _indexes) create(o types.Object) {
	sks := o.SecondaryKeys()
	pk := o.PrimaryKey()
	for _, sk := range sks {
		store := s.storeFromSecondaryKey(sk)
		store.Set(pk.Key(), []byte{})
	}
	s.createIndexList(pk, sks)
}

// delete deletes the indexes that point to the given primary key
func (s _indexes) delete(pk types.PrimaryKey) {
	secondaryKeys := s.getSecondaryKeysFromPrimary(pk)
	for _, sk := range secondaryKeys {
		store := s.storeFromSecondaryKey(sk)
		store.Delete(pk.Key())
	}
	// remove indexes
	s.deleteIndexList(pk)
}

// iterate finds all the primary keys to which the given secondary key points to
func (s _indexes) iterate(sk types.SecondaryKey, do func(pk types.PrimaryKey) bool) {
	store := s.storeFromSecondaryKey(sk)
	it := store.Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		primaryKey := types.NewPrimaryKey(it.Key())
		if !do(primaryKey) {
			break
		}
	}
}

// getSecondaryKeysFromPrimary returns all the secondary keys associated with an object with the given primary key
func (s _indexes) getSecondaryKeysFromPrimary(pk types.PrimaryKey) []types.SecondaryKey {
	v := s.list.Get(pk.Key())
	if v == nil {
		panic(fmt.Sprintf("no index exists for given key: %x", pk))
	}
	indexes := new(marshalledIndexes)
	s.cdc.MustUnmarshalBinaryBare(v, indexes)
	secondaryKeys := make([]types.SecondaryKey, len(indexes.Indexes))
	for i, index := range indexes.Indexes {
		sk := types.NewSecondaryKeyFromBytes(index)
		secondaryKeys[i] = sk
	}
	return secondaryKeys
}

// createIndexList creates the pointer to the whole list of secondary keys associated with an object's primary key
func (s _indexes) createIndexList(pk types.PrimaryKey, sks []types.SecondaryKey) {
	indexes := make([][]byte, len(sks))
	for i, sk := range sks {
		key := sk.Marshal()
		indexes[i] = key
	}
	// order it
	sort.Slice(indexes, func(i, j int) bool {
		return bytes.Compare(indexes[i], indexes[j]) < 0
	})
	// check if create/delete flow is correctly applied
	if s.list.Has(pk.Key()) {
		panic(fmt.Sprintf("index list for primary key %x should have been deleted and then reset", pk.Key()))
	}
	// set them
	s.list.Set(pk.Key(), s.cdc.MustMarshalBinaryBare(marshalledIndexes{Indexes: indexes}))
}

// deleteIndexList deletes the pointer of the list of secondary keys associated with an object's primary key
func (s _indexes) deleteIndexList(pk types.PrimaryKey) {
	if !s.list.Has(pk.Key()) {
		panic(fmt.Sprintf("cannot remove index list because it does not exist for key: %x", pk))
	}
	s.list.Delete(pk.Key())
}

// marshalledIndexes is the byte array containing
type marshalledIndexes struct {
	Indexes [][]byte
}