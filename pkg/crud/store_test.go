package crud

import (
	"reflect"
	"testing"
)

func TestNewStore(t *testing.T) {
	_ = NewStore(testCtx, testKey, testCdc, []byte{0x0})
}

type testStoreObject struct {
	Primary   string `crud:"primaryKey"`
	Secondary string `crud:"secondaryKey,01"`
}

func TestStore(t *testing.T) {
	store := NewStore(testCtx, testKey, testCdc, []byte{0x0})
	obj := &testStoreObject{
		Primary:   "account",
		Secondary: "address",
	}
	store.Create(obj)
	var readObj = new(testStoreObject)
	ok := store.Read([]byte("account"), readObj)
	if !ok {
		t.Fatal("object not found")
	}
	updateObj := &testStoreObject{
		Primary:   "account",
		Secondary: "third",
	}
	store.Update(updateObj)
	ok = store.Read([]byte("account"), readObj)
	if !ok {
		t.Fatal("object not found")
	}
	if !reflect.DeepEqual(updateObj, readObj) {
		t.Fatal("updates do not match")
	}
	store.Delete(updateObj)
	ok = store.Read([]byte("account"), readObj)
	if ok {
		t.Fatal("object still exists")
	}
}
