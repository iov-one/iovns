package index

import "fmt"

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
