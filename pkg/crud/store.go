package crud

import (
	"bytes"
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/tutils"
)

var indexPrefix = []byte{0x01}
var objectPrefix = []byte{0x00}

const ReservedSeparator = 0xFF

// Store defines a crud object store
// the store creates two sub-stores
// using prefixing, one is used to store objects
// the other one is used to store the indexes of
// the object.
type Store struct {
	cdc *codec.Codec

	indexes sdk.KVStore
	objects sdk.KVStore
}

// NewStore generates a new crud.Store given a context, a store key, the codec and a unique prefix
// that can be specified as nil if not required, the prefix generally serves the purpose of splitting
// a store into different stores in case different objects have to coexist in the same store.
func NewStore(ctx sdk.Context, key sdk.StoreKey, cdc *codec.Codec, uniquePrefix []byte) Store {
	store := ctx.KVStore(key)
	if len(uniquePrefix) != 0 {
		store = prefix.NewStore(store, uniquePrefix)
	}
	return Store{
		indexes: prefix.NewStore(store, indexPrefix),
		cdc:     cdc,
		objects: prefix.NewStore(store, objectPrefix),
	}
}

// Create creates a new object in the object store and writes its indexes
func (s Store) Create(o interface{}) {
	// inspect
	primaryKey, secondaryKeys, err := inspect(o)
	if err != nil {
		panic(err)
	}
	// marshal object
	objectBytes := s.cdc.MustMarshalBinaryBare(o)
	// save object to object store using its primary key
	s.objects.Set(primaryKey, objectBytes)
	// generate indexes
	s.index(primaryKey, secondaryKeys)
}

// Read reads in the object store and returns false if the object is not found
// if it is found then the binary is unmarshalled into the Object.
// CONTRACT: Object must be a pointer for the unmarshalling to take effect.
func (s Store) Read(key []byte, o interface{}) (ok bool) {
	v := s.objects.Get(key)
	if v == nil {
		return
	}
	s.cdc.MustUnmarshalBinaryBare(v, o)
	return true
}

func (s Store) ReadFilter(filter interface{}, o interface{}) (ok bool) {
	pk, sk, err := getKeys(filter)
	if err != nil {
		panic(err)
	}
	if pk != nil {
		return s.Read(pk, o)
	}
	if len(sk) != 0 {
		return s.ReadFromIndex(sk[0], o)
	}
	panic("no valid filter provided")
}

// IterateFilter takes care of iterating over an index given a filter
func (s Store) IterateFilter(filter interface{}, do func(pk PrimaryKey) bool) {
	pk, sk, err := getKeys(filter)
	if err != nil {
		panic(err)
	}
	if pk != nil {
		panic("primary keys are not expected in index iterations")
	}
	if len(sk) != 0 {
		s.IterateIndex(sk[0], do)
		return
	}
	panic("no valid filter provided")
}

// ReadFromIndex gets the first primary key of the given object from the index
func (s Store) ReadFromIndex(index SecondaryKey, o interface{}) (ok bool) {
	var primaryKey PrimaryKey
	s.IterateIndex(index, func(key PrimaryKey) bool {
		primaryKey = key
		return false
	})
	if primaryKey == nil {
		return false
	}
	ok = s.Read(primaryKey, o)
	if !ok {
		panic("key found in index but not on store")
	}
	return
}

func (s Store) IterateIndex(index SecondaryKey, do func(key PrimaryKey) bool) {
	indexStore := prefix.NewStore(s.indexes, index.StorePrefix)
	iterator := sdk.KVStorePrefixIterator(indexStore, index.Key)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if primaryKey := iterator.Key(); !do(primaryKey) {
			return
		}
	}
}

// Update updates the given Object in the objects store
// after clearing the indexes and reapplying them based on the
// new update.
// To achieve so a zeroed copy of Object is created which is used to
// unmarshal the old object contents which is necessary for the un-indexing.
func (s Store) Update(o interface{}) {
	primaryKey, _, err := inspect(o)
	if err != nil {
		panic(err)
	}
	// get old copy of the object marshalValue
	oldObjBytes := s.objects.Get(primaryKey)
	if oldObjBytes == nil {
		panic("trying to update a non existing object")
	}
	// copy the object
	objCopy := tutils.CloneFromValue(o)
	// unmarshal
	s.cdc.MustUnmarshalBinaryBare(oldObjBytes, objCopy)
	// remove old object
	s.Delete(objCopy)
	// create new object
	s.Create(o)
	// add new indexes

	// update object
	s.objects.Set(primaryKey, s.cdc.MustMarshalBinaryBare(o))
}

// Delete deletes an object from the object store after
// clearing its indexes.
func (s Store) Delete(o interface{}) {
	primaryKey, secondaryKey, err := inspect(o)
	if err != nil {
		panic(err)
	}
	s.unindex(primaryKey, secondaryKey)
	s.objects.Delete(primaryKey)
}

// unindex removes the indexes values related to the given object
func (s Store) unindex(primaryKey PrimaryKey, secondaryKeys []SecondaryKey) {
	s.opIndex(secondaryKeys, func(s sdk.KVStore) bool {
		s.Delete(primaryKey)
		return false
	})
}

// index indexes the secondary key values related to the object
func (s Store) index(primaryKey PrimaryKey, secondaryKeys []SecondaryKey) {
	s.opIndex(secondaryKeys, func(s sdk.KVStore) bool {
		s.Set(primaryKey, []byte{})
		return true
	})
}

// opIndex does operations on indexes given an object and a function to process indexed objects
func (s Store) opIndex(secondaryKeys []SecondaryKey, do func(s sdk.KVStore) bool) {
	for _, sk := range secondaryKeys {
		fixSecondaryKey(&sk)
		// move into the prefixed store of the index
		store := prefix.NewStore(s.indexes, sk.StorePrefix)
		// move into the prefixed store of the index value, the index is hence a set
		store = prefix.NewStore(store, sk.Key)
		if !do(store) {
			break
		}
	}
}

// Object defines an object in which we can do crud operations
type Object interface {
	// PrimaryKey returns the unique key of the object
	PrimaryKey() PrimaryKey
	// SecondaryKeys returns the secondary keys used to index the object
	SecondaryKeys() []SecondaryKey
}

// PrimaryKey defines a primary key, which is a secondary key, under the hood, but with a fixed 0x0 prefix
type PrimaryKey []byte

// SecondaryKey defines a secondary key for the object
type SecondaryKey struct {
	// Key is the byte key which identifies the byte key prefix used to iterate of the index of the secondary key
	Key []byte
	// StorePrefix is the prefix of the index, necessary to divide one index from another
	StorePrefix []byte
}

// fixSecondaryKey encodes key value which contains the reserved separator by base64-encoding them
// this is necessary because we're dealing with a prefixed KVStore which, if we iterate, is going to
// iterate over bytes contained in a key, so if we assume we have:
// KeyA = [0x1]
// KeyB = [0x1, 0x2]
// during iteration, in case we wanted to iterate over KeyA only we'd end up in KeyB domain too because
// KeyB starts with KeyA, so to avoid this we put a full stop separator which we know other keys can not contain
func fixSecondaryKey(sk *SecondaryKey) {
	if bytes.Contains(sk.Key, []byte{ReservedSeparator}) {
		dst := make([]byte, base64.RawStdEncoding.EncodedLen(len(sk.Key)))
		base64.RawStdEncoding.Encode(dst, sk.Key)
		sk.Key = dst
	}
	sk.Key = append(sk.Key, ReservedSeparator)
}
