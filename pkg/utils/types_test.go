package utils

import (
	"reflect"
	"testing"
)

func TestUnderlyingValue(t *testing.T) {
	expectedKind := reflect.Int
	i := new(int)
	ii := &i
	iii := &ii
	x := UnderlyingValue(reflect.ValueOf(iii))
	if x.Kind() != expectedKind {
		t.Fatalf("unexpected kind %s", x.Kind())
	}
}
