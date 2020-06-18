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
