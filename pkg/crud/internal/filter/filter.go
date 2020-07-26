package filter

import (
	"bytes"
	"fmt"
	"github.com/iov-one/iovns/pkg/crud/types"
	"sort"
)

type store interface {
	Read(key types.PrimaryKey, o types.Object) bool
	Delete(key types.PrimaryKey)
	Update(o types.Object)
}

type Filtered struct {
	counter     int
	nKeys       int
	primaryKeys []types.PrimaryKey
	store       store
}

func NewFiltered(keys []types.PrimaryKey, store store) *Filtered {
	// sort deterministically
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i].Key(), keys[j].Key()) < 0
	})
	return &Filtered{
		counter:     0,
		nKeys:       len(keys),
		primaryKeys: keys,
		store:       store,
	}
}

func (f *Filtered) Read(o types.Object) {
	ok := f.store.Read(f.currKey(), o)
	if !ok {
		panic(fmt.Sprintf("can't find object using primary key: %s", f.currKey()))
	}
}

func (f *Filtered) Update(o types.Object) {
	if !bytes.Equal(o.PrimaryKey().Key(), f.currKey().Key()) {
		panic("trying to update objects with unmatching primary keys")
	}
	f.store.Update(o)
}

func (f *Filtered) Delete() {
	if !f.Valid() {
		return
	}
	f.store.Delete(f.currKey())
}

func (f *Filtered) currKey() types.PrimaryKey {
	if f.counter == 0 && f.nKeys == 0 {
		panic("iterating an empty filter is not valid")
	}
	return f.primaryKeys[f.counter]
}

func (f *Filtered) Next() {
	f.counter++
}

func (f *Filtered) Valid() bool {
	return f.counter < f.nKeys
}

