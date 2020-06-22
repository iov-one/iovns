package crud

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/iov-one/iovns/tutils"
	"reflect"
	"strings"
)

// TagName is the tag used to marshal crud types
const TagName = "crud"

// PrimaryKeyTag is the tag used to define a primary key
const PrimaryKeyTag = "primaryKey"

// SecondarKeyTag is the value used to define a secondary key in a tag
const SecondaryKeyTag = "secondaryKey"

// PrimaryKeyPrefix
const PrimaryKeyPrefix = 0x0

// these errors exist for testing purposes
var errNotAPointer = errors.New("crud: not a pointer")
var errPointerToStruct = errors.New("crud: pointer to struct expected")
var errIsPrimaryKeyPrefix = errors.New("secondary key store prefix equals to reserved primary key prefix")
var errMultiplePrimaryKeys = errors.New("only one primary key is allowed in each type")
var errNoPrimaryKey = errors.New("no primary key specified in type")
var errEmptyKey = errors.New("provided key is empty")
var errDuplicateKey = errors.New("duplicate key in same prefix")
var errNotAllowedKind = errors.New("kind is not allowed")

func inspect(o interface{}) (primaryKey PrimaryKey, secondaryKeys []SecondaryKey, err error) {
	// check if type implements object interface
	if object, ok := o.(Object); ok {
		primaryKey = object.PrimaryKey()
		secondaryKeys = object.SecondaryKeys()
		return
	}
	v := reflect.ValueOf(o)
	if v.Kind() != reflect.Ptr {
		err = errNotAPointer
		return
	}
	v = tutils.UnderlyingValue(v)
	if v.Kind() != reflect.Struct {
		err = errPointerToStruct
		return
	}
	// find primary keys and secondary keys
	primaryKey, secondaryKeys, err = getKeys(v)
	if err != nil {
		err = fmt.Errorf("crud: %w", err)
		return
	}
	err = validateKeys(primaryKey, secondaryKeys)
	if err != nil {
		err = fmt.Errorf("crud: %w", err)
		return
	}
	return
}

func validateKeys(pk PrimaryKey, sk []SecondaryKey) error {
	// check if pk is empty
	if len(pk) == 0 {
		return fmt.Errorf("primary key: %w", errEmptyKey)
	}
	// check if secondary keys are valid
	var keySet = make(map[string]struct{}, len(sk)) // maps the concatenation of sk.Key and sk.StorePrefix to check for dups

	for _, key := range sk {
		var keyString = string(append(key.StorePrefix, key.Key...))
		if _, ok := keySet[keyString]; ok {
			return fmt.Errorf("prefix %x with key %x: %w", key.StorePrefix, key.Key, errDuplicateKey)
		}
		err := isValidSecondaryKey(key)
		if err != nil {
			return err
		}
		keySet[keyString] = struct{}{}
	}
	return nil
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
		iface := fieldValue.Interface()
		// check if type inherently implements Index
		if index, ok := iface.(Index); ok {
			sk := index.SecondaryKey()
			// append
			secondaryKeys = append(secondaryKeys, sk)
			// check if slice and if every single element implements index
		} else if fieldValue.Kind() == reflect.Slice {
			slLen := fieldValue.Len()
			for i := 0; i < slLen; i++ {
				index, ok := fieldValue.Index(i).Interface().(Index)
				if !ok {
					continue
				}
				sk := index.SecondaryKey()
				secondaryKeys = append(secondaryKeys, sk)
			}
		}
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
				err = fmt.Errorf("%w: %s", errMultiplePrimaryKeys, typ.String())
				return
			}
			var valueBytes []byte
			valueBytes, err = marshalValue(fieldValue)
			if err != nil {
				return
			}
			primaryKey = valueBytes
		case SecondaryKeyTag:
			var prefix []byte
			prefix, err = hex.DecodeString(split[1])
			if err != nil {
				err = fmt.Errorf("invalid hex prefix in secondary key in field %s on type %T", field.Name, v.Interface())
				return
			}
			var valueBytes []byte
			valueBytes, err = marshalValue(fieldValue)
			if err != nil {
				return
			}
			secondaryKey := SecondaryKey{
				StorePrefix: prefix,
				Key:         valueBytes,
			}
			err = isValidSecondaryKey(secondaryKey)
			if err != nil {
				err = fmt.Errorf("invalid secondary key in field %s on type %T: %w", field.Name, v.Interface(), err)
			}
			secondaryKeys = append(secondaryKeys, secondaryKey)
		}
	}
	if primaryKey == nil {
		err = fmt.Errorf("%w: %T", errNoPrimaryKey, v.Interface())
		return
	}
	return
}

var typesToBytes = map[reflect.Kind]func(v reflect.Value) ([]byte, error){
	reflect.String: func(v reflect.Value) ([]byte, error) {
		return []byte(v.Interface().(string)), nil
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
func marshalValue(v reflect.Value) ([]byte, error) {
	v = tutils.UnderlyingValue(v)
	// check if forbidden type
	kind := v.Kind()
	if kind == reflect.Slice {
		return marshalSlice(v)
	}
	if _, ok := notAllowedIndexType[kind]; ok {
		return nil, fmt.Errorf("value of type %s cannot be turned into a byte key", kind)
	}
	// now index based on type
	marshaler, ok := typesToBytes[kind]
	if !ok {
		return nil, fmt.Errorf("value of type %s cannot be marshaled to bytes", v.Type().String())
	}
	return marshaler(v)
}

func marshalSlice(v reflect.Value) ([]byte, error) {
	if b, ok := v.Interface().([]byte); ok {
		return b, nil
	}
	return nil, fmt.Errorf("only slice types allowed are byte ones, got: %T", v.Interface())
}

func isValidSecondaryKey(sk SecondaryKey) (err error) {
	if bytes.Equal(sk.StorePrefix, []byte{PrimaryKeyPrefix}) {
		return errIsPrimaryKeyPrefix
	}
	if len(sk.Key) == 0 {
		return errEmptyKey
	}
	return
}
