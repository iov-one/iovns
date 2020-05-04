package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/iov-one/iovns/mock"
	"github.com/iov-one/iovns/x/configuration"
	confCdc "github.com/iov-one/iovns/x/configuration/types"
	"github.com/iov-one/iovns/x/domain/types"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tm-db"
	"testing"
	"time"
)

// NewTestCodec generates aliceAddr mock codec for keeper module
func NewTestCodec() *codec.Codec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	confCdc.RegisterCodec(cdc)
	return cdc
}

type Mocks struct {
	Supply *mock.SupplyKeeperMock
}

// NewTestKeeper generates aliceAddr keeper and aliceAddr context from it
func NewTestKeeper(t testing.TB, isCheckTx bool) (Keeper, sdk.Context, *Mocks) {
	cdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey) // configuration module store key
	accountStoreKey := sdk.NewKVStoreKey(types.DomainStoreKey)         // account module store key
	domainStoreKey := sdk.NewKVStoreKey(types.AccountStoreKey)         // domain module store key
	indexStoreKey := sdk.NewKVStoreKey(types.IndexStoreKey)            // index store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	ms.MountStoreWithDB(accountStoreKey, sdk.StoreTypeIAVL, mdb)       // mount account module
	ms.MountStoreWithDB(domainStoreKey, sdk.StoreTypeIAVL, mdb)        // mount domain module
	ms.MountStoreWithDB(indexStoreKey, sdk.StoreTypeIAVL, mdb)
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create Mocks
	mocks := new(Mocks)
	// create mock supply keeper
	mocks.Supply = mock.NewSupplyKeeper()
	// create config keeper
	confKeeper := configuration.NewKeeper(cdc, configurationStoreKey, subspace.NewSubspace(cdc, nil, nil, "test"))
	// create context
	ctx := sdk.NewContext(ms, tmtypes.Header{Time: time.Now()}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return NewKeeper(cdc, domainStoreKey, accountStoreKey, indexStoreKey, confKeeper, mocks.Supply.Mock(), nil), ctx, mocks
}
