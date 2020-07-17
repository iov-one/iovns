package crud

import (
	"reflect"
	"testing"
)

func Test_filter(t *testing.T) {
	set1 := make(keySet)
	set1.Insert([]byte("1"))
	set1.Insert([]byte("2"))
	set1.Insert([]byte("5"))
	set2 := make(keySet)
	set2.Insert([]byte("2"))
	set2.Insert([]byte("3"))
	set3 := make(keySet)
	set3.Insert([]byte("5"))
	set3.Insert([]byte("2"))
	expected := []PrimaryKey{PrimaryKey("2")}
	if !reflect.DeepEqual(expected, filter([]set{set1, set2, set3})) {
		t.Fatal("unexpected result")
	}
}
