package domain

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
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

func TestHandleMsgRegisterDomain(t *testing.T) {
	keyBase := keys.NewInMemory()
	aliceAddr, _, _ := keyBase.CreateMnemonic("alice", keys.English, "", keys.Secp256k1)
	bobAddr, _, _ := keyBase.CreateMnemonic("bob", keys.English, "", keys.Secp256k1)
	testCases := map[string]subTest{
		"success": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				configSetter := getConfigSetter(k.ConfigurationKeeper)
				// set config
				configSetter.SetConfig(ctx, configuration.Config{
					Owner:       nil,
					ValidDomain: "^(.*?)?",
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// register domain with superuser
				_, err := handleMsgRegisterDomain(ctx, k, MsgRegisterDomain{
					Name:         "domain",
					HasSuperuser: true,
					AccountRenew: 10,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterDomain() with superuser, got error: %s", err)
				}
				// TODO register domain without superuser
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// TODO add check domains exists
			},
		},
		"fail domain name exists": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				k.SetDomain(ctx, types.Domain{
					Name:         "exists",
					Admin:        nil,
					ValidUntil:   0,
					HasSuperuser: false,
					AccountRenew: 0,
					Broker:       nil,
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterDomain(ctx, k, MsgRegisterDomain{
					Name:         "exists",
					Admin:        nil,
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 0,
				})
				if !errors.Is(err, types.ErrDomainAlreadyExists) {
					t.Fatalf("handleMsgRegisterDomain() expected: %s got: %s", types.ErrDomainAlreadyExists, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain does not match valid domain regexp": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// get set config function
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidDomain: "$^",
					DomainRenew: 0,
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				_, err := handleMsgRegisterDomain(ctx, k, MsgRegisterDomain{
					Name:         "invalid-name",
					Admin:        nil,
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 0,
				})
				if !errors.Is(err, types.ErrInvalidDomainName) {
					t.Fatalf("handleMsgRegisterDomain() expected error: %s, got: %s", types.ErrInvalidDomainName, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain with no super user must be registered by configuration owner": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// add config with owner
				config := configuration.Config{
					Owner:                  aliceAddr.GetAddress(),
					ValidDomain:            "^(.*?)?",
					ValidName:              "",
					ValidBlockchainID:      "",
					ValidBlockchainAddress: "",
					DomainRenew:            0,
				}
				setConfig := getConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, config)
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context) {
				// try to register domain with no super user
				_, err := handleMsgRegisterDomain(ctx, k, MsgRegisterDomain{
					Name:         "some-domain",
					Admin:        bobAddr.GetAddress(),
					HasSuperuser: false,
					Broker:       nil,
					AccountRenew: 10,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterDomain() expecter error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
	}

	// run all test cases
	runTests(t, testCases)
}
