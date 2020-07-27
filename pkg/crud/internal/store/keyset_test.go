package store

import (
	"github.com/iov-one/iovns/pkg/crud/types"
	"reflect"
	"testing"
)

func Test_filter(t *testing.T) {
	set1 := make(keySet)
	set1.Insert(types.NewPrimaryKeyFromString("1"))
	set1.Insert(types.NewPrimaryKeyFromString("2"))
	set1.Insert(types.NewPrimaryKeyFromString("5"))
	set2 := make(keySet)
	set2.Insert(types.NewPrimaryKeyFromString("2"))
	set2.Insert(types.NewPrimaryKeyFromString("3"))
	set3 := make(keySet)
	set3.Insert(types.NewPrimaryKeyFromString("5"))
	set3.Insert(types.NewPrimaryKeyFromString("2"))
	expected := []types.PrimaryKey{types.NewPrimaryKeyFromString("2")}
	result := primaryKeysFromSets([]set{set1, set2, set3})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("unexpected result got: %#v", result)
	}
}
