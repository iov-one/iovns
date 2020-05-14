package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Unpacker defines a type that
// can unpack itself from a byte key
type Unpacker interface {
	Unpack(b []byte) error
}

// Indexed defines an object that can save itself
// into byte data using Pack, and retrive unique info
// about himself from Pack through Unpack
type Indexed interface {
	// Pack marshals the object into a unique byte key
	Pack() ([]byte, error)
}

// Unpack takes an unpacker and fills it based on key
func Unpack(key []byte, unpacker Unpacker) error {
	return unpacker.Unpack(key)
}

// MustUnpack panics if Unpack fails
func MustUnpack(key []byte, unpacker Unpacker) {
	err := unpacker.Unpack(key)
	if err != nil {
		panic(fmt.Sprintf("failure in unpacking %x key at %T: %s", key, unpacker, err))
	}
}

// PackBytes takes a bytes slice and packs it
func PackBytes(keys [][]byte) ([]byte, error) {
	var packedKey []byte
	for _, key := range keys {
		pKey, err := packBytes(key)
		if err != nil {
			return nil, fmt.Errorf("index: PackBytes: %w", err)
		}
		packedKey = append(packedKey, pKey...)
	}
	return packedKey, nil
}

func packBytes(k []byte) ([]byte, error) {
	packed := new(bytes.Buffer)
	if len(k) == 0 {
		return nil, errors.New("empty key")
	}
	if len(k) > math.MaxUint8-1 {
		return nil, fmt.Errorf("key length exceeded: %d", len(k))
	}
	// write size
	err := binary.Write(packed, binary.BigEndian, uint8(len(k)))
	if err != nil {
		return nil, err
	}
	// write content
	_, err = packed.Write(k)
	if err != nil {
		return nil, err
	}
	return packed.Bytes(), nil
}

// UnpackBytes reads a key and returns the
// byte arrays composing said key, it fails
// only if the given key length is not valid
// or if the byte representing the size does not
// match the real length of the key.
func UnpackBytes(k []byte) ([][]byte, error) {
	// check if minimum length is matched
	if len(k) < 2 {
		return nil, fmt.Errorf("minimum length not reached: %d", len(k))
	}
	kCopy := make([]byte, len(k))
	copy(kCopy, k)
	// read size
	var u8size uint8
	packed := bytes.NewBuffer(kCopy)
	err := binary.Read(packed, binary.BigEndian, &u8size)
	if err != nil {
		return nil, err
	}
	size := int(u8size)

	// check if key length minus size byte matches size
	if len(k)-1 < size {
		return nil, fmt.Errorf("invalid key length %d, wanted at least: %d", len(k)-1, size)
	}
	// get key
	packed.Reset()
	var result [][]byte
	// get bytes of key from size_byte_index->size_index
	result = append(result, packed.Bytes()[1:size+1])
	// check if there are more keys
	if len(k) > 1+size {
		remainder := make([]byte, len(k)-1-size)
		// get remainder and process it
		copy(remainder, k[1+size:])
		otherKeys, err := UnpackBytes(remainder)
		if err != nil {
			return nil, err
		}
		result = append(result, otherKeys...)
	}
	return result, nil
}
