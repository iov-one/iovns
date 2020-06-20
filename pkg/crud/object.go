package crud

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/iov-one/iovns/tutils"
	"reflect"
	"strings"
)

const TagName = "crud"
const PrimaryKeyTag = "primaryKey"
const SecondaryKeyTag = "secondaryKey"
const PrimaryKeyPrefix = 0x0

func inspect(o interface{}) (primaryKey PrimaryKey, secondaryKey []SecondaryKey, err error) {
	// check if type implements object interface
	if object, ok := o.(Object); ok {
		primaryKey = object.PrimaryKey()
		secondaryKey = object.SecondaryKeys()
		return
	}
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr {
		err = fmt.Errorf("crud: pointer expected")
		return
	}
	v = tutils.UnderlyingValue(v)
	if v.Kind() != reflect.Struct {
		err = fmt.Errorf("crud: pointer to struct expected")
		return
	}
	// find primary keys and secondary keys
	primaryKey, secondaryKey, err = getKeys(v)
	return
}

func getKeys(v reflect.Value) (primaryKey PrimaryKey, secondaryKeys []SecondaryKey, err error) {
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// ignore unexported fields
		if field.Anonymous {
			continue
		}
		// get field value
		fieldValue := v.FieldByName(field.Name)
		// check if type implements Index interface
		if index, ok := fieldValue.Interface().(Index); ok {
			sk := index.SecondaryKey()
			err = isValidSecondaryKey(sk)
			if err != nil {
				err = fmt.Errorf("crud: invalid secondary key in field %s on type %T: %w", field.Name, typ.Name(), err)
				return
			}
			// append
			secondaryKeys = append(secondaryKeys, sk)
		}
		// check if type implements
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
			if primaryKey != nil {
				err = fmt.Errorf("crud: only one primary key can be specified for each object")
				return
			}
			valueBytes := marshalValue(fieldValue)
			primaryKey = valueBytes
		case SecondaryKeyTag:
			var prefix []byte
			prefix, err = hex.DecodeString(split[1])
			if err != nil {
				err = fmt.Errorf("crud: invalid hex prefix in secondary key in field %s on type %T", field.Name, typ.Name())
				return
			}
			secondaryKey := SecondaryKey{
				StorePrefix: prefix,
				Key:         marshalValue(fieldValue),
			}
			err = isValidSecondaryKey(secondaryKey)
			if err != nil {
				err = fmt.Errorf("crud: invalid secondary key in field %s on type %T: %w", field.Name, typ.Name(), err)
			}
			secondaryKeys = append(secondaryKeys, secondaryKey)
		}
	}
	if primaryKey == nil {
		err = fmt.Errorf("crud: no primary key specified in type: %T", v.Interface())
		return
	}
	return
}

var typesToBytes = map[reflect.Kind]func(v reflect.Value) []byte{
	reflect.String: func(v reflect.Value) []byte {
		return []byte(v.Interface().(string))
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
		panic("crud: value of type %s cannot be marshaled to bytes")
	}
	return marshaler(v)
}

func marshalSlice(v reflect.Value) []byte {
	if b, ok := v.Interface().([]byte); ok {
		return b
	}
	// todo check if it implements Indexable interface
	panic(fmt.Sprintf("crud: only slice types allowed are byte ones, got: %T", v.Interface()))
}

func isValidSecondaryKey(sk SecondaryKey) (err error) {
	if bytes.Equal(sk.StorePrefix, []byte{PrimaryKeyPrefix}) {
		err = fmt.Errorf("secondary key store prefix contains reserved primary key prefix")
		return
	}
	return
}
