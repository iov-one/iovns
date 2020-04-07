package domain

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovnsd/x/account"
	"github.com/iov-one/iovnsd/x/configuration"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

// subTest defines a test runner
type subTest struct {
	// Name is the name of the test
	Name string
	// BeforeTest is the function run before doing the test,
	// used for example to store state, like configurations etc.
	// Ignored if nil
	BeforeTest func(t *testing.T, k Keeper, ctx sdk.Context)
	// Test runs the actual test
	Test func(t *testing.T, k Keeper, ctx sdk.Context)
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k Keeper, ctx sdk.Context)
}

func runTests(t *testing.T, tests []subTest) {
	keeper, ctx := newTestKeeper(t, true)
	for _, test := range tests {
		// run sub subTest
		t.Run(test.Name, func(t *testing.T) {
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

// newTestCodec generates a codec for keeper module
func newTestCodec() *codec.Codec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	configuration.RegisterCodec(cdc)
	account.RegisterCodec(cdc)
	return cdc
}

func newTestKeeper(t *testing.T, isCheckTx bool) (Keeper, sdk.Context) {
	cdc := newTestCodec()
	// generate store
	mdb := dbm.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey) // configuration module store key
	accountStoreKey := sdk.NewKVStoreKey(account.StoreKey)             // account module store key
	domainStoreKey := sdk.NewKVStoreKey(StoreKey)                      // domain module store key
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

func TestMsgRegisterDomain(t *testing.T) {
	test := subTest{
		Name: "success",
		BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
			// add config
			type configurationSetter interface {
				SetConfig(ctx sdk.Context, config configuration.Config)
			}
			// check if the configuration keeper is also a config setter
			configSetter, ok := k.ConfigurationKeeper.(configurationSetter)
			if !ok {
				t.Fatalf("handleMsgRegisterDomain() cannot cast configuration keeper to configuration setter: got uncastable type: %T", k.ConfigurationKeeper)
			}
			// set config
			configSetter.SetConfig(ctx, configuration.Config{
				Owner:                  nil,
				ValidDomain:            "^(.*?)?",
				ValidName:              "",
				ValidBlockchainID:      "",
				ValidBlockchainAddress: "",
				DomainRenew:            0,
			})
		},
		Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
			_, err := handleMsgRegisterDomain(ctx, k, MsgRegisterDomain{
				Name:         "domain",
				Admin:        nil,
				HasSuperuser: true,
				Broker:       nil,
				AccountRenew: 10,
			})
			if err != nil {
				t.Fatalf("handleMsgRegisterDomain() got error: %s", err)
			}
		},
		AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context) {

		},
	}
	runTests(t, []subTest{test})
}
