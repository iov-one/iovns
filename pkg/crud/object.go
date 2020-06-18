package crud

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/iov-one/iovns/tutils"
	"math"
	"reflect"
)

const TagName = "crud"
const PrimaryKeyTag = "primaryKey"
const SecondaryKeyTag = "secondaryKey"

func inspect(o interface{}) {
	// TODO check if type implements object
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr {
		panic("crud: pointer expected")
	}
	v = tutils.UnderlyingValue(v)
	if v.Kind() != reflect.Struct {
		panic("crud: pointer to struct expected")
	}
	// find primary keys and secondary keys
	// primaryKey, secondaryKeys := getKeys(v)
}

func getKeys(v reflect.Value) ([]byte, [][]byte) {
	typ := v.Type()
	var primaryKey []byte
	var secondaryKeys [][]byte
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// TODO check if type implements indexer interface
		// ignore unexported fields
		if field.Anonymous {
			continue
		}
		// check field type by tags
		_, ok1 := field.Tag.Lookup(PrimaryKeyTag)
		_, ok2 := field.Tag.Lookup(SecondaryKeyTag)
		if ok1 && ok2 {
			panic(fmt.Sprintf("crud: field %s in type %s is both primary and secondary key", field.Name, field.Type.Name()))
		}
		if ok1 {
			primaryKey = marshalValue(v.FieldByName(field.Name))
		}
	}
	return primaryKey, secondaryKeys
}

var typesToBytes = map[reflect.Kind]func(v reflect.Value) []byte{
	reflect.String: func(v reflect.Value) []byte {
		return []byte(v.Interface().(string))
	},
	reflect.Float64: func(v reflect.Value) []byte {
		f64 := v.Interface().(float64)
		var buf []byte
		binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f64))
		return buf
	},
}

// notAllowedIndexType contains a set of first class citizen types that cannot be indexed in marshalValue
var notAllowedIndexType = map[reflect.Kind]struct{}{
	reflect.Struct:        {},
	reflect.UnsafePointer: {},
	reflect.Invalid:       {},
	reflect.Map:           {},
	reflect.Array:         {},
	reflect.Func:          {},
	reflect.Chan:          {},
	reflect.Interface:     {},
	reflect.Ptr:           {},
}

// typeToBytes converts an arbitrary type to bytes
func typeToBytes(v reflect.Value) []byte {
	i := v.Interface()
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		panic(fmt.Errorf("crud: unable to set type %T to bytes: %w", i, err))
	}
	return buf.Bytes()
}

// marshalValue gets marshalValue from reflect.Value
func marshalValue(v reflect.Value) []byte {
	v = tutils.UnderlyingValue(v)
	// check if forbidden type
	kind := v.Kind()
	if kind == reflect.Slice {
		return marshalSlice(v)
	}
	if _, ok := notAllowedIndexType[kind]; ok {
		panic(fmt.Sprintf("crud: value of type %s cannot be turned into a byte key", kind))
	}
	// now index based on type
	marshaler, ok := typesToBytes[kind]
	if !ok {
		return typeToBytes(v)
	}
	return marshaler(v)
}

func marshalSlice(v reflect.Value) []byte {
	if b, ok := v.Interface().([]byte); ok {
		return b
	}
	panic(fmt.Sprintf("crud: only slice types allowed are byte ones, got: %T", v.Interface()))
}

type secondaryKey struct {
	key    []byte
	prefix []byte
}
