package crud

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/iov-one/iovns/tutils"
	"reflect"
)

// TagName is the tag used to marshal crud types
const TagName = "crud"

// PrimaryKeyTag is the tag used to define a primary key
const PrimaryKeyTag = "primaryKey"

// PrimaryKeyPrefix
const PrimaryKeyPrefix = 0x0

var primaryKeyPrefix = []byte{PrimaryKeyPrefix}

// these errors exist for testing purposes
var errNotAPointer = errors.New("crud: not a pointer")
var errPointerToStruct = errors.New("crud: pointer to struct expected")
var errIsPrimaryKeyPrefix = errors.New("secondary key store prefix equals to reserved primary key prefix")
var errMultiplePrimaryKeys = errors.New("only one primary key is allowed in each type")
var errNoPrimaryKey = errors.New("no primary key specified in type")
var errEmptyKey = errors.New("provided key is empty")
var errDuplicateKey = errors.New("duplicate key in same prefix")
var errNotAllowedKind = errors.New("kind is not allowed")
var errNotValidSliceType = errors.New("provided slice type is not valid")

func inspect(o interface{}) (primaryKey PrimaryKey, secondaryKeys []SecondaryKey, err error) {
	// find primary keys and secondary keys
	primaryKey, secondaryKeys, err = getKeys(o)
	if err != nil {
		err = fmt.Errorf("crud: %w", err)
		return
	}
	// validate
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

func getKeys(o interface{}) (primaryKey PrimaryKey, secondaryKeys []SecondaryKey, err error) {
	// check if type implements object interface
	if object, ok := o.(Object); ok {
		primaryKey = object.PrimaryKey()
		secondaryKeys = object.SecondaryKeys()
		return
	}
	v := reflect.ValueOf(o)
	v = tutils.UnderlyingValue(v)
	if v.Kind() != reflect.Struct {
		err = errPointerToStruct
		return
	}
	typ := v.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// ignore unexported fields
		if field.Anonymous {
			continue
		}
		// check if field must be parsed
		tagValue, ok := field.Tag.Lookup(TagName)
		// if tagValue is missing then no indexing is required
		if !ok {
			continue
		}
		// get prefix
		var prefix []byte
		prefix, err = getPrefixFromTag(tagValue)
		if err != nil {
			return
		}
		// get field value
		fieldValue := v.FieldByName(field.Name)
		var keys [][]byte
		keys, err = marshal(fieldValue)
		if err != nil {
			return
		}
		if len(keys) == 0 {
			continue
		}

		switch bytes.Equal(prefix, primaryKeyPrefix) {
		case true:
			if primaryKey != nil {
				err = errMultiplePrimaryKeys
				return
			}
			if len(keys) != 1 {
				err = fmt.Errorf("unexpected number of byte keys")
				return
			}
			primaryKey = keys[0]
		default:
			sks := make([]SecondaryKey, len(keys))
			for i, key := range keys {
				sks[i] = SecondaryKey{
					Key:         key,
					StorePrefix: prefix,
				}
			}
			secondaryKeys = append(secondaryKeys, sks...)
		}
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

type Hashable interface {
	Hashes() [][]byte
}

func marshal(v reflect.Value) ([][]byte, error) {
	kind := v.Kind()
	// check if slice
	if kind == reflect.Slice {
		return marshalArray(v)
	}
	// check if type implements Hashable
	i := v.Interface()
	if hashable, ok := i.(Hashable); ok {
		return hashable.Hashes(), nil
	}
	// otherwise if it does not implement the interface we need to automatically generate the hash
	key, err := marshalValue(v)
	if err != nil {
		return nil, err
	}
	return [][]byte{key}, nil
}

func marshalArray(v reflect.Value) (keys [][]byte, err error) {
	// check if array type implements interface
	l := v.Len()
	// no element to marshal
	if l == 0 {
		return nil, nil
	}
	// otherwise check if elements inside the slice implement Hashable interface
	for i := 0; i < l; i++ {
		obj := v.Index(i)
		iface := obj.Interface()
		hashable, ok := iface.(Hashable)
		if !ok {
			// technically we could 'break' here but if we want to support []interface{} in which some implement Hashable and some not then we have to continue
			continue
		}
		keys = append(keys, hashable.Hashes()...)
	}
	// if length of keys is > 0 then we have already done the marshalling and we can quit
	if len(keys) != 0 {
		return
	}
	// otherwise keep going
	iface := v.Interface()
	if b, ok := iface.([]byte); ok {
		keys = [][]byte{b}
		return
	}
	// un-marshalable type
	err = fmt.Errorf("%T: %w", iface, errNotValidSliceType)
	return
}

// getPrefixFromTag returns the kv store prefix extracted from the tag
func getPrefixFromTag(value string) ([]byte, error) {
	if value == PrimaryKeyTag {
		return primaryKeyPrefix, nil
	}
	return hex.DecodeString(value)
}
