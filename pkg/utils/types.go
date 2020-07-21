package utils

import "reflect"

// UnderlyingValue gets the underlying value of a reflect.Value
func UnderlyingValue(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}
	v = v.Elem()
	return UnderlyingValue(v)
}

func StrPtr(str string) *string {
	return &str
}
