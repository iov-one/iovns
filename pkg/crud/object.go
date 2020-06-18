package crud

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/iov-one/iovns/tutils"
	"math"
	"reflect"
	"strings"
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

func getKeys(v reflect.Value) (key, []key) {
	typ := v.Type()
	var primaryKey key
	var secondaryKeys []key
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// TODO check if type implements indexer interface
		// ignore unexported fields
		if field.Anonymous {
			continue
		}
		// get field value
		fieldValue := v.FieldByName(field.Name)
		// check field type by tags
		tag, ok := field.Tag.Lookup(TagName)
		// if tag is missing then no indexing is required
		if !ok {
			continue
		}
		// check tag type
		split := strings.Split(tag, ",")
		switch split[0] {
		// check if primary key or secondary key
		case PrimaryKeyTag:
			// check if a primary key was already specified
			if primaryKey.value != nil {
				panic("crud: only one primary key can be specified for each object")
			}
			valueBytes := marshalValue(fieldValue)
			primaryKey = key{
				prefix: []byte{0x0},
				value:  valueBytes,
			}
		case SecondaryKeyTag:
			prefix, err := hex.DecodeString(split[1])
			if err != nil {
				panic("crud: invalid hex prefix in key")
			}
			secondaryKey := key{
				prefix: prefix,
				value:  marshalValue(fieldValue),
			}
			secondaryKeys = append(secondaryKeys, secondaryKey)
		}
	}
	if primaryKey.value == nil {
		panic(fmt.Sprintf("crud: no primary key specified in type: %T", v.Interface()))
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

// key defines a database key
// adapted for key value stores
// that use byte prefixing
type key struct {
	prefix []byte
	value  []byte
}
