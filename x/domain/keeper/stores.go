package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/index"
	"github.com/iov-one/iovns/x/domain/types"
)

// ownerToDomainIndexStore returns the store that indexes all the domains owned by an sdk.AccAddress
func ownerToDomainIndexStore(kvstore sdk.KVStore, domain types.Domain) (index.Store, error) {
	// check if admin is provided
	if domain.Admin.Empty() {
		return index.Store{}, fmt.Errorf("cannot index empty address")
	}
	// get index store
	indexPrefixedStore := indexStore(kvstore)
	// get address to domain index
	store, err := index.NewAddressIndex(indexPrefixedStore, ownerToDomainPrefix, domain.Admin)
	if err != nil {
		return index.Store{}, err
	}
	return store, err
}
