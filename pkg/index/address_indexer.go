package index

import sdk "github.com/cosmos/cosmos-sdk/types"

type addrIndexer []byte

func (a addrIndexer) Index() ([]byte, error) {
	return a, nil
}

func NewAddressIndex(store sdk.KVStore, prefix []byte, address sdk.AccAddress) (Store, error) {
	return NewIndexedStore(store, prefix, addrIndexer(address))
}
