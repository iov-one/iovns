package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tm-db"
	"testing"
	"time"
)

// NewTestCodec generates a mock codec for keeper module
func NewTestCodec() *codec.Codec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	configuration.RegisterCodec(cdc)
	return cdc
}

// NewTestKeeper generates a keeper and a context from it
func NewTestKeeper(t *testing.T, isCheckTx bool) (Keeper, sdk.Context) {
	cdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey) // configuration module store key
	accountStoreKey := sdk.NewKVStoreKey(types.DomainStoreKey)         // account module store key
	domainStoreKey := sdk.NewKVStoreKey(types.AccountStoreKey)         // domain module store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	ms.MountStoreWithDB(accountStoreKey, sdk.StoreTypeIAVL, mdb)       // mount account module
	ms.MountStoreWithDB(domainStoreKey, sdk.StoreTypeIAVL, mdb)        // mount domain module
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create config keeper
	confKeeper := configuration.NewKeeper(cdc, configurationStoreKey, nil)
	// create context
	ctx := sdk.NewContext(ms, tmtypes.Header{Time: time.Now()}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return NewKeeper(cdc, domainStoreKey, accountStoreKey, confKeeper, nil), ctx
}
