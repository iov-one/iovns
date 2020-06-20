package crud

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_marshalValue(t *testing.T) {
	x := "hello"
	b := marshalValue(reflect.ValueOf(x))
	t.Logf("%x", b)
}

func Test_marshalSlice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		x := []byte{0x01}
		got := marshalSlice(reflect.ValueOf(x))
		if !bytes.Equal(got, x) {
			t.Fatal("unexpected result")
		}
	})
	t.Run("panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("panic expected")
			}
		}()
		x := []string{"hi"}
		_ = marshalSlice(reflect.ValueOf(x))

	})
}

func Test_getKeys(t *testing.T) {
	type obj struct {
		PK string `crud:"primaryKey"`
		SK string `crud:"secondaryKey,01"`
	}
	pk, sk, err := getKeys(reflect.ValueOf(obj{
		PK: "test1",
		SK: "test2",
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v %#v", pk, sk)
}
