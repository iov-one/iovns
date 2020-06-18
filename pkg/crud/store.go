package crud

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/tutils"
)

var indexPrefix = []byte{0x01}
var objectPrefix = []byte{0x02}

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
func (s Store) Create(o Object) {
	// index then create
	s.index(o)
	key := o.PrimaryKey()
	s.objects.Set(key, s.cdc.MustMarshalBinaryBare(o))
}

// Read reads in the object store and returns false if the object is not found
// if it is found then the binary is unmarshalled into the Object.
// CONTRACT: Object must be a pointer for the unmarshalling to take effect.
func (s Store) Read(key []byte, o Object) (ok bool) {
	v := s.objects.Get(key)
	if v == nil {
		return
	}
	s.cdc.MustUnmarshalBinaryBare(v, o)
	return true
}

// ReadFromIndex gets the first primary key of the given object from the index
func (s Store) ReadFromIndex(index SecondaryKey, o Object) (ok bool) {
	var primaryKey []byte
	s.IterateIndex(index, func(key []byte) bool {
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

func (s Store) IterateIndex(index SecondaryKey, do func(key []byte) bool) {
	indexStore := prefix.NewStore(s.indexes, index.StorePrefix())
	iterator := sdk.KVStorePrefixIterator(indexStore, index.Key())
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
func (s Store) Update(o Object) {
	key := o.PrimaryKey()
	// get old copy of the object marshalValue
	oldObjBytes := s.objects.Get(key)
	if oldObjBytes == nil {
		panic("trying to update a non existing object")
	}
	// copy the object
	objCopy := tutils.CloneFromValue(o)
	// unmarshal
	s.cdc.MustUnmarshalBinaryBare(oldObjBytes, objCopy)
	// remove old indexes
	s.unindex(objCopy.(Object))
	// update object
	s.objects.Set(key, s.cdc.MustMarshalBinaryBare(o))
}

// Delete deletes an object from the object store after
// clearing its indexes.
func (s Store) Delete(o Object) {
	s.unindex(o)
	s.objects.Delete(o.PrimaryKey())
}

// unindex removes the indexes values related to the given object
func (s Store) unindex(o Object) {
	s.opIndex(o, func(idx index.Store, obj index.Indexed) {
		err := idx.Delete(obj)
		if err != nil {
			panic(err)
		}
	})
}

// index indexes the values related to the object
func (s Store) index(o Object) {
	s.opIndex(o, func(idx index.Store, obj index.Indexed) {
		err := idx.Set(obj)
		if err != nil {
			panic(err)
		}
	})
}

// opIndex does operations on indexes given an object and a function to process indexed objects
func (s Store) opIndex(o Object, op func(idx index.Store, obj index.Indexed)) {
	for _, idx := range o.SecondaryKeys() {
		indx, err := index.NewIndexedStore(s.indexes, idx.Prefix, nil)
		if err != nil {
			panic(fmt.Sprintf("unable to index object: %s", err))
		}
		op(indx, idx.Indexed)
	}
}

// Index defines a simple index which consists of index and indexed object
type Index struct {
	// Prefix is the unique prefix for the index
	// CONTRACT: this must be constant across different objects of the same type
	Prefix []byte
	// Indexed defines the object that has to be indexed
	Indexed index.Indexed
}

// Object defines an object in which we can do crud operations
type Object interface {
	// PrimaryKey returns the unique key of the object
	PrimaryKey() []byte
	// SecondaryKeys returns the secondary keys used to index the object
	SecondaryKeys() []Index
}

// SecondaryKey defines a secondary key for the object
type SecondaryKey interface {
	Key() []byte
	StorePrefix() []byte
}
