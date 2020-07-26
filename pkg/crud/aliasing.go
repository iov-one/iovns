package crud

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/pkg/crud/internal/store"
	"github.com/iov-one/iovns/pkg/crud/types"
)

// NewStore returns a new CRUD key value store
func NewStore(ctx sdk.Context, key sdk.StoreKey, cdc *codec.Codec, uniquePrefix []byte) types.Store {
	return store.NewStore(ctx, key, cdc, uniquePrefix)
}

