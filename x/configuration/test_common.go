package configuration

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration/types"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

// NewTestCodec generates aliceAddr mock codec for keeper module
func NewTestCodec() *codec.Codec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return cdc
}

// NewTestKeeper generates aliceAddr keeper and aliceAddr context from it
func NewTestKeeper(t testing.TB, isCheckTx bool) (Keeper, sdk.Context) {
	cdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(StoreKey) // configuration module store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create context
	ctx := sdk.NewContext(ms, tmtypes.Header{Time: time.Now()}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return NewKeeper(cdc, configurationStoreKey, nil), ctx
}
