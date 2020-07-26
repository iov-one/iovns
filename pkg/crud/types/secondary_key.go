package types

import (
	"bytes"
	"encoding/base64"
)

const ReservedSeparator = 0xFF

// SecondaryKey defines the secondary key behaviour
type SecondaryKey interface {
	Marshal() []byte
	Unmarshal(b []byte)
	Key() []byte
	Prefix() byte
}


func NewSecondaryKey(storePrefix byte, key []byte) SecondaryKey {
	return &secondaryKey{
		key:         fixKey(key),
		storePrefix: storePrefix,
	}
}

func NewSecondaryKeyFromBytes(b []byte) SecondaryKey {
	sk := &secondaryKey{}
	sk.Unmarshal(b)
	return sk
}

// SecondaryKey defines a secondary key for the object
type secondaryKey struct {
	// key is the byte key which identifies the byte key prefix used to iterate of the index of the secondary key
	key []byte
	// storePrefix is the prefix of the index, necessary to divide one index from another
	storePrefix byte
}

func (s secondaryKey) Marshal() []byte {
	result := make([]byte, 0, 1 + len(s.key))
	result = append(result, s.storePrefix)
	result = append(result, s.key...)
	return result
}

func (s *secondaryKey) Unmarshal(b []byte) {
	// at least three bytes define an index
	// store prefix + index key + separator
	if len(b) < 3 {
		panic("cannot unmarshal invalid length byte slice")
	}
	s.storePrefix = b[0]
	s.key = b[1:]
}

func (s *secondaryKey) Prefix() byte {
	return s.storePrefix
}

func (s *secondaryKey) Key() []byte {
	cpy := make([]byte, len(s.key))
	copy(cpy, s.key)
	return cpy
}

// fixKey encodes key value which contains the reserved separator by base64-encoding them
// this is necessary because we're dealing with a prefixed KVStore which, if we iterate, is going to
// iterate over bytes contained in a key, so if we assume we have:
// KeyA = [0x1]
// KeyB = [0x1, 0x2]
// during iteration, in case we wanted to iterate over KeyA only we'd end up in KeyB domain too because
// KeyB starts with KeyA, so to avoid this we put a full stop separator which we know other keys can not contain
func fixKey(b []byte) []byte {
	if bytes.Contains(b, []byte{ReservedSeparator}) {
		dst := make([]byte, base64.RawStdEncoding.EncodedLen(len(b)))
		base64.RawStdEncoding.Encode(dst, b)
		return append(dst, ReservedSeparator)
	}
	return append(b, ReservedSeparator)
}
