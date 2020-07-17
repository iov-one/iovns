package crud

import (
	"crypto/rand"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"reflect"
	"testing"
)

func TestNewStore(t *testing.T) {
	_ = NewStore(testCtx, testKey, testCdc, []byte{0x0})
}

type testStoreObject struct {
	Key    string
	Index1 string
	Index2 string
}

func (t testStoreObject) PrimaryKey() PrimaryKey {
	return []byte(t.Key)
}

func (t testStoreObject) SecondaryKeys() []SecondaryKey {
	var sk []SecondaryKey
	if t.Index1 != "" {
		sk = append(sk, SecondaryKey{
			Key:         []byte(t.Index1),
			StorePrefix: []byte{0x1},
		})
	}
	if t.Index2 != "" {
		sk = append(sk, SecondaryKey{
			Key:         []byte(t.Index2),
			StorePrefix: []byte{0x2},
		})
	}
	return sk
}

func newTestStoreObject() *testStoreObject {
	key := make([]byte, 8)
	index1 := make([]byte, 8)
	index2 := make([]byte, 8)
	rand.Read(key)
	rand.Read(index1)
	rand.Read(index2)
	return &testStoreObject{
		Key:    fmt.Sprintf("%x", key),
		Index1: fmt.Sprintf("%x", index1),
		Index2: fmt.Sprintf("%x", index2),
	}
}

func TestStore(t *testing.T) {
	store := NewStore(testCtx, testKey, testCdc, []byte{0x0})
	obj := newTestStoreObject()
	obj.Index2 = string([]byte{ReservedSeparator, 0x2, 0x3}) // put reserved separator in
	store.Create(obj)
	cpy := new(testStoreObject)
	if !store.Read(obj.PrimaryKey(), cpy) {
		t.Fatal("object not found")
	}
	if !reflect.DeepEqual(cpy, obj) {
		t.Fatal("objects do not match")
	}
	// update object
	oldIndex := obj.Index2
	obj.Index2 = "updated"
	store.Update(obj.PrimaryKey(), obj)
	if !store.Read(obj.PrimaryKey(), cpy) {
		t.Fatal("object deleted after update")
	}
	// check if indexes were updated
	filter := store.Filter(&testStoreObject{Index2: "updated"})
	if !filter.Valid() {
		t.Fatal("index was not updated")
	}
	filter.Read(cpy)
	if !reflect.DeepEqual(cpy, obj) {
		t.Fatal("objects do not match")
	}
	// try read from deleted index
	filter = store.Filter(&testStoreObject{Index2: oldIndex})
	if filter.Valid() {
		t.Fatal("old index was not removed")
	}
	// delete object
	store.Delete(obj.PrimaryKey(), obj)
	// check if anything was left in the store
	it := store.raw.Iterator(nil, nil)
	defer it.Close()
	if it.Valid() {
		t.Fatal("nothing should be in the store")
	}
}

func TestFilterAndIterateKeys(t *testing.T) {
	x := 5
	objs := make([]*testStoreObject, x)
	for i := 0; i < x; i++ {
		objs[i] = newTestStoreObject()
		objs[i].Index1 = objs[0].Index1 // set same index
	}
	// create objects
	store := NewStore(testCtx, testKey, testCdc, []byte{0x0})
	// create objects
	for _, obj := range objs {
		store.Create(obj)
	}
	// check if object number is correct
	keys := make([]PrimaryKey, 0, x)
	store.IterateKeys(func(pk PrimaryKey) bool {
		keys = append(keys, pk)
		return true
	})
	if len(keys) != x {
		t.Fatal("unexpected number of keys", len(keys), x)
	}
	// delete based on filter
	filter := store.Filter(&testStoreObject{Index1: objs[0].Index1}) // delete based on same index
	for ; filter.Valid(); filter.Next() {
		filter.Delete()
	}
	// try to read primary keys and check if they're deleted
	for i, obj := range objs {
		if store.objects.Has(obj.PrimaryKey()) {
			t.Fatalf("key not deleted %d %#v", i, obj)
		}
	}
	// now try to iterate indexes and check if index keys were removed
	for _, obj := range objs {
		for _, sk := range obj.SecondaryKeys() {
			if prefix.NewStore(store.indexes, sk.StorePrefix).Has(sk.Key) {
				t.Fatal("index not removed")
			}
		}
	}
	// try to filter again
	filter = store.Filter(&testStoreObject{Index1: objs[0].Index1})
	if filter.Valid() {
		t.Fatal("no valid keys should exist")
	}
	// try to filter by secondary index
	filter = store.Filter(&testStoreObject{Index2: objs[0].Index2})
	if filter.Valid() {
		t.Fatalf("no valid keys should exist %s", filter.primaryKeys[0])
	}
	iterator := store.raw.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		t.Logf("%s", iterator.Key())
	}

}
