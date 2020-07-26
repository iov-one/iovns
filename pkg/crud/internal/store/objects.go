package store

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/crud/types"
)

// newObjectsStore builds the object store
func newObjectsStore(cdc *codec.Codec, kv sdk.KVStore) _objects {
	return _objects{cdc: cdc, store: prefix.NewStore(kv, objectPrefix)}
}

// _objects is the objects store
type _objects struct {
	cdc *codec.Codec
	store sdk.KVStore
}

// create creates an object in the store
func (s _objects) create(o types.Object) {
	pk := o.PrimaryKey()
	if s.store.Has(pk.Key()) {
		panic(fmt.Errorf("cannot re-create an existing object with primary key: %x", pk.Key()))
	}
	s.store.Set(pk.Key(), s.encode(o))
}

// delete deletes an object from the store
func (s _objects) delete(pk types.PrimaryKey) {
	if !s.store.Has(pk.Key()) {
		panic(fmt.Errorf("cannot delete non existing object with primary key: %x", pk.Key()))
	}
	s.store.Delete(pk.Key())
}

// update updates the object in the store
func (s _objects) update(o types.Object) {
	pk := o.PrimaryKey()
	if !s.store.Has(pk.Key()) {
		panic(fmt.Errorf("cannot update non existing object with primary key: %x", pk.Key()))
	}
	s.store.Set(o.PrimaryKey().Key(), s.encode(o))
}

// read reads to the given target using the provided primary key
func (s _objects) read(pk types.PrimaryKey, target types.Object) bool {
	v := s.store.Get(pk.Key())
	if v == nil {
		return false
	}
	s.decode(v, target)
	return true
}

// iterate iterates
func (s _objects) iterate(do func(pk types.PrimaryKey) bool) {
	it := s.store.Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		key := types.NewPrimaryKey(it.Key())
		if !do(key) {
			break
		}
	}
}

func (s _objects) encode(o interface{}) []byte {
	if e, ok := o.(encoder); ok {
		o = e.MarshalCRUD()
	}
	return s.cdc.MustMarshalBinaryBare(o)
}

func (s _objects) decode(b []byte, o interface{}) {
	e, ok := o.(encoder)
	if !ok {
		s.cdc.MustUnmarshalBinaryBare(b, o)
		return
	}
	e.UnmarshalCRUD(s.cdc, b)
}

type encoder interface {
	MarshalCRUD() interface{}
	UnmarshalCRUD(cdc *codec.Codec, b []byte)
}