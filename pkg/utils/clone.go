package utils

import (
	"reflect"
)

// CloneFromValue clones an arbitrary type
// and returns a new zeroed type of the
// instance, IT DOES NOT COPY THE CONTENTS.
// CONTRACT: must be a pointer
func CloneFromValue(x interface{}) interface{} {
	return CloneFromType(GetPtrType(x))
}

// CloneFromType returns a new
// zeroed type of the given type value
// CONTRACT: typ must be returned from reflect.Value.Type().Elem()
func CloneFromType(typ reflect.Type) interface{} {
	return reflect.New(typ).Interface()
}

// GetPtrType returns the pointer type
// CONTRACT: ptr must be a pointer
func GetPtrType(ptr interface{}) reflect.Type {
	return reflect.ValueOf(ptr).Type().Elem()
}
