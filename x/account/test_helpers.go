package account

/*
type subTest struct {
	// BeforeTest is the function run before doing the test,
	// used for example to store state, like configurations etc.
	// Ignored if nil
	BeforeTest func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
	// Test is the function that runs the actual test
	Test func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
	// AfterTest performs actions after the test is run, it can
	// be used to check if the state after Test is run matches
	// the result we expect.
	// Ignored if nil
	AfterTest func(t *testing.T, k keeper.Keeper, ctx sdk.Context)
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




func newTestCodec() *codec.Codec {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	configuration.RegisterCodec(cdc)
	domain.RegisterCodec(cdc)
	account.RegisterCodec(cdc)
	return cdc
}
func newTestKeeper(t *testing.T) {
	cdc := newTestCodec()
	// gen store
	mdb := dbm.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey)
	accountStoreKey := sdk.NewKVStoreKey(account.StoreKey)
	domainStoreKey := sdk.NewKVStoreKey(domain.StoreKey)
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	ms.MountStoreWithDB(accountStoreKey, sdk.StoreTypeIAVL, mdb)       // mount account module
	ms.MountStoreWithDB(domainStoreKey, sdk.StoreTypeIAVL, mdb)        // mount domain module
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create config keeper
	confKeeper := configuration.NewKeeper(cdc, configurationStoreKey, nil)
	// create account keeper
	accountKeeper := keeper.NewKeeper(cdc, accountStoreKey, nil)
	// create context
	ctx := sdk.NewContext(ms, abci.Header{}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return types.NewKeeper(cdc, domainStoreKey, accountKeeper, confKeeper, nil), ctx
}

*/
