package index

import sdk "github.com/cosmos/cosmos-sdk/types"

// addrIndexer is just a type alias for sdk.AccAddress
type addrIndexer []byte

// Index wraps normal sdk.AccAddress into a valid Indexer
func (a addrIndexer) Index() ([]byte, error) {
	return a, nil
}

// NewAddressIndex builds an address indexer given a prefix and the address itself
func NewAddressIndex(store sdk.KVStore, prefix []byte, address sdk.AccAddress) (Store, error) {
	return NewIndexedStore(store, prefix, addrIndexer(address))
}

type stringIndexer string

func (s stringIndexer) Index() ([]byte, error) {
	return []byte(s), nil
}

// NewStringIndex builds a string indexer given a prefix and a string
func NewStringIndex(store sdk.KVStore, prefix []byte, str string) (Store, error) {
	return NewIndexedStore(store, prefix, stringIndexer(str))
}
