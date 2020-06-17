package crud

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/tutils"
)

var indexPrefix = []byte{0x03}
var objectPrefix = []byte{0x0}

type Index struct {
	Prefix  []byte
	Indexed index.Indexed
}

// Store defines the crud objects
type Store struct {
	cdc *codec.Codec

	indexes sdk.KVStore
	objects sdk.KVStore
}

// NewStore generates a new crud.Store
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

func (s Store) Create(o Object) {
	// index then create
	s.index(o)
	key := o.Key()
	s.objects.Set(key, s.cdc.MustMarshalBinaryBare(o))
}

func (s Store) Read(key []byte, o Object) (ok bool) {
	v := s.objects.Get(key)
	if v == nil {
		return
	}
	s.cdc.MustUnmarshalBinaryBare(v, o)
	return true
}

func (s Store) Update(o Object) {
	key := o.Key()
	// get old copy of the object bytes
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

func (s Store) Delete(o Object) {
	s.unindex(o)
	s.objects.Delete(o.Key())
}

func (s Store) unindex(o Object) {
	s.opIndex(o, func(idx index.Store, obj index.Indexed) {
		err := idx.Delete(obj)
		if err != nil {
			panic(err)
		}
	})
}

func (s Store) index(o Object) {
	s.opIndex(o, func(idx index.Store, obj index.Indexed) {
		err := idx.Set(obj)
		if err != nil {
			panic(err)
		}
	})
}

// opIndex defines an operation on an index
func (s Store) opIndex(o Object, op func(idx index.Store, obj index.Indexed)) {
	for _, idx := range o.Indexes() {
		indx, err := index.NewIndexedStore(s.indexes, idx.Prefix, o)
		if err != nil {
			panic(fmt.Sprintf("unable to index object: %s", err))
		}
		op(indx, idx.Indexed)
	}
}

// Object defines an object in which we can do crud operations
type Object interface {
	// Key returns the unique key of the object
	Key() []byte
	// Indexes returns the indexes of the object
	Indexes() []Index
	index.Indexer
}
