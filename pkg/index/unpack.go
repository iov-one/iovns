package index

import "fmt"

// Unpacker defines a type that
// can unpack itself from a byte key
type Unpacker interface {
	Unpack(b []byte) error
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
