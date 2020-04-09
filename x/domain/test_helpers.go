package domain

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/account"
	"github.com/iov-one/iovnsd/x/configuration"
	"github.com/iov-one/iovnsd/x/domain/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

// subTest defines a test runner
type subTest struct {
	// BeforeTest is the function run before doing the test,
	// used for example to store state, like configurations etc.
	// Ignored if nil
	BeforeTest func(t *testing.T, k Keeper, ctx sdk.Context)
	// Test is the function that runs the actual test
	Test func(t *testing.T, k Keeper, ctx sdk.Context)
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k Keeper, ctx sdk.Context)
}

// runTests run tests cases after generating a new keeper and context for each test case
func runTests(t *testing.T, tests map[string]subTest) {
	for name, test := range tests {
		keeper, ctx := newTestKeeper(t, true)
		// run sub subTest
		t.Run(name, func(t *testing.T) {
			// run before subTest
			if test.BeforeTest != nil {
				test.BeforeTest(t, keeper, ctx)
			}
			// run actual subTest
			test.Test(t, keeper, ctx)
			// run after subTest
			if test.AfterTest != nil {
				test.AfterTest(t, keeper, ctx)
			}
		})
	}
}

// newTestCodec generates a mock codec for keeper module
func newTestCodec() *codec.Codec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	configuration.RegisterCodec(cdc)
	account.RegisterCodec(cdc)
	return cdc
}

// newTestKeeper generates a keeper and a context from it
func newTestKeeper(t *testing.T, isCheckTx bool) (Keeper, sdk.Context) {
	cdc := newTestCodec()
	// generate store
	mdb := dbm.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey) // configuration module store key
	accountStoreKey := sdk.NewKVStoreKey(account.StoreKey)             // account module store key
	domainStoreKey := sdk.NewKVStoreKey(types.StoreKey)                // domain module store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	ms.MountStoreWithDB(accountStoreKey, sdk.StoreTypeIAVL, mdb)       // mount account module
	ms.MountStoreWithDB(domainStoreKey, sdk.StoreTypeIAVL, mdb)        // mount domain module
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create config keeper
	confKeeper := configuration.NewKeeper(cdc, configurationStoreKey, nil)
	// create account keeper
	accountKeeper := account.NewKeeper(cdc, accountStoreKey, nil)
	// create context
	ctx := sdk.NewContext(ms, abci.Header{}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return NewKeeper(cdc, domainStoreKey, accountKeeper, confKeeper, nil), ctx
}

// since the exposed interface for configuration keeper
// does not include set config as the domain module should
// not be able to change configuration state, then only
// in test cases we expose this method
type configurationSetter interface {
	SetConfig(ctx sdk.Context, config configuration.Config)
}

// getConfigSetter exposes the configurationSetter interface
// allowing the module to set configuration state, this should only
// be used for tests and will panic if the keeper provided can not
// be cast to configurationSetter
func getConfigSetter(keeper ConfigurationKeeper) configurationSetter {
	// check if the configuration keeper is also a config setter
	configSetter, ok := keeper.(configurationSetter)
	if !ok {
		panic(fmt.Sprintf("cannot cast configuration keeper to configuration setter: got uncastable type: %T", keeper))
	}
	return configSetter
}
