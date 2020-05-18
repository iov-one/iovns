package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
)

// ownerToDomainIndexStore returns the store that indexes all the domains owned by an sdk.AccAddress
func ownerToDomainIndexStore(kvstore sdk.KVStore, addr sdk.AccAddress) (index.Store, error) {
	// check if admin is provided
	if addr.Empty() {
		return index.Store{}, fmt.Errorf("cannot index empty address")
	}
	// get index store
	indexPrefixedStore := indexStore(kvstore)
	// get address to domain index
	return index.NewAddressIndex(indexPrefixedStore, ownerToDomainPrefix, addr)
}

// ownerToAccountIndexStore returns the indexed store that indexes all the accounts owned by an sdk.AccAddress
func ownerToAccountIndexStore(kvstore sdk.KVStore, addr sdk.AccAddress) (index.Store, error) {
	if addr.Empty() {
		return index.Store{}, fmt.Errorf("cannot index empty address")
	}
	// get index store
	indexPrefixedStore := indexStore(kvstore)
	// get address to account index
	return index.NewAddressIndex(indexPrefixedStore, ownerToAccountPrefix, addr)
}
